package parser

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	P "github.com/bestnite/sub2clash/model/proxy"
	"github.com/bestnite/sub2clash/utils"
)

type VmessJson struct {
	V    any    `json:"v"`
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Port any    `json:"port"`
	Id   string `json:"id"`
	Aid  any    `json:"aid"`
	Scy  string `json:"scy"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	Tls  string `json:"tls"`
	Sni  string `json:"sni"`
	Alpn string `json:"alpn"`
	Fp   string `json:"fp"`
}

type VmessParser struct{}

func (p *VmessParser) SupportClash() bool {
	return true
}

func (p *VmessParser) SupportMeta() bool {
	return true
}

func (p *VmessParser) GetPrefixes() []string {
	return []string{"vmess://"}
}

func (p *VmessParser) GetType() string {
	return "vmess"
}

func (p *VmessParser) Parse(config ParseConfig, proxy string) (P.Proxy, error) {
	if !hasPrefix(proxy, p.GetPrefixes()) {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidPrefix, proxy)
	}

	for _, prefix := range p.GetPrefixes() {
		if strings.HasPrefix(proxy, prefix) {
			proxy = strings.TrimPrefix(proxy, prefix)
			break
		}
	}
	base64, err := utils.DecodeBase64(proxy, true)
	if err != nil {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidBase64, err.Error())
	}

	var vmess VmessJson
	err = json.Unmarshal([]byte(base64), &vmess)
	if err != nil {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidStruct, err.Error())
	}

	var port int
	switch vmess.Port.(type) {
	case string:
		port, err = ParsePort(vmess.Port.(string))
		if err != nil {
			return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidPort, err.Error())
		}
	case float64:
		port = int(vmess.Port.(float64))
	}

	aid := 0
	switch vmess.Aid.(type) {
	case string:
		aid, err = strconv.Atoi(vmess.Aid.(string))
		if err != nil {
			return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidStruct, err.Error())
		}
	case float64:
		aid = int(vmess.Aid.(float64))
	}

	if vmess.Scy == "" {
		vmess.Scy = "auto"
	}

	name, err := url.QueryUnescape(vmess.Ps)
	if err != nil {
		name = vmess.Ps
	}

	var alpn []string
	if strings.Contains(vmess.Alpn, ",") {
		alpn = strings.Split(vmess.Alpn, ",")
	} else {
		alpn = nil
	}

	result := P.Vmess{
		Server:  vmess.Add,
		Port:    P.IntOrString(port),
		UUID:    vmess.Id,
		AlterID: P.IntOrString(aid),
		Cipher:  vmess.Scy,
		UDP:     config.UseUDP,
	}

	if len(alpn) > 0 {
		result.ALPN = alpn
	}

	if vmess.Fp != "" {
		result.ClientFingerprint = vmess.Fp
	}

	if vmess.Sni != "" {
		result.ServerName = vmess.Sni
	}

	if vmess.Tls == "tls" {
		result.TLS = true
	}

	if vmess.Net == "ws" {
		if vmess.Path == "" {
			vmess.Path = "/"
		}
		if vmess.Host == "" {
			vmess.Host = vmess.Add
		}
		result.Network = "ws"
		result.WSOpts = P.WSOptions{
			Path: vmess.Path,
			Headers: map[string]string{
				"Host": vmess.Host,
			},
		}
	}

	if vmess.Net == "grpc" {
		result.GrpcOpts = P.GrpcOptions{
			GrpcServiceName: vmess.Path,
		}
		result.Network = "grpc"
	}

	if vmess.Net == "h2" {
		result.HTTP2Opts = P.HTTP2Options{
			Host: strings.Split(vmess.Host, ","),
			Path: vmess.Path,
		}
		result.Network = "h2"
	}

	return P.Proxy{
		Type:  p.GetType(),
		Name:  name,
		Vmess: result,
	}, nil
}

func init() {
	RegisterParser(&VmessParser{})
}
