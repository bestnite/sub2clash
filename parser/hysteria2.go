package parser

import (
	"fmt"
	"net/url"
	"strings"

	P "github.com/bestnite/sub2clash/model/proxy"
)

type Hysteria2Parser struct{}

func (p *Hysteria2Parser) SupportClash() bool {
	return false
}

func (p *Hysteria2Parser) SupportMeta() bool {
	return true
}

func (p *Hysteria2Parser) GetPrefixes() []string {
	return []string{"hysteria2://", "hy2://"}
}

func (p *Hysteria2Parser) GetType() string {
	return "hysteria2"
}

func (p *Hysteria2Parser) Parse(config ParseConfig, proxy string) (P.Proxy, error) {
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
	obfs, obfsPassword, insecure, sni := query.Get("obfs"), query.Get("obfs-password"), query.Get("insecure"), query.Get("sni")
	insecureBool := insecure == "1"
	remarks := link.Fragment
	if remarks == "" {
		remarks = fmt.Sprintf("%s:%s", server, portStr)
	}
	remarks = strings.TrimSpace(remarks)

	result := P.Proxy{
		Type: p.GetType(),
		Name: remarks,
		Hysteria2: P.Hysteria2{
			Server:         server,
			Port:           P.IntOrString(port),
			Password:       password,
			Obfs:           obfs,
			ObfsPassword:   obfsPassword,
			SNI:            sni,
			SkipCertVerify: insecureBool,
		},
	}
	return result, nil
}

func init() {
	RegisterParser(&Hysteria2Parser{})
}
