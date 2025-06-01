package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/nitezs/sub2clash/common"
	"github.com/nitezs/sub2clash/logger"
	"github.com/nitezs/sub2clash/model"
	"github.com/nitezs/sub2clash/parser"
	"github.com/nitezs/sub2clash/validator"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func BuildSub(clashType model.ClashType, query validator.SubValidator, template string) (
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
		templateBytes, err = common.LoadSubscription(template, query.Refresh, query.UserAgent)
		if err != nil {
			logger.Logger.Debug(
				"load template failed", zap.String("template", template), zap.Error(err),
			)
			return nil, errors.New("加载模板失败: " + err.Error())
		}
	} else {
		unescape, err := url.QueryUnescape(template)
		if err != nil {
			return nil, errors.New("加载模板失败: " + err.Error())
		}
		templateBytes, err = common.LoadTemplate(unescape)
		if err != nil {
			logger.Logger.Debug(
				"load template failed", zap.String("template", template), zap.Error(err),
			)
			return nil, errors.New("加载模板失败: " + err.Error())
		}
	}

	err = yaml.Unmarshal(templateBytes, &temp)
	if err != nil {
		logger.Logger.Debug("parse template failed", zap.Error(err))
		return nil, errors.New("解析模板失败: " + err.Error())
	}
	var proxyList []model.Proxy

	for i := range query.Subs {
		data, err := common.LoadSubscription(query.Subs[i], query.Refresh, query.UserAgent)
		subName := ""
		if strings.Contains(query.Subs[i], "#") {
			subName = query.Subs[i][strings.LastIndex(query.Subs[i], "#")+1:]
		}
	
		if err != nil {
			logger.Logger.Debug(
				"load subscription failed", zap.String("url", query.Subs[i]), zap.Error(err),
			)
			if !query.SkipErrors {
				return nil, errors.New("加载订阅失败: " + err.Error())
			}
			// 打印错误并跳过当前订阅
			continue
		}

		// 打印订阅内容
		logger.Logger.Debug("subscription content", 
			zap.String("url", query.Subs[i]),
			zap.String("content", string(data)),
			zap.Int("content_length", len(data)))
	
		err = yaml.Unmarshal(data, &sub)
		var newProxies []model.Proxy
		if err != nil {
			reg, _ := regexp.Compile("(ssr|ss|vmess|trojan|vless|hysteria|hy2|hysteria2)://")
			if reg.Match(data) {
				p := common.ParseProxy(strings.Split(string(data), "\n")...)
				newProxies = p
			} else {
				base64, err := parser.DecodeBase64(string(data))
				if err != nil {
					logger.Logger.Debug(
						"parse subscription failed", zap.String("url", query.Subs[i]),
						zap.String("data", string(data)),
						zap.Error(err),
					)
					if !query.SkipErrors {
						return nil, errors.New("解析订阅失败: " + err.Error())
					}
					continue // 继续下一个订阅
				}
				p := common.ParseProxy(strings.Split(base64, "\n")...)
				newProxies = p
			}
		} else {
			newProxies = sub.Proxies
			// 打印解析后的订阅内容
			logger.Logger.Debug("parsed subscription content",
				zap.String("url", query.Subs[i]),
				zap.Int("proxy_count", len(newProxies)),
				zap.Any("proxies", newProxies))
		}
	
		if subName != "" {
			for i := range newProxies {
				newProxies[i].SubName = subName
			}
		}
		proxyList = append(proxyList, newProxies...)
	}
	

	if len(query.Proxies) != 0 {
		proxyList = append(proxyList, common.ParseProxy(query.Proxies...)...)
	}

	for i := range proxyList {
		if proxyList[i].SubName != "" {
			proxyList[i].Name = strings.TrimSpace(proxyList[i].SubName) + " " + strings.TrimSpace(proxyList[i].Name)
		}
	}

	proxies := make(map[string]*model.Proxy)
	newProxies := make([]model.Proxy, 0, len(proxyList))
	for i := range proxyList {
		key := proxyList[i].Server + strconv.Itoa(proxyList[i].Port) + proxyList[i].Type + proxyList[i].UUID + proxyList[i].Password
		if proxyList[i].Network == "ws" {
			key += proxyList[i].WSOpts.Path + proxyList[i].WSOpts.Headers["Host"]
		}
		if _, exist := proxies[key]; !exist {
			proxies[key] = &proxyList[i]
			newProxies = append(newProxies, proxyList[i])
		}
	}
	proxyList = newProxies

	if strings.TrimSpace(query.Remove) != "" {
		newProxyList := make([]model.Proxy, 0, len(proxyList))
		for i := range proxyList {
			removeReg, err := regexp.Compile(query.Remove)
			if err != nil {
				logger.Logger.Debug("remove regexp compile failed", zap.Error(err))
				return nil, errors.New("remove 参数非法: " + err.Error())
			}

			if removeReg.MatchString(proxyList[i].Name) {
				continue
			}
			newProxyList = append(newProxyList, proxyList[i])
		}
		proxyList = newProxyList
	}

	if len(query.ReplaceKeys) != 0 {

		replaceRegs := make([]*regexp.Regexp, 0, len(query.ReplaceKeys))
		for _, v := range query.ReplaceKeys {
			replaceReg, err := regexp.Compile(v)
			if err != nil {
				logger.Logger.Debug("replace regexp compile failed", zap.Error(err))
				return nil, errors.New("replace 参数非法: " + err.Error())
			}
			replaceRegs = append(replaceRegs, replaceReg)
		}
		for i := range proxyList {

			for j, v := range replaceRegs {
				if v.MatchString(proxyList[i].Name) {
					proxyList[i].Name = v.ReplaceAllString(
						proxyList[i].Name, query.ReplaceTo[j],
					)
				}
			}
		}
	}

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
	common.AddProxy(t, query.AutoTest, query.Lazy, clashType, proxyList...)

	switch query.Sort {
	case "sizeasc":
		sort.Sort(model.ProxyGroupsSortBySize(t.ProxyGroups))
	case "sizedesc":
		sort.Sort(sort.Reverse(model.ProxyGroupsSortBySize(t.ProxyGroups)))
	case "nameasc":
		sort.Sort(model.ProxyGroupsSortByName(t.ProxyGroups))
	case "namedesc":
		sort.Sort(sort.Reverse(model.ProxyGroupsSortByName(t.ProxyGroups)))
	default:
		sort.Sort(model.ProxyGroupsSortByName(t.ProxyGroups))
	}

	MergeSubAndTemplate(temp, t, query.IgnoreCountryGrooup)

	for _, v := range query.Rules {
		if v.Prepend {
			common.PrependRules(temp, v.Rule)
		} else {
			common.AppendRules(temp, v.Rule)
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
			common.PrependRuleProvider(
				temp, v.Name, v.Group, provider,
			)
		} else {
			common.AppenddRuleProvider(
				temp, v.Name, v.Group, provider,
			)
		}
	}
	return temp, nil
}

func fetchSubscriptionUserInfo(url string, userAgent string) (string, error) {
	resp, err := common.Head(url, common.WithUserAgent(userAgent))
	if err != nil {
		logger.Logger.Debug("创建 HEAD 请求失败", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()
	if userInfo := resp.Header.Get("subscription-userinfo"); userInfo != "" {
		return userInfo, nil
	}

	logger.Logger.Debug("目标 URL 未返回 subscription-userinfo 头", zap.Error(err))
	return "", err
}

func MergeSubAndTemplate(temp *model.Subscription, sub *model.Subscription, igcg bool) {

	var countryGroupNames []string
	for _, proxyGroup := range sub.ProxyGroups {
		if proxyGroup.IsCountryGrop {
			countryGroupNames = append(
				countryGroupNames, proxyGroup.Name,
			)
		}
	}
	var proxyNames []string
	for _, proxy := range sub.Proxies {
		proxyNames = append(proxyNames, proxy.Name)
	}

	temp.Proxies = append(temp.Proxies, sub.Proxies...)

	for i := range temp.ProxyGroups {
		if temp.ProxyGroups[i].IsCountryGrop {
			continue
		}
		newProxies := make([]string, 0)
		countryGroupMap := make(map[string]model.ProxyGroup)
		for _, v := range sub.ProxyGroups {
			if v.IsCountryGrop {
				countryGroupMap[v.Name] = v
			}
		}
		for j := range temp.ProxyGroups[i].Proxies {
			reg := regexp.MustCompile("<(.*?)>")
			if reg.Match([]byte(temp.ProxyGroups[i].Proxies[j])) {
				key := reg.FindStringSubmatch(temp.ProxyGroups[i].Proxies[j])[1]
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
								newProxies, countryGroupMap[common.GetContryName(key)].Proxies...,
							)
						}
					}
				}
			} else {
				newProxies = append(newProxies, temp.ProxyGroups[i].Proxies[j])
			}
		}
		temp.ProxyGroups[i].Proxies = newProxies
	}
	if !igcg {
		temp.ProxyGroups = append(temp.ProxyGroups, sub.ProxyGroups...)
	}
}
