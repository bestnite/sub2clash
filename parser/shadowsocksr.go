package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	P "github.com/bestnite/sub2clash/model/proxy"
	"github.com/bestnite/sub2clash/utils"
)

type ShadowsocksRParser struct{}

func (p *ShadowsocksRParser) SupportClash() bool {
	return true
}

func (p *ShadowsocksRParser) SupportMeta() bool {
	return true
}

func (p *ShadowsocksRParser) GetPrefixes() []string {
	return []string{"ssr://"}
}

func (p *ShadowsocksRParser) GetType() string {
	return "ssr"
}

func (p *ShadowsocksRParser) Parse(config ParseConfig, proxy string) (P.Proxy, error) {
	if !hasPrefix(proxy, p.GetPrefixes()) {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidPrefix, proxy)
	}

	for _, prefix := range p.GetPrefixes() {
		if strings.HasPrefix(proxy, prefix) {
			proxy = strings.TrimPrefix(proxy, prefix)
			break
		}
	}

	proxy, err := utils.DecodeBase64(proxy, true)
	if err != nil {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidBase64, err.Error())
	}
	serverInfoAndParams := strings.SplitN(proxy, "/?", 2)
	if len(serverInfoAndParams) != 2 {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidStruct, proxy)
	}
	parts := SplitNRight(serverInfoAndParams[0], ":", 6)
	if len(parts) < 6 {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidStruct, proxy)
	}
	server := parts[0]
	protocol := parts[2]
	method := parts[3]
	obfs := parts[4]
	password, err := utils.DecodeBase64(parts[5], true)
	if err != nil {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidStruct, err.Error())
	}
	port, err := ParsePort(parts[1])
	if err != nil {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidPort, err.Error())
	}

	var obfsParam string
	var protoParam string
	var remarks string
	if len(serverInfoAndParams) == 2 {
		params, err := url.ParseQuery(serverInfoAndParams[1])
		if err != nil {
			return P.Proxy{}, fmt.Errorf("%w: %s", ErrCannotParseParams, err.Error())
		}
		if params.Get("obfsparam") != "" {
			obfsParam, err = utils.DecodeBase64(params.Get("obfsparam"), true)
		}
		if params.Get("protoparam") != "" {
			protoParam, err = utils.DecodeBase64(params.Get("protoparam"), true)
		}
		if params.Get("remarks") != "" {
			remarks, err = utils.DecodeBase64(params.Get("remarks"), true)
		} else {
			remarks = server + ":" + strconv.Itoa(port)
		}
		if err != nil {
			return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidStruct, err.Error())
		}
	}

	result := P.Proxy{
		Type: p.GetType(),
		Name: remarks,
		ShadowSocksR: P.ShadowSocksR{
			Server:        server,
			Port:          port,
			Protocol:      protocol,
			Cipher:        method,
			Obfs:          obfs,
			Password:      password,
			ObfsParam:     obfsParam,
			ProtocolParam: protoParam,
			UDP:           config.UseUDP,
		},
	}
	return result, nil
}

func init() {
	RegisterParser(&ShadowsocksRParser{})
}

func SplitNRight(s, sep string, n int) []string {
	if n <= 0 {
		return nil
	}
	if n == 1 {
		return []string{s}
	}
	parts := strings.Split(s, sep)
	if len(parts) <= n {
		return parts
	}
	result := make([]string, n)
	for i, j := len(parts)-1, 0; i >= 0; i, j = i-1, j+1 {
		if j < n-1 {
			result[n-j-1] = parts[len(parts)-j-1]
		} else {
			result[0] = strings.Join(parts[:i+1], sep)
			break
		}
	}
	return result
}
