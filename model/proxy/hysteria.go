package proxy

// https://github.com/MetaCubeX/mihomo/blob/Meta/adapter/outbound/hysteria.go
type Hysteria struct {
	Server              string      `yaml:"server"`
	Port                IntOrString `yaml:"port,omitempty"`
	Ports               string      `yaml:"ports,omitempty"`
	Protocol            string      `yaml:"protocol,omitempty"`
	ObfsProtocol        string      `yaml:"obfs-protocol,omitempty"` // compatible with Stash
	Up                  string      `yaml:"up"`
	UpSpeed             int         `yaml:"up-speed,omitempty"` // compatible with Stash
	Down                string      `yaml:"down"`
	DownSpeed           int         `yaml:"down-speed,omitempty"` // compatible with Stash
	Auth                string      `yaml:"auth,omitempty"`
	AuthString          string      `yaml:"auth-str,omitempty"`
	Obfs                string      `yaml:"obfs,omitempty"`
	SNI                 string      `yaml:"sni,omitempty"`
	ECHOpts             ECHOptions  `yaml:"ech-opts,omitempty"`
	SkipCertVerify      bool        `yaml:"skip-cert-verify,omitempty"`
	Fingerprint         string      `yaml:"fingerprint,omitempty"`
	ALPN                []string    `yaml:"alpn,omitempty"`
	CustomCA            string      `yaml:"ca,omitempty"`
	CustomCAString      string      `yaml:"ca-str,omitempty"`
	ReceiveWindowConn   int         `yaml:"recv-window-conn,omitempty"`
	ReceiveWindow       int         `yaml:"recv-window,omitempty"`
	DisableMTUDiscovery bool        `yaml:"disable-mtu-discovery,omitempty"`
	FastOpen            bool        `yaml:"fast-open,omitempty"`
	HopInterval         int         `yaml:"hop-interval,omitempty"`
}
