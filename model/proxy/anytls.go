package proxy

// https://github.com/MetaCubeX/mihomo/blob/Meta/adapter/outbound/anytls.go
type Anytls struct {
	Server                   string      `yaml:"server"`
	Port                     IntOrString `yaml:"port"`
	Password                 string      `yaml:"password"`
	ALPN                     []string    `yaml:"alpn,omitempty"`
	SNI                      string      `yaml:"sni,omitempty"`
	ECHOpts                  ECHOptions  `yaml:"ech-opts,omitempty"`
	ClientFingerprint        string      `yaml:"client-fingerprint,omitempty"`
	SkipCertVerify           bool        `yaml:"skip-cert-verify,omitempty"`
	Fingerprint              string      `yaml:"fingerprint,omitempty"`
	UDP                      bool        `yaml:"udp,omitempty"`
	IdleSessionCheckInterval int         `yaml:"idle-session-check-interval,omitempty"`
	IdleSessionTimeout       int         `yaml:"idle-session-timeout,omitempty"`
	MinIdleSession           int         `yaml:"min-idle-session,omitempty"`
}
