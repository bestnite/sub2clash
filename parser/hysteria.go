package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/nitezs/sub2clash/constant"
	"github.com/nitezs/sub2clash/model"
)

func ParseHysteria(proxy string) (model.Proxy, error) {
	if !strings.HasPrefix(proxy, constant.HysteriaPrefix) {
		return model.Proxy{}, &ParseError{Type: ErrInvalidPrefix, Raw: proxy}
	}

	link, err := url.Parse(proxy)
	if err != nil {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidStruct,
			Message: "url parse error",
			Raw:     proxy,
		}
	}
	server := link.Hostname()
	if server == "" {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidStruct,
			Message: "missing server host",
			Raw:     proxy,
		}
	}

	portStr := link.Port()
	if portStr == "" {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidStruct,
			Message: "missing server port",
			Raw:     proxy,
		}
	}

	port, err := ParsePort(portStr)
	if err != nil {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidPort,
			Message: err.Error(),
			Raw:     proxy,
		}
	}

	query := link.Query()

	protocol, authStr, insecure, upmbps, downmbps, obfs, alpnStr := query.Get("protocol"), query.Get("auth_str"), query.Get("insecure"), query.Get("upmbps"), query.Get("downmbps"), query.Get("obfs"), query.Get("alpn")
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

	result := model.Proxy{
		Type:           "hysteria",
		Name:           remarks,
		Server:         server,
		Port:           port,
		Up:             upmbps,
		Down:           downmbps,
		AuthString:     authStr,
		Obfs:           obfs,
		SkipCertVerify: insecureBool,
		Alpn:           alpn,
		Protocol:       protocol,
		AllowInsecure:  insecureBool,
	}
	return result, nil
}
