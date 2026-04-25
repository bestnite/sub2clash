package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bestnite/sub2clash/logger"
	"github.com/bestnite/sub2clash/model"
	P "github.com/bestnite/sub2clash/model/proxy"
	"github.com/bestnite/sub2clash/parser"
	"github.com/bestnite/sub2clash/utils"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var subsDir = "subs"
var fileLock sync.RWMutex

func LoadSubscription(url string, refresh bool, userAgent string, cacheExpire int64, retryTimes int) ([]byte, error) {
	if refresh {
		return FetchSubscriptionFromAPI(url, userAgent, retryTimes)
	}
	hash := sha256.Sum224([]byte(url))
	fileName := filepath.Join(subsDir, hex.EncodeToString(hash[:]))
	stat, err := os.Stat(fileName)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return FetchSubscriptionFromAPI(url, userAgent, retryTimes)
	}
	lastGetTime := stat.ModTime().Unix()
	if lastGetTime+cacheExpire > time.Now().Unix() {
		file, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		defer func(file *os.File) {
			if file != nil {
				_ = file.Close()
			}
		}(file)
		fileLock.RLock()
		defer fileLock.RUnlock()
		subContent, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		return subContent, nil
	}
	return FetchSubscriptionFromAPI(url, userAgent, retryTimes)
}

func FetchSubscriptionFromAPI(url string, userAgent string, retryTimes int) ([]byte, error) {
	hash := sha256.Sum224([]byte(url))
	fileName := filepath.Join(subsDir, hex.EncodeToString(hash[:]))
	client := Request(retryTimes)
	defer client.Close()
	resp, err := client.R().SetHeader("User-Agent", userAgent).Get(url)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		if file != nil {
			_ = file.Close()
		}
	}(file)
	fileLock.Lock()
	defer fileLock.Unlock()
	_, err = file.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write to sub.yaml: %w", err)
	}
	return data, nil
}

// BuildSub 是当前配置转换链路的核心入口。
//
// 当前设计分为三层：
// 1. templateDoc：模板 YAML 的完整语法树，也是最终输出真源
// 2. generatedConfig：本项目运行期最小叠加层，只保存参与业务计算的字段
// 3. proxy.Proxy：节点解析后的 typed 模型，用于过滤、去重、重命名和输出
//
// 这个函数的目标不是“重建一整份 mihomo 配置”，而是：
// - 保留模板中绝大部分原始字段
// - 只对 proxies / proxy-groups / rules / rule-providers 做定点 patch
func BuildSub(clashType model.ClashType, query model.ConvertConfig, template string, cacheExpire int64, retryTimes int) (
	*BuiltSub, error,
) {
	templateDoc, templateBytes, err := loadTemplateDocument(query, template, cacheExpire, retryTimes)
	if err != nil {
		return nil, err
	}

	temp, err := extractTemplateOverlay(templateDoc)
	if err != nil {
		logger.Logger.Debug("extract template overlay failed", zap.Error(err))
		return nil, NewTemplateParseError(templateBytes, err)
	}
	proxyList, err := collectQueryProxies(query, cacheExpire, retryTimes)
	if err != nil {
		return nil, err
	}

	proxyList, err = normalizeProxyList(query, proxyList)
	if err != nil {
		return nil, err
	}

	// t 仅承载“由节点生成出来的新内容”，例如国家组。
	// 模板里原有的组、规则等则保存在 temp 中。
	generated, err := buildGeneratedConfig(clashType, query, proxyList)
	if err != nil {
		return nil, err
	}

	MergeSubAndTemplate(temp, generated, query.IgnoreCountryGrooup)

	applyRulePatches(temp, query)
	addedRuleProviders := buildRuleProviderPatches(query)

	if err := mergeTemplateProxies(templateDoc, generated.Proxy); err != nil {
		return nil, NewError(ErrConfigInvalid, "failed to update template path: proxies", err)
	}

	if temp.ProxyGroup == nil {
		temp.ProxyGroup = make([]generatedGroup, 0)
	}
	if err := mergeTemplateProxyGroups(templateDoc, temp.ProxyGroup); err != nil {
		return nil, NewError(ErrConfigInvalid, "failed to update template path: proxy-groups", err)
	}

	rulesChanged := len(query.Rules) != 0 || len(query.RuleProviders) != 0
	if rulesChanged {
		if temp.Rule == nil {
			temp.Rule = make([]string, 0)
		}
		if err := SetYAMLPath(templateDoc, "rules", temp.Rule); err != nil {
			return nil, NewError(ErrConfigInvalid, "failed to update template path: rules", err)
		}
	}

	if len(query.RuleProviders) != 0 {
		if err := mergeTemplateRuleProviders(templateDoc, addedRuleProviders); err != nil {
			return nil, NewError(ErrConfigInvalid, "failed to update template path: rule-providers", err)
		}
	}

	return &BuiltSub{root: templateDoc}, nil
}

// loadTemplateDocument 负责统一加载模板来源，并返回：
// 1. 解析后的 YAML 语法树
// 2. 原始模板字节，用于错误报告
func loadTemplateDocument(query model.ConvertConfig, template string, cacheExpire int64, retryTimes int) (*yaml.Node, []byte, error) {
	var err error
	var templateBytes []byte

	if query.Template != "" {
		template = query.Template
	}
	if strings.HasPrefix(template, "http") {
		templateBytes, err = LoadSubscription(template, query.Refresh, query.UserAgent, cacheExpire, retryTimes)
		if err != nil {
			logger.Logger.Debug(
				"load template failed", zap.String("template", template), zap.Error(err),
			)
			return nil, nil, NewTemplateLoadError(template, err)
		}
	} else {
		unescape, err := url.QueryUnescape(template)
		if err != nil {
			return nil, nil, NewTemplateLoadError(template, err)
		}
		templateBytes, err = LoadTemplate(unescape)
		if err != nil {
			logger.Logger.Debug(
				"load template failed", zap.String("template", template), zap.Error(err),
			)
			return nil, nil, NewTemplateLoadError(unescape, err)
		}
	}

	templateDoc, err := ParseYAMLDocument(templateBytes)
	if err != nil {
		logger.Logger.Debug("parse template yaml node failed", zap.Error(err))
		return nil, templateBytes, NewTemplateParseError(templateBytes, err)
	}

	return templateDoc, templateBytes, nil
}

// collectQueryProxies 汇总来自订阅链接和直接传入代理链接的所有节点。
func collectQueryProxies(query model.ConvertConfig, cacheExpire int64, retryTimes int) ([]P.Proxy, error) {
	proxyList := make([]P.Proxy, 0)
	for i := range query.Subs {
		newProxies, err := loadSubscriptionProxies(query, query.Subs[i], cacheExpire, retryTimes)
		if err != nil {
			return nil, err
		}
		proxyList = append(proxyList, newProxies...)
	}

	if len(query.Proxies) != 0 {
		p, err := parser.ParseProxies(parser.ParseConfig{UseUDP: query.UseUDP}, query.Proxies...)
		if err != nil {
			return nil, err
		}
		proxyList = append(proxyList, p...)
	}

	return proxyList, nil
}

// loadSubscriptionProxies 负责加载单条订阅并应用订阅名作为节点前缀。
func loadSubscriptionProxies(query model.ConvertConfig, subscriptionURL string, cacheExpire int64, retryTimes int) ([]P.Proxy, error) {
	data, err := LoadSubscription(subscriptionURL, query.Refresh, query.UserAgent, cacheExpire, retryTimes)
	if err != nil {
		logger.Logger.Debug(
			"load subscription failed", zap.String("url", subscriptionURL), zap.Error(err),
		)
		return nil, NewSubscriptionLoadError(subscriptionURL, err)
	}

	subName := ""
	if strings.Contains(subscriptionURL, "#") {
		subName = subscriptionURL[strings.LastIndex(subscriptionURL, "#")+1:]
	}

	newProxies, err := parseSubscriptionProxies(data, query.UseUDP, subscriptionURL)
	if err != nil {
		return nil, err
	}

	if subName != "" {
		for i := range newProxies {
			newProxies[i].SubName = subName
		}
	}

	return newProxies, nil
}

// parseSubscriptionProxies 按“Clash YAML -> URI 列表 -> Base64 文本”的顺序容错解析节点。
func parseSubscriptionProxies(data []byte, useUDP bool, subscriptionURL string) ([]P.Proxy, error) {
	sub := &proxyListDoc{}
	if err := yaml.Unmarshal(data, sub); err == nil {
		return sub.Proxy, nil
	}

	reg, err := regexp.Compile("(" + strings.Join(parser.GetAllPrefixes(), "|") + ")://")
	if err != nil {
		logger.Logger.Debug("compile regex failed", zap.Error(err))
		return nil, NewRegexInvalidError("prefix", err)
	}

	if reg.Match(data) {
		return parser.ParseProxies(parser.ParseConfig{UseUDP: useUDP}, strings.Split(string(data), "\n")...)
	}

	base64, err := utils.DecodeBase64(string(data), false)
	if err != nil {
		logger.Logger.Debug(
			"parse subscription failed", zap.String("url", subscriptionURL),
			zap.String("data", string(data)),
			zap.Error(err),
		)
		return nil, NewSubscriptionParseError(data, err)
	}

	return parser.ParseProxies(parser.ParseConfig{UseUDP: useUDP}, strings.Split(base64, "\n")...)
}

// normalizeProxyList 汇总所有节点标准化步骤，确保后续分组和 patch 使用的是稳定结果。
func normalizeProxyList(query model.ConvertConfig, proxyList []P.Proxy) ([]P.Proxy, error) {
	applySubscriptionPrefixes(proxyList)

	var err error
	proxyList, err = dedupeProxies(proxyList)
	if err != nil {
		return nil, err
	}

	proxyList, err = removeProxiesByPattern(proxyList, query.Remove)
	if err != nil {
		return nil, err
	}

	proxyList, err = replaceProxyNames(proxyList, query.Replace)
	if err != nil {
		return nil, err
	}

	ensureUniqueProxyNames(proxyList)
	trimProxyNames(proxyList)
	return proxyList, nil
}

func applySubscriptionPrefixes(proxyList []P.Proxy) {
	for i := range proxyList {
		if proxyList[i].SubName != "" {
			proxyList[i].Name = strings.TrimSpace(proxyList[i].SubName) + " " + strings.TrimSpace(proxyList[i].Name)
		}
	}
}

// dedupeProxies 通过 YAML 序列化结果判定两个节点是否完全相同。
func dedupeProxies(proxyList []P.Proxy) ([]P.Proxy, error) {
	proxies := make(map[string]*P.Proxy)
	newProxies := make([]P.Proxy, 0, len(proxyList))
	for i := range proxyList {
		yamlBytes, err := yaml.Marshal(proxyList[i])
		if err != nil {
			logger.Logger.Debug("marshal proxy failed", zap.Error(err))
			return nil, fmt.Errorf("marshal proxy failed: %w", err)
		}
		key := string(yamlBytes)
		if _, exist := proxies[key]; !exist {
			proxies[key] = &proxyList[i]
			newProxies = append(newProxies, proxyList[i])
		}
	}
	return newProxies, nil
}

func removeProxiesByPattern(proxyList []P.Proxy, pattern string) ([]P.Proxy, error) {
	if strings.TrimSpace(pattern) == "" {
		return proxyList, nil
	}

	removeReg, err := regexp.Compile(pattern)
	if err != nil {
		logger.Logger.Debug("remove regexp compile failed", zap.Error(err))
		return nil, NewRegexInvalidError("remove", err)
	}

	newProxyList := make([]P.Proxy, 0, len(proxyList))
	for i := range proxyList {
		if removeReg.MatchString(proxyList[i].Name) {
			continue
		}
		newProxyList = append(newProxyList, proxyList[i])
	}
	return newProxyList, nil
}

func replaceProxyNames(proxyList []P.Proxy, replacements map[string]string) ([]P.Proxy, error) {
	if len(replacements) == 0 {
		return proxyList, nil
	}

	for pattern, replacement := range replacements {
		replaceReg, err := regexp.Compile(pattern)
		if err != nil {
			logger.Logger.Debug("replace regexp compile failed", zap.Error(err))
			return nil, NewRegexInvalidError("replace", err)
		}
		for i := range proxyList {
			if replaceReg.MatchString(proxyList[i].Name) {
				proxyList[i].Name = replaceReg.ReplaceAllString(proxyList[i].Name, replacement)
			}
		}
	}

	return proxyList, nil
}

func ensureUniqueProxyNames(proxyList []P.Proxy) {
	names := make(map[string]int)
	for i := range proxyList {
		if _, exist := names[proxyList[i].Name]; exist {
			names[proxyList[i].Name] = names[proxyList[i].Name] + 1
			proxyList[i].Name = proxyList[i].Name + " " + strconv.Itoa(names[proxyList[i].Name])
		} else {
			names[proxyList[i].Name] = 0
		}
	}
}

func trimProxyNames(proxyList []P.Proxy) {
	for i := range proxyList {
		proxyList[i].Name = strings.TrimSpace(proxyList[i].Name)
	}
}

// buildGeneratedConfig 只生成“新增内容”，例如国家组和最终可输出的节点集合。
func buildGeneratedConfig(clashType model.ClashType, query model.ConvertConfig, proxyList []P.Proxy) (*generatedConfig, error) {
	generated := &generatedConfig{}
	AddProxy(generated, query.AutoTest, query.Lazy, clashType, proxyList...)
	sortGeneratedGroups(generated, query.Sort)
	return generated, nil
}

func sortGeneratedGroups(generated *generatedConfig, sortMode string) {
	switch sortMode {
	case "sizeasc":
		sort.Sort(generatedGroupsSortBySize(generated.ProxyGroup))
	case "sizedesc":
		sort.Sort(sort.Reverse(generatedGroupsSortBySize(generated.ProxyGroup)))
	case "nameasc":
		sort.Sort(generatedGroupsSortByName(generated.ProxyGroup))
	case "namedesc":
		sort.Sort(sort.Reverse(generatedGroupsSortByName(generated.ProxyGroup)))
	default:
		sort.Sort(generatedGroupsSortByName(generated.ProxyGroup))
	}
}

// applyRulePatches 只修改运行期 overlay 中的 rules 切片，不直接写 YAML。
func applyRulePatches(temp *generatedConfig, query model.ConvertConfig) {
	for _, v := range query.Rules {
		if v.Prepend {
			PrependRules(temp, v.Rule)
		} else {
			AppendRules(temp, v.Rule)
		}
	}
	for _, v := range query.RuleProviders {
		if v.Prepend {
			PrependRuleProvider(temp, v.Name, v.Group)
		} else {
			AppenddRuleProvider(temp, v.Name, v.Group)
		}
	}
}

// buildRuleProviderPatches 把 API 请求中的 rule-provider 参数转换成 YAML patch payload。
func buildRuleProviderPatches(query model.ConvertConfig) map[string]generatedRulePatch {
	if len(query.RuleProviders) == 0 {
		return nil
	}

	patches := make(map[string]generatedRulePatch, len(query.RuleProviders))
	for _, v := range query.RuleProviders {
		hash := sha256.Sum224([]byte(v.Url))
		name := hex.EncodeToString(hash[:])
		patches[v.Name] = generatedRulePatch{
			Type:     "http",
			Behavior: v.Behavior,
			Url:      v.Url,
			Path:     "./" + name + ".yaml",
			Interval: 3600,
		}
	}
	return patches
}

// extractTemplateOverlay 只从模板 YAML 树中提取本项目真正会参与计算的局部字段。
// 这让模板读取完全基于 yaml.Node，而不再依赖任何整份配置的 typed unmarshal。
func extractTemplateOverlay(templateDoc *yaml.Node) (*generatedConfig, error) {
	overlay := &generatedConfig{}

	if err := decodeOptionalYAMLPath(templateDoc, "proxy-groups", &overlay.ProxyGroup); err != nil {
		return nil, err
	}
	if err := decodeOptionalYAMLPath(templateDoc, "rules", &overlay.Rule); err != nil {
		return nil, err
	}

	return overlay, nil
}

// decodeOptionalYAMLPath 在路径存在且非 null 时才执行 Decode，
// 路径不存在时保持目标值为零值。
func decodeOptionalYAMLPath(doc *yaml.Node, path string, target any) error {
	node, err := GetYAMLPath(doc, path)
	if err != nil {
		return err
	}
	if node == nil || isNullYAMLNode(node) {
		return nil
	}
	if err := node.Decode(target); err != nil {
		return fmt.Errorf("decode template path %q failed: %w", path, err)
	}
	return nil
}

// mergeTemplateProxies 只负责把本项目生成出的代理追加到模板现有 proxies 后面。
// 模板中已有代理节点原样保留，不做 struct round-trip。
func mergeTemplateProxies(templateDoc *yaml.Node, generated []P.Proxy) error {
	if len(generated) == 0 && !HasYAMLPath(templateDoc, "proxies") {
		return nil
	}

	proxiesNode, err := EnsureYAMLSequencePath(templateDoc, "proxies")
	if err != nil {
		return err
	}

	for _, proxy := range generated {
		if err := AppendYAMLSequenceValue(proxiesNode, proxy); err != nil {
			return err
		}
	}

	return nil
}

// mergeTemplateProxyGroups 负责两类更新：
// 1. 对模板中同名组，仅覆盖 proxies 字段，保留其他字段
// 2. 追加本项目新生成的国家组
func mergeTemplateProxyGroups(templateDoc *yaml.Node, groups []generatedGroup) error {
	if len(groups) == 0 && !HasYAMLPath(templateDoc, "proxy-groups") {
		return nil
	}

	groupNodes, err := EnsureYAMLSequencePath(templateDoc, "proxy-groups")
	if err != nil {
		return err
	}

	for _, group := range groups {
		if group.IsCountry {
			if existing := FindYAMLSequenceMappingByStringField(groupNodes, "name", group.Name); existing != nil {
				continue
			}
			if err := AppendYAMLSequenceValue(groupNodes, group); err != nil {
				return err
			}
			continue
		}

		existing := FindYAMLSequenceMappingByStringField(groupNodes, "name", group.Name)
		if existing == nil {
			if err := AppendYAMLSequenceValue(groupNodes, group); err != nil {
				return err
			}
			continue
		}

		if findMappingValue(existing, "proxies") == nil {
			continue
		}

		if err := SetYAMLMappingField(existing, "proxies", group.Proxies); err != nil {
			return err
		}
	}

	return nil
}

// mergeTemplateRuleProviders 以字段级 patch 的方式更新/插入 rule-provider，
// 以避免覆盖模板中已有 provider 的未知字段。
func mergeTemplateRuleProviders(templateDoc *yaml.Node, providers map[string]generatedRulePatch) error {
	if len(providers) == 0 && !HasYAMLPath(templateDoc, "rule-providers") {
		return nil
	}

	providerNodes, err := EnsureYAMLMappingPath(templateDoc, "rule-providers")
	if err != nil {
		return err
	}

	for name, provider := range providers {
		existing := findMappingValue(providerNodes, name)
		if existing != nil && existing.Kind == yaml.MappingNode {
			if err := SetYAMLMappingField(existing, "type", provider.Type); err != nil {
				return err
			}
			if err := SetYAMLMappingField(existing, "behavior", provider.Behavior); err != nil {
				return err
			}
			if err := SetYAMLMappingField(existing, "url", provider.Url); err != nil {
				return err
			}
			if err := SetYAMLMappingField(existing, "path", provider.Path); err != nil {
				return err
			}
			if err := SetYAMLMappingField(existing, "interval", provider.Interval); err != nil {
				return err
			}
			if provider.Format != "" {
				if err := SetYAMLMappingField(existing, "format", provider.Format); err != nil {
					return err
				}
			}
			continue
		}

		if err := SetYAMLMappingField(providerNodes, name, provider); err != nil {
			return err
		}
	}

	return nil
}

func FetchSubscriptionUserInfo(url string, userAgent string, retryTimes int) (string, error) {
	client := Request(retryTimes)
	defer client.Close()
	resp, err := client.R().SetHeader("User-Agent", userAgent).Head(url)
	if err != nil {
		logger.Logger.Debug("创建 HEAD 请求失败", zap.Error(err))
		return "", NewNetworkRequestError(url, err)
	}
	defer resp.Body.Close()
	if userInfo := resp.Header().Get("subscription-userinfo"); userInfo != "" {
		return userInfo, nil
	}

	logger.Logger.Debug("subscription-userinfo header not found in response")
	return "", NewNetworkResponseError("subscription-userinfo header not found", nil)
}

// MergeSubAndTemplate 把“模板侧需要参与计算的最小叠加层”和“本项目生成结果”合并。
// 它只处理本项目关心的运行期结构，不负责最终 YAML 输出。
func MergeSubAndTemplate(temp *generatedConfig, sub *generatedConfig, igcg bool) {
	var countryGroupNames []string
	for _, proxyGroup := range sub.ProxyGroup {
		if proxyGroup.IsCountry {
			countryGroupNames = append(
				countryGroupNames, proxyGroup.Name,
			)
		}
	}
	var proxyNames []string
	for _, proxy := range sub.Proxy {
		proxyNames = append(proxyNames, proxy.Name)
	}

	for i := range temp.ProxyGroup {
		if temp.ProxyGroup[i].IsCountry {
			continue
		}
		newProxies := make([]string, 0)
		countryGroupMap := make(map[string]generatedGroup)
		for _, v := range sub.ProxyGroup {
			if v.IsCountry {
				countryGroupMap[v.Name] = v
			}
		}
		for j := range temp.ProxyGroup[i].Proxies {
			reg := regexp.MustCompile("<(.*?)>")
			if reg.Match([]byte(temp.ProxyGroup[i].Proxies[j])) {
				key := reg.FindStringSubmatch(temp.ProxyGroup[i].Proxies[j])[1]
				switch key {
				case "all":
					newProxies = append(newProxies, proxyNames...)
				case "countries":
					if !igcg {
						newProxies = append(newProxies, countryGroupNames...)
					}
				default:
					if !igcg {
						if len(key) == 2 {
							newProxies = append(
								newProxies, countryGroupMap[GetContryName(key)].Proxies...,
							)
						}
					}
				}
			} else {
				newProxies = append(newProxies, temp.ProxyGroup[i].Proxies[j])
			}
		}
		temp.ProxyGroup[i].Proxies = newProxies
	}
	if !igcg {
		temp.ProxyGroup = append(temp.ProxyGroup, sub.ProxyGroup...)
	}
}
