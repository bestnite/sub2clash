package proxy

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type HTTPOptions struct {
	Method  string              `yaml:"method,omitempty"`
	Path    []string            `yaml:"path,omitempty"`
	Headers map[string][]string `yaml:"headers,omitempty"`
}

type HTTP2Options struct {
	Host []string `yaml:"host,omitempty"`
	Path string   `yaml:"path,omitempty"`
}

type GrpcOptions struct {
	GrpcServiceName string `yaml:"grpc-service-name,omitempty"`
}

type RealityOptions struct {
	PublicKey string `yaml:"public-key"`
	ShortID   string `yaml:"short-id,omitempty"`
}

type WSOptions struct {
	Path                string            `yaml:"path,omitempty"`
	Headers             map[string]string `yaml:"headers,omitempty"`
	MaxEarlyData        int               `yaml:"max-early-data,omitempty"`
	EarlyDataHeaderName string            `yaml:"early-data-header-name,omitempty"`
}

type SmuxStruct struct {
	Enabled bool `yaml:"enable"`
}

type WireGuardPeerOption struct {
	Server       string   `yaml:"server"`
	Port         int      `yaml:"port"`
	PublicKey    string   `yaml:"public-key,omitempty"`
	PreSharedKey string   `yaml:"pre-shared-key,omitempty"`
	Reserved     []uint8  `yaml:"reserved,omitempty"`
	AllowedIPs   []string `yaml:"allowed-ips,omitempty"`
}

type ECHOptions struct {
	Enable bool   `yaml:"enable,omitempty" obfs:"enable,omitempty"`
	Config string `yaml:"config,omitempty" obfs:"config,omitempty"`
}

type Proxy struct {
	Type    string
	Name    string
	SubName string `yaml:"-"`
	Anytls
	Hysteria
	Hysteria2
	ShadowSocks
	ShadowSocksR
	Trojan
	Vless
	Vmess
	Socks
}

func (p Proxy) MarshalYAML() (any, error) {
	switch p.Type {
	case "anytls":
		return struct {
			Type   string `yaml:"type"`
			Name   string `yaml:"name"`
			Anytls `yaml:",inline"`
		}{
			Type:   p.Type,
			Name:   p.Name,
			Anytls: p.Anytls,
		}, nil
	case "hysteria":
		return struct {
			Type     string `yaml:"type"`
			Name     string `yaml:"name"`
			Hysteria `yaml:",inline"`
		}{
			Type:     p.Type,
			Name:     p.Name,
			Hysteria: p.Hysteria,
		}, nil
	case "hysteria2":
		return struct {
			Type      string `yaml:"type"`
			Name      string `yaml:"name"`
			Hysteria2 `yaml:",inline"`
		}{
			Type:      p.Type,
			Name:      p.Name,
			Hysteria2: p.Hysteria2,
		}, nil
	case "ss":
		return struct {
			Type        string `yaml:"type"`
			Name        string `yaml:"name"`
			ShadowSocks `yaml:",inline"`
		}{
			Type:        p.Type,
			Name:        p.Name,
			ShadowSocks: p.ShadowSocks,
		}, nil
	case "ssr":
		return struct {
			Type         string `yaml:"type"`
			Name         string `yaml:"name"`
			ShadowSocksR `yaml:",inline"`
		}{
			Type:         p.Type,
			Name:         p.Name,
			ShadowSocksR: p.ShadowSocksR,
		}, nil
	case "trojan":
		return struct {
			Type   string `yaml:"type"`
			Name   string `yaml:"name"`
			Trojan `yaml:",inline"`
		}{
			Type:   p.Type,
			Name:   p.Name,
			Trojan: p.Trojan,
		}, nil
	case "vless":
		return struct {
			Type  string `yaml:"type"`
			Name  string `yaml:"name"`
			Vless `yaml:",inline"`
		}{
			Type:  p.Type,
			Name:  p.Name,
			Vless: p.Vless,
		}, nil
	case "vmess":
		return struct {
			Type  string `yaml:"type"`
			Name  string `yaml:"name"`
			Vmess `yaml:",inline"`
		}{
			Type:  p.Type,
			Name:  p.Name,
			Vmess: p.Vmess,
		}, nil
	case "socks5":
		return struct {
			Type  string `yaml:"type"`
			Name  string `yaml:"name"`
			Socks `yaml:",inline"`
		}{
			Type:  p.Type,
			Name:  p.Name,
			Socks: p.Socks,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported proxy type: %s", p.Type)
	}
}

func (p *Proxy) UnmarshalYAML(node *yaml.Node) error {
	var temp struct {
		Type string `yaml:"type"`
		Name string `yaml:"name"`
	}

	if err := node.Decode(&temp); err != nil {
		return err
	}

	p.Type = temp.Type
	p.Name = temp.Name

	switch temp.Type {
	case "anytls":
		var data struct {
			Type   string `yaml:"type"`
			Name   string `yaml:"name"`
			Anytls `yaml:",inline"`
		}
		if err := node.Decode(&data); err != nil {
			return err
		}
		p.Anytls = data.Anytls

	case "hysteria":
		var data struct {
			Type     string `yaml:"type"`
			Name     string `yaml:"name"`
			Hysteria `yaml:",inline"`
		}
		if err := node.Decode(&data); err != nil {
			return err
		}
		p.Hysteria = data.Hysteria

	case "hysteria2":
		var data struct {
			Type      string `yaml:"type"`
			Name      string `yaml:"name"`
			Hysteria2 `yaml:",inline"`
		}
		if err := node.Decode(&data); err != nil {
			return err
		}
		p.Hysteria2 = data.Hysteria2

	case "ss":
		var data struct {
			Type        string `yaml:"type"`
			Name        string `yaml:"name"`
			ShadowSocks `yaml:",inline"`
		}
		if err := node.Decode(&data); err != nil {
			return err
		}
		p.ShadowSocks = data.ShadowSocks

	case "ssr":
		var data struct {
			Type         string `yaml:"type"`
			Name         string `yaml:"name"`
			ShadowSocksR `yaml:",inline"`
		}
		if err := node.Decode(&data); err != nil {
			return err
		}
		p.ShadowSocksR = data.ShadowSocksR

	case "trojan":
		var data struct {
			Type   string `yaml:"type"`
			Name   string `yaml:"name"`
			Trojan `yaml:",inline"`
		}
		if err := node.Decode(&data); err != nil {
			return err
		}
		p.Trojan = data.Trojan

	case "vless":
		var data struct {
			Type  string `yaml:"type"`
			Name  string `yaml:"name"`
			Vless `yaml:",inline"`
		}
		if err := node.Decode(&data); err != nil {
			return err
		}
		p.Vless = data.Vless

	case "vmess":
		var data struct {
			Type  string `yaml:"type"`
			Name  string `yaml:"name"`
			Vmess `yaml:",inline"`
		}
		if err := node.Decode(&data); err != nil {
			return err
		}
		p.Vmess = data.Vmess

	case "socks5":
		var data struct {
			Type  string `yaml:"type"`
			Name  string `yaml:"name"`
			Socks `yaml:",inline"`
		}
		if err := node.Decode(&data); err != nil {
			return err
		}
		p.Socks = data.Socks

	default:
		return fmt.Errorf("unsupported proxy type: %s", temp.Type)
	}

	return nil
}
