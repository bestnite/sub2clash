package parser

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"

	P "github.com/bestnite/sub2clash/model/proxy"
	"github.com/bestnite/sub2clash/utils"
)

func hasPrefix(proxy string, prefixes []string) bool {
	hasPrefix := false
	for _, prefix := range prefixes {
		if strings.HasPrefix(proxy, prefix) {
			hasPrefix = true
			break
		}
	}
	return hasPrefix
}

func ParsePort(portStr string) (int, error) {
	port, err := strconv.Atoi(portStr)

	if err != nil {
		return 0, err
	}
	if port < 1 || port > 65535 {
		return 0, errors.New("invaild port range")
	}
	return port, nil
}

// isLikelyBase64 不严格判断是否是合法的 Base64, 很多分享链接不符合 Base64 规范
func isLikelyBase64(s string) bool {
	if strings.TrimSpace(s) == "" {
		return false
	}

	if !strings.Contains(strings.TrimSuffix(s, "="), "=") {
		s = strings.TrimSuffix(s, "=")
		chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
		for _, c := range s {
			if !strings.ContainsRune(chars, c) {
				return false
			}
		}
	}

	decoded, err := utils.DecodeBase64(s, true)
	if err != nil {
		return false
	}
	if !utf8.ValidString(decoded) {
		return false
	}

	return true
}

type ParseConfig struct {
	UseUDP bool
}

func ParseProxies(config ParseConfig, proxies ...string) ([]P.Proxy, error) {
	var result []P.Proxy
	for _, proxy := range proxies {
		if proxy != "" {
			var proxyItem P.Proxy
			var err error

			proxyItem, err = ParseProxyWithRegistry(config, proxy)
			if err != nil {
				return nil, err
			}
			result = append(result, proxyItem)
		}
	}
	return result, nil
}
