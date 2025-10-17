package proxy

// https://github.com/MetaCubeX/mihomo/blob/Meta/adapter/outbound/socks5.go
type Socks struct {
	Server         string      `yaml:"server"`
	Port           IntOrString `yaml:"port"`
	UserName       string      `yaml:"username,omitempty"`
	Password       string      `yaml:"password,omitempty"`
	TLS            bool        `yaml:"tls,omitempty"`
	UDP            bool        `yaml:"udp,omitempty"`
	SkipCertVerify bool        `yaml:"skip-cert-verify,omitempty"`
	Fingerprint    string      `yaml:"fingerprint,omitempty"`
}
