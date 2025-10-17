package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	P "github.com/bestnite/sub2clash/model/proxy"
)

type HysteriaParser struct{}

func (p *HysteriaParser) SupportClash() bool {
	return false
}

func (p *HysteriaParser) SupportMeta() bool {
	return true
}

func (p *HysteriaParser) GetPrefixes() []string {
	return []string{"hysteria://"}
}

func (p *HysteriaParser) GetType() string {
	return "hysteria"
}

func (p *HysteriaParser) Parse(config ParseConfig, proxy string) (P.Proxy, error) {
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
	if portStr == "" {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidStruct, "missing server port")
	}

	port, err := ParsePort(portStr)
	if err != nil {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidPort, err.Error())
	}

	query := link.Query()

	protocol, auth, auth_str, insecure, upmbps, downmbps, obfs, alpnStr := query.Get("protocol"), query.Get("auth"), query.Get("auth-str"), query.Get("insecure"), query.Get("upmbps"), query.Get("downmbps"), query.Get("obfs"), query.Get("alpn")
	insecureBool, err := strconv.ParseBool(insecure)
	if err != nil {
		insecureBool = false
	}

	var alpn []string
	alpnStr = strings.TrimSpace(alpnStr)
	if alpnStr != "" {
		alpn = strings.Split(alpnStr, ",")
	}

	remarks := link.Fragment
	if remarks == "" {
		remarks = fmt.Sprintf("%s:%s", server, portStr)
	}
	remarks = strings.TrimSpace(remarks)

	result := P.Proxy{
		Type: p.GetType(),
		Name: remarks,
		Hysteria: P.Hysteria{
			Server:         server,
			Port:           P.IntOrString(port),
			Up:             upmbps,
			Down:           downmbps,
			Auth:           auth,
			AuthString:     auth_str,
			Obfs:           obfs,
			SkipCertVerify: insecureBool,
			ALPN:           alpn,
			Protocol:       protocol,
		},
	}
	return result, nil
}

func init() {
	RegisterParser(&HysteriaParser{})
}
