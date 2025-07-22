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

func BuildSub(clashType model.ClashType, query model.ConvertConfig, template string, cacheExpire int64, retryTimes int) (
	*model.Subscription, error,
) {
	var temp = &model.Subscription{}
	var sub = &model.Subscription{}
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
			return nil, NewTemplateLoadError(template, err)
		}
	} else {
		unescape, err := url.QueryUnescape(template)
		if err != nil {
			return nil, NewTemplateLoadError(template, err)
		}
		templateBytes, err = LoadTemplate(unescape)
		if err != nil {
			logger.Logger.Debug(
				"load template failed", zap.String("template", template), zap.Error(err),
			)
			return nil, NewTemplateLoadError(unescape, err)
		}
	}

	err = yaml.Unmarshal(templateBytes, &temp)
	if err != nil {
		logger.Logger.Debug("parse template failed", zap.Error(err))
		return nil, NewTemplateParseError(err)
	}
	var proxyList []P.Proxy

	for i := range query.Subs {
		data, err := LoadSubscription(query.Subs[i], query.Refresh, query.UserAgent, cacheExpire, retryTimes)
		subName := ""
		if strings.Contains(query.Subs[i], "#") {
			subName = query.Subs[i][strings.LastIndex(query.Subs[i], "#")+1:]
		}
		if err != nil {
			logger.Logger.Debug(
				"load subscription failed", zap.String("url", query.Subs[i]), zap.Error(err),
			)
			return nil, NewSubscriptionLoadError(query.Subs[i], err)
		}

		err = yaml.Unmarshal(data, &sub)
		var newProxies []P.Proxy
		if err != nil {
			reg, err := regexp.Compile("(" + strings.Join(parser.GetAllPrefixes(), "|") + ")://")
			if err != nil {
				logger.Logger.Debug("compile regex failed", zap.Error(err))
				return nil, NewRegexInvalidError("prefix", err)
			}
			if reg.Match(data) {
				p, err := parser.ParseProxies(parser.ParseConfig{UseUDP: query.UseUDP}, strings.Split(string(data), "\n")...)
				if err != nil {
					return nil, err
				}
				newProxies = p
			} else {
				base64, err := utils.DecodeBase64(string(data), true)
				if err != nil {
					logger.Logger.Debug(
						"parse subscription failed", zap.String("url", query.Subs[i]),
						zap.String("data", string(data)),
						zap.Error(err),
					)
					return nil, NewSubscriptionParseError(err)
				}
				p, err := parser.ParseProxies(parser.ParseConfig{UseUDP: query.UseUDP}, strings.Split(base64, "\n")...)
				if err != nil {
					return nil, err
				}
				newProxies = p
			}
		} else {
			newProxies = sub.Proxy
		}
		if subName != "" {
			for i := range newProxies {
				newProxies[i].SubName = subName
			}
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

	for i := range proxyList {
		if proxyList[i].SubName != "" {
			proxyList[i].Name = strings.TrimSpace(proxyList[i].SubName) + " " + strings.TrimSpace(proxyList[i].Name)
		}
	}

	// 去重
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
	proxyList = newProxies

	// 移除
	if strings.TrimSpace(query.Remove) != "" {
		newProxyList := make([]P.Proxy, 0, len(proxyList))
		for i := range proxyList {
			removeReg, err := regexp.Compile(query.Remove)
			if err != nil {
				logger.Logger.Debug("remove regexp compile failed", zap.Error(err))
				return nil, NewRegexInvalidError("remove", err)
			}

			if removeReg.MatchString(proxyList[i].Name) {
				continue
			}
			newProxyList = append(newProxyList, proxyList[i])
		}
		proxyList = newProxyList
	}

	// 替换
	if len(query.Replace) != 0 {
		for k, v := range query.Replace {
			replaceReg, err := regexp.Compile(k)
			if err != nil {
				logger.Logger.Debug("replace regexp compile failed", zap.Error(err))
				return nil, NewRegexInvalidError("replace", err)
			}
			for i := range proxyList {
				if replaceReg.MatchString(proxyList[i].Name) {
					proxyList[i].Name = replaceReg.ReplaceAllString(
						proxyList[i].Name, v,
					)
				}
			}
		}
	}

	// 重命名有相同名称的节点
	names := make(map[string]int)
	for i := range proxyList {
		if _, exist := names[proxyList[i].Name]; exist {
			names[proxyList[i].Name] = names[proxyList[i].Name] + 1
			proxyList[i].Name = proxyList[i].Name + " " + strconv.Itoa(names[proxyList[i].Name])
		} else {
			names[proxyList[i].Name] = 0
		}
	}

	for i := range proxyList {
		proxyList[i].Name = strings.TrimSpace(proxyList[i].Name)
	}

	var t = &model.Subscription{}
	AddProxy(t, query.AutoTest, query.Lazy, clashType, proxyList...)

	// 排序
	switch query.Sort {
	case "sizeasc":
		sort.Sort(model.ProxyGroupsSortBySize(t.ProxyGroup))
	case "sizedesc":
		sort.Sort(sort.Reverse(model.ProxyGroupsSortBySize(t.ProxyGroup)))
	case "nameasc":
		sort.Sort(model.ProxyGroupsSortByName(t.ProxyGroup))
	case "namedesc":
		sort.Sort(sort.Reverse(model.ProxyGroupsSortByName(t.ProxyGroup)))
	default:
		sort.Sort(model.ProxyGroupsSortByName(t.ProxyGroup))
	}

	MergeSubAndTemplate(temp, t, query.IgnoreCountryGrooup)

	for _, v := range query.Rules {
		if v.Prepend {
			PrependRules(temp, v.Rule)
		} else {
			AppendRules(temp, v.Rule)
		}
	}

	for _, v := range query.RuleProviders {
		hash := sha256.Sum224([]byte(v.Url))
		name := hex.EncodeToString(hash[:])
		provider := model.RuleProvider{
			Type:     "http",
			Behavior: v.Behavior,
			Url:      v.Url,
			Path:     "./" + name + ".yaml",
			Interval: 3600,
		}
		if v.Prepend {
			PrependRuleProvider(
				temp, v.Name, v.Group, provider,
			)
		} else {
			AppenddRuleProvider(
				temp, v.Name, v.Group, provider,
			)
		}
	}
	return temp, nil
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

func MergeSubAndTemplate(temp *model.Subscription, sub *model.Subscription, igcg bool) {
	var countryGroupNames []string
	for _, proxyGroup := range sub.ProxyGroup {
		if proxyGroup.IsCountryGrop {
			countryGroupNames = append(
				countryGroupNames, proxyGroup.Name,
			)
		}
	}
	var proxyNames []string
	for _, proxy := range sub.Proxy {
		proxyNames = append(proxyNames, proxy.Name)
	}

	temp.Proxy = append(temp.Proxy, sub.Proxy...)

	for i := range temp.ProxyGroup {
		if temp.ProxyGroup[i].IsCountryGrop {
			continue
		}
		newProxies := make([]string, 0)
		countryGroupMap := make(map[string]model.ProxyGroup)
		for _, v := range sub.ProxyGroup {
			if v.IsCountryGrop {
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
