package parser

import (
	"fmt"
	"net/url"
	"strings"

	P "github.com/bestnite/sub2clash/model/proxy"
)

type TrojanParser struct{}

func (p *TrojanParser) SupportClash() bool {
	return true
}

func (p *TrojanParser) SupportMeta() bool {
	return true
}

func (p *TrojanParser) GetPrefixes() []string {
	return []string{"trojan://"}
}

func (p *TrojanParser) GetType() string {
	return "trojan"
}

func (p *TrojanParser) Parse(config ParseConfig, proxy string) (P.Proxy, error) {
	if !hasPrefix(proxy, p.GetPrefixes()) {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidPrefix, proxy)
	}

	link, err := url.Parse(proxy)
	if err != nil {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidStruct, err.Error())
	}

	password := link.User.Username()
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
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidPort, err.Error())
	}

	remarks := link.Fragment
	if remarks == "" {
		remarks = fmt.Sprintf("%s:%s", server, portStr)
	}
	remarks = strings.TrimSpace(remarks)

	query := link.Query()
	network, security, alpnStr, sni, pbk, sid, fp, path, host, serviceName, udp, insecure := query.Get("type"), query.Get("security"), query.Get("alpn"), query.Get("sni"), query.Get("pbk"), query.Get("sid"), query.Get("fp"), query.Get("path"), query.Get("host"), query.Get("serviceName"), query.Get("udp"), query.Get("allowInsecure")

	insecureBool := insecure == "1"
	result := P.Trojan{
		Server:         server,
		Port:           P.IntOrString(port),
		Password:       password,
		Network:        network,
		UDP:            udp == "true" || config.UseUDP,
		SkipCertVerify: insecureBool,
	}

	var alpn []string
	if strings.Contains(alpnStr, ",") {
		alpn = strings.Split(alpnStr, ",")
	} else {
		alpn = nil
	}
	if len(alpn) > 0 {
		result.ALPN = alpn
	}

	if fp != "" {
		result.ClientFingerprint = fp
	}

	if sni != "" {
		result.SNI = sni
	}

	if security == "reality" {
		result.RealityOpts = P.RealityOptions{
			PublicKey: pbk,
			ShortID:   sid,
		}
	}

	if network == "ws" {
		result.Network = "ws"
		result.WSOpts = P.WSOptions{
			Path: path,
			Headers: map[string]string{
				"Host": host,
			},
		}
	}

	if network == "grpc" {
		result.GrpcOpts = P.GrpcOptions{
			GrpcServiceName: serviceName,
		}
	}

	return P.Proxy{
		Type:   p.GetType(),
		Name:   remarks,
		Trojan: result,
	}, nil
}

func init() {
	RegisterParser(&TrojanParser{})
}
