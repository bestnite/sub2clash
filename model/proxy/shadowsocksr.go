package proxy

// https://github.com/MetaCubeX/mihomo/blob/Meta/adapter/outbound/shadowsocksr.go
type ShadowSocksR struct {
	Server        string      `yaml:"server"`
	Port          IntOrString `yaml:"port"`
	Password      string      `yaml:"password"`
	Cipher        string      `yaml:"cipher"`
	Obfs          string      `yaml:"obfs"`
	ObfsParam     string      `yaml:"obfs-param,omitempty"`
	Protocol      string      `yaml:"protocol"`
	ProtocolParam string      `yaml:"protocol-param,omitempty"`
	UDP           bool        `yaml:"udp,omitempty"`
}
