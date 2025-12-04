package proxy

type Tuic struct {
	Server                string   `proxy:"server"`
	Port                  int      `proxy:"port"`
	Token                 string   `proxy:"token,omitempty"`
	UUID                  string   `proxy:"uuid,omitempty"`
	Password              string   `proxy:"password,omitempty"`
	Ip                    string   `proxy:"ip,omitempty"`
	HeartbeatInterval     int      `proxy:"heartbeat-interval,omitempty"`
	ALPN                  []string `proxy:"alpn,omitempty"`
	ReduceRtt             bool     `proxy:"reduce-rtt,omitempty"`
	RequestTimeout        int      `proxy:"request-timeout,omitempty"`
	UdpRelayMode          string   `proxy:"udp-relay-mode,omitempty"`
	CongestionController  string   `proxy:"congestion-controller,omitempty"`
	DisableSni            bool     `proxy:"disable-sni,omitempty"`
	MaxUdpRelayPacketSize int      `proxy:"max-udp-relay-packet-size,omitempty"`

	FastOpen             bool       `proxy:"fast-open,omitempty"`
	MaxOpenStreams       int        `proxy:"max-open-streams,omitempty"`
	CWND                 int        `proxy:"cwnd,omitempty"`
	SkipCertVerify       bool       `proxy:"skip-cert-verify,omitempty"`
	Fingerprint          string     `proxy:"fingerprint,omitempty"`
	Certificate          string     `proxy:"certificate,omitempty"`
	PrivateKey           string     `proxy:"private-key,omitempty"`
	ReceiveWindowConn    int        `proxy:"recv-window-conn,omitempty"`
	ReceiveWindow        int        `proxy:"recv-window,omitempty"`
	DisableMTUDiscovery  bool       `proxy:"disable-mtu-discovery,omitempty"`
	MaxDatagramFrameSize int        `proxy:"max-datagram-frame-size,omitempty"`
	SNI                  string     `proxy:"sni,omitempty"`
	ECHOpts              ECHOptions `proxy:"ech-opts,omitempty"`

	UDPOverStream        bool `proxy:"udp-over-stream,omitempty"`
	UDPOverStreamVersion int  `proxy:"udp-over-stream-version,omitempty"`
}
