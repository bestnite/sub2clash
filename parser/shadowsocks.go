package parser

import (
	"fmt"
	"net/url"
	"strings"

	P "github.com/bestnite/sub2clash/model/proxy"
	"github.com/bestnite/sub2clash/utils"
)

// ShadowsocksParser Shadowsocks协议解析器
type ShadowsocksParser struct{}

func (p *ShadowsocksParser) SupportClash() bool {
	return true
}

func (p *ShadowsocksParser) SupportMeta() bool {
	return true
}

// GetPrefixes 返回支持的协议前缀
func (p *ShadowsocksParser) GetPrefixes() []string {
	return []string{"ss://"}
}

// GetType 返回协议类型
func (p *ShadowsocksParser) GetType() string {
	return "ss"
}

// Parse 解析Shadowsocks代理
func (p *ShadowsocksParser) Parse(config ParseConfig, proxy string) (P.Proxy, error) {
	if !hasPrefix(proxy, p.GetPrefixes()) {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidPrefix, proxy)
	}

	if !strings.Contains(proxy, "@") {
		s := strings.SplitN(proxy, "#", 2)
		for _, prefix := range p.GetPrefixes() {
			if strings.HasPrefix(s[0], prefix) {
				s[0] = strings.TrimPrefix(s[0], prefix)
				break
			}
		}
		d, err := utils.DecodeBase64(s[0], true)
		if err != nil {
			return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidStruct, err.Error())
		}
		if len(s) == 2 {
			proxy = "ss://" + d + "#" + s[1]
		} else {
			proxy = "ss://" + d
		}
	}
	link, err := url.Parse(proxy)
	if err != nil {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidStruct, err.Error())
	}

	server := link.Hostname()
	if server == "" {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidStruct, "missing server host")
	}

	portStr := link.Port()
	if portStr == "" {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidStruct, "missing server port")
	}
	port, err := ParsePort(portStr)
	if err != nil {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidStruct, err.Error())
	}

	method := link.User.Username()
	password, hasPassword := link.User.Password()

	if !hasPassword && isLikelyBase64(method) {
		decodedStr, err := utils.DecodeBase64(method, true)
		if err == nil {
			methodAndPass := strings.SplitN(decodedStr, ":", 2)
			if len(methodAndPass) == 2 {
				method = methodAndPass[0]
				password = methodAndPass[1]
			} else {
				method = decodedStr
			}
		}
	}
	if password != "" && isLikelyBase64(password) {
		password, err = utils.DecodeBase64(password, true)
		if err != nil {
			return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidStruct, err.Error())
		}
	}

	remarks := link.Fragment
	if remarks == "" {
		remarks = fmt.Sprintf("%s:%s", server, portStr)
	}
	remarks = strings.TrimSpace(remarks)

	result := P.Proxy{
		Type: p.GetType(),
		Name: remarks,
		ShadowSocks: P.ShadowSocks{
			Cipher:   method,
			Password: password,
			Server:   server,
			Port:     port,
			UDP:      config.UseUDP,
		},
	}
	return result, nil
}

// 注册解析器
func init() {
	RegisterParser(&ShadowsocksParser{})
}
