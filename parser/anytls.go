package parser

import (
	"fmt"
	"net/url"
	"strings"

	P "github.com/bestnite/sub2clash/model/proxy"
)

type AnytlsParser struct{}

func (p *AnytlsParser) SupportClash() bool {
	return false
}

func (p *AnytlsParser) SupportMeta() bool {
	return true
}

func (p *AnytlsParser) GetPrefixes() []string {
	return []string{"anytls://"}
}

func (p *AnytlsParser) GetType() string {
	return "anytls"
}

func (p *AnytlsParser) Parse(config ParseConfig, proxy string) (P.Proxy, error) {
	if !hasPrefix(proxy, p.GetPrefixes()) {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidPrefix, proxy)
	}

	link, err := url.Parse(proxy)
	if err != nil {
		return P.Proxy{}, fmt.Errorf("%w: %s", ErrInvalidStruct, err.Error())
	}

	username := link.User.Username()
	password, exist := link.User.Password()
	if !exist {
		password = username
	}

	query := link.Query()
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
	insecure, sni := query.Get("insecure"), query.Get("sni")
	insecureBool := insecure == "1"
	remarks := link.Fragment
	if remarks == "" {
		remarks = fmt.Sprintf("%s:%s", server, portStr)
	}
	remarks = strings.TrimSpace(remarks)

	result := P.Proxy{
		Type: p.GetType(),
		Name: remarks,
		Anytls: P.Anytls{
			Server:         server,
			Port:           P.IntOrString(port),
			Password:       password,
			SNI:            sni,
			SkipCertVerify: insecureBool,
			UDP:            config.UseUDP,
		},
	}
	return result, nil
}

func init() {
	RegisterParser(&AnytlsParser{})
}
