package parser

import (
	"fmt"
	"net/url"
	"strings"

	P "github.com/bestnite/sub2clash/model/proxy"
)

type VlessParser struct{}

func (p *VlessParser) SupportClash() bool {
	return false
}

func (p *VlessParser) SupportMeta() bool {
	return true
}

func (p *VlessParser) GetPrefixes() []string {
	return []string{"vless://"}
}

func (p *VlessParser) GetType() string {
	return "vless"
}

func (p *VlessParser) Parse(config ParseConfig, proxy string) (P.Proxy, error) {
	if !hasPrefix(proxy, p.GetPrefixes()) {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidPrefix, proxy)
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
	port, err := ParsePort(portStr)
	if err != nil {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidPort, err.Error())
	}

	query := link.Query()
	uuid := link.User.Username()
	flow, security, alpnStr, sni, insecure, fp, pbk, sid, path, host, serviceName, _type, udp := query.Get("flow"), query.Get("security"), query.Get("alpn"), query.Get("sni"), query.Get("allowInsecure"), query.Get("fp"), query.Get("pbk"), query.Get("sid"), query.Get("path"), query.Get("host"), query.Get("serviceName"), query.Get("type"), query.Get("udp")

	insecureBool := insecure == "1"
	var alpn []string
	if strings.Contains(alpnStr, ",") {
		alpn = strings.Split(alpnStr, ",")
	} else {
		alpn = nil
	}

	remarks := link.Fragment
	if remarks == "" {
		remarks = fmt.Sprintf("%s:%s", server, portStr)
	}
	remarks = strings.TrimSpace(remarks)

	result := P.Vless{
		Server:         server,
		Port:           P.IntOrString(port),
		UUID:           uuid,
		Flow:           flow,
		UDP:            udp == "true" || config.UseUDP,
		SkipCertVerify: insecureBool,
	}

	if len(alpn) > 0 {
		result.ALPN = alpn
	}

	if fp != "" {
		result.ClientFingerprint = fp
	}

	if sni != "" {
		result.ServerName = sni
	}

	if security == "tls" {
		result.TLS = true
	}

	if security == "reality" {
		result.TLS = true
		result.RealityOpts = P.RealityOptions{
			PublicKey: pbk,
			ShortID:   sid,
		}
	}

	if _type == "ws" {
		result.Network = "ws"
		result.WSOpts = P.WSOptions{
			Path: path,
		}
		if host != "" {
			result.WSOpts.Headers = make(map[string]string)
			result.WSOpts.Headers["Host"] = host
		}
	}

	if _type == "grpc" {
		result.Network = "grpc"
		result.GrpcOpts = P.GrpcOptions{
			GrpcServiceName: serviceName,
		}
	}

	if _type == "http" {
		result.HTTPOpts = P.HTTPOptions{}
		result.HTTPOpts.Headers = map[string][]string{}

		result.HTTPOpts.Path = strings.Split(path, ",")

		hosts, err := url.QueryUnescape(host)
		if err != nil {
			return P.Proxy{}, fmt.Errorf("%w: %s", ErrCannotParseParams, err.Error())
		}
		result.Network = "http"
		if hosts != "" {
			result.HTTPOpts.Headers["host"] = strings.Split(host, ",")
		}
	}

	return P.Proxy{
		Type:  p.GetType(),
		Name:  remarks,
		Vless: result,
	}, nil
}

func init() {
	RegisterParser(&VlessParser{})
}
