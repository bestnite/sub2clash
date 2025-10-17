package proxy

// https://github.com/MetaCubeX/mihomo/blob/Meta/adapter/outbound/hysteria2.go
type Hysteria2 struct {
	Server         string      `yaml:"server"`
	Port           IntOrString `yaml:"port,omitempty"`
	Ports          string      `yaml:"ports,omitempty"`
	HopInterval    int         `yaml:"hop-interval,omitempty"`
	Up             string      `yaml:"up,omitempty"`
	Down           string      `yaml:"down,omitempty"`
	Password       string      `yaml:"password,omitempty"`
	Obfs           string      `yaml:"obfs,omitempty"`
	ObfsPassword   string      `yaml:"obfs-password,omitempty"`
	SNI            string      `yaml:"sni,omitempty"`
	ECHOpts        ECHOptions  `yaml:"ech-opts,omitempty"`
	SkipCertVerify bool        `yaml:"skip-cert-verify,omitempty"`
	Fingerprint    string      `yaml:"fingerprint,omitempty"`
	ALPN           []string    `yaml:"alpn,omitempty"`
	CustomCA       string      `yaml:"ca,omitempty"`
	CustomCAString string      `yaml:"ca-str,omitempty"`
	CWND           int         `yaml:"cwnd,omitempty"`
	UdpMTU         int         `yaml:"udp-mtu,omitempty"`

	// quic-go special config
	InitialStreamReceiveWindow     uint64 `yaml:"initial-stream-receive-window,omitempty"`
	MaxStreamReceiveWindow         uint64 `yaml:"max-stream-receive-window,omitempty"`
	InitialConnectionReceiveWindow uint64 `yaml:"initial-connection-receive-window,omitempty"`
	MaxConnectionReceiveWindow     uint64 `yaml:"max-connection-receive-window,omitempty"`
}
