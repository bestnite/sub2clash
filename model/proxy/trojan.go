package proxy

// https://github.com/MetaCubeX/mihomo/blob/Meta/adapter/outbound/trojan.go
type Trojan struct {
	Server            string         `yaml:"server"`
	Port              IntOrString    `yaml:"port"`
	Password          string         `yaml:"password"`
	ALPN              []string       `yaml:"alpn,omitempty"`
	SNI               string         `yaml:"sni,omitempty"`
	SkipCertVerify    bool           `yaml:"skip-cert-verify,omitempty"`
	Fingerprint       string         `yaml:"fingerprint,omitempty"`
	UDP               bool           `yaml:"udp,omitempty"`
	Network           string         `yaml:"network,omitempty"`
	ECHOpts           ECHOptions     `yaml:"ech-opts,omitempty"`
	RealityOpts       RealityOptions `yaml:"reality-opts,omitempty"`
	GrpcOpts          GrpcOptions    `yaml:"grpc-opts,omitempty"`
	WSOpts            WSOptions      `yaml:"ws-opts,omitempty"`
	SSOpts            TrojanSSOption `yaml:"ss-opts,omitempty"`
	ClientFingerprint string         `yaml:"client-fingerprint,omitempty"`
}

type TrojanSSOption struct {
	Enabled  bool   `yaml:"enabled,omitempty"`
	Method   string `yaml:"method,omitempty"`
	Password string `yaml:"password,omitempty"`
}
