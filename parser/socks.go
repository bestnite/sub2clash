package parser

import (
	"fmt"
	"net/url"
	"strings"

	P "github.com/bestnite/sub2clash/model/proxy"
)

type SocksParser struct{}

func (p *SocksParser) SupportClash() bool {
	return true
}
func (p *SocksParser) SupportMeta() bool {
	return true
}

func (p *SocksParser) GetPrefixes() []string {
	return []string{"socks://", "socks5://"}
}

func (p *SocksParser) GetType() string {
	return "socks5"
}

func (p *SocksParser) Parse(config ParseConfig, proxy string) (P.Proxy, error) {
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

	remarks := link.Fragment
	if remarks == "" {
		remarks = fmt.Sprintf("%s:%s", server, portStr)
	}
	remarks = strings.TrimSpace(remarks)

	var username, password string

	username = link.User.Username()
	password, hasPassword := link.User.Password()

	if !hasPassword && isLikelyBase64(username) {
		decodedStr, err := DecodeBase64(username)
		if err == nil {
			usernameAndPassword := strings.SplitN(decodedStr, ":", 2)
			if len(usernameAndPassword) == 2 {
				username = usernameAndPassword[0]
				password = usernameAndPassword[1]
			} else {
				username = decodedStr
			}
		}
	}

	tls, udp := link.Query().Get("tls"), link.Query().Get("udp")

	return P.Proxy{
		Type: p.GetType(),
		Name: remarks,
		Socks: P.Socks{
			Server:   server,
			Port:     port,
			UserName: username,
			Password: password,
			TLS:      tls == "true",
			UDP:      udp == "true",
		},
	}, nil
}

func init() {
	RegisterParser(&SocksParser{})
}
