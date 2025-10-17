package proxy

// https://github.com/MetaCubeX/mihomo/blob/Meta/adapter/outbound/vless.go
type Vless struct {
	Server            string            `yaml:"server"`
	Port              IntOrString       `yaml:"port"`
	UUID              string            `yaml:"uuid"`
	Flow              string            `yaml:"flow,omitempty"`
	TLS               bool              `yaml:"tls,omitempty"`
	ALPN              []string          `yaml:"alpn,omitempty"`
	UDP               bool              `yaml:"udp,omitempty"`
	PacketAddr        bool              `yaml:"packet-addr,omitempty"`
	XUDP              bool              `yaml:"xudp,omitempty"`
	PacketEncoding    string            `yaml:"packet-encoding,omitempty"`
	Network           string            `yaml:"network,omitempty"`
	ECHOpts           ECHOptions        `yaml:"ech-opts,omitempty"`
	RealityOpts       RealityOptions    `yaml:"reality-opts,omitempty"`
	HTTPOpts          HTTPOptions       `yaml:"http-opts,omitempty"`
	HTTP2Opts         HTTP2Options      `yaml:"h2-opts,omitempty"`
	GrpcOpts          GrpcOptions       `yaml:"grpc-opts,omitempty"`
	WSOpts            WSOptions         `yaml:"ws-opts,omitempty"`
	WSPath            string            `yaml:"ws-path,omitempty"`
	WSHeaders         map[string]string `yaml:"ws-headers,omitempty"`
	SkipCertVerify    bool              `yaml:"skip-cert-verify,omitempty"`
	Fingerprint       string            `yaml:"fingerprint,omitempty"`
	ServerName        string            `yaml:"servername,omitempty"`
	ClientFingerprint string            `yaml:"client-fingerprint,omitempty"`
}
