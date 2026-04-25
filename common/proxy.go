package common

import (
	"strings"

	"github.com/bestnite/sub2clash/model"
	"github.com/bestnite/sub2clash/model/proxy"
)

func GetContryName(countryKey string) string {

	countryMaps := []map[string]string{
		model.CountryFlag,
		model.CountryChineseName,
		model.CountryISO,
		model.CountryEnglishName,
	}

	for i, countryMap := range countryMaps {
		if i == 2 {

			splitChars := []string{"-", "_", " "}
			key := make([]string, 0)
			for _, splitChar := range splitChars {
				slic := strings.Split(countryKey, splitChar)
				for _, v := range slic {
					if len(v) == 2 {
						key = append(key, v)
					}
				}
			}

			for _, v := range key {

				if country, ok := countryMap[strings.ToUpper(v)]; ok {
					return country
				}
			}
		}
		for k, v := range countryMap {
			if strings.Contains(countryKey, k) {
				return v
			}
		}
	}
	return "其他地区"
}

func AddProxy(
	sub *generatedConfig, autotest bool,
	lazy bool, clashType model.ClashType, proxies ...proxy.Proxy,
) {
	proxyTypes := model.GetSupportProxyTypes(clashType)

	for _, proxy := range proxies {
		if !proxyTypes[proxy.Type] {
			continue
		}
		sub.Proxy = append(sub.Proxy, proxy)
		haveProxyGroup := false
		countryName := GetContryName(proxy.Name)
		for i := range sub.ProxyGroup {
			group := &sub.ProxyGroup[i]
			if group.Name == countryName {
				group.Proxies = append(group.Proxies, proxy.Name)
				group.Size++
				haveProxyGroup = true
			}
		}
		if !haveProxyGroup {
			var newGroup generatedGroup
			if !autotest {
				newGroup = generatedGroup{
					Name:      countryName,
					Type:      "select",
					Proxies:   []string{proxy.Name},
					IsCountry: true,
					Size:      1,
				}
			} else {
				newGroup = generatedGroup{
					Name:      countryName,
					Type:      "url-test",
					Proxies:   []string{proxy.Name},
					IsCountry: true,
					Url:       "http://www.gstatic.com/generate_204",
					Interval:  300,
					Tolerance: 50,
					Lazy:      lazy,
					Size:      1,
				}
			}
			sub.ProxyGroup = append(sub.ProxyGroup, newGroup)
		}
	}
}
