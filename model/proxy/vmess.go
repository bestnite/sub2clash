package proxy

// https://github.com/MetaCubeX/mihomo/blob/Meta/adapter/outbound/vmess.go
type Vmess struct {
	Server              string         `yaml:"server"`
	Port                IntOrString    `yaml:"port"`
	UUID                string         `yaml:"uuid"`
	AlterID             IntOrString    `yaml:"alterId"`
	Cipher              string         `yaml:"cipher"`
	UDP                 bool           `yaml:"udp,omitempty"`
	Network             string         `yaml:"network,omitempty"`
	TLS                 bool           `yaml:"tls,omitempty"`
	ALPN                []string       `yaml:"alpn,omitempty"`
	SkipCertVerify      bool           `yaml:"skip-cert-verify,omitempty"`
	Fingerprint         string         `yaml:"fingerprint,omitempty"`
	ServerName          string         `yaml:"servername,omitempty"`
	ECHOpts             ECHOptions     `yaml:"ech-opts,omitempty"`
	RealityOpts         RealityOptions `yaml:"reality-opts,omitempty"`
	HTTPOpts            HTTPOptions    `yaml:"http-opts,omitempty"`
	HTTP2Opts           HTTP2Options   `yaml:"h2-opts,omitempty"`
	GrpcOpts            GrpcOptions    `yaml:"grpc-opts,omitempty"`
	WSOpts              WSOptions      `yaml:"ws-opts,omitempty"`
	PacketAddr          bool           `yaml:"packet-addr,omitempty"`
	XUDP                bool           `yaml:"xudp,omitempty"`
	PacketEncoding      string         `yaml:"packet-encoding,omitempty"`
	GlobalPadding       bool           `yaml:"global-padding,omitempty"`
	AuthenticatedLength bool           `yaml:"authenticated-length,omitempty"`
	ClientFingerprint   string         `yaml:"client-fingerprint,omitempty"`
}
