package model

import (
	"net/netip"

	"github.com/bestnite/sub2clash/model/proxy"
	C "github.com/metacubex/mihomo/config"
	CC "github.com/metacubex/mihomo/constant"
	LC "github.com/metacubex/mihomo/listener/config"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type NodeList struct {
	Proxy []proxy.Proxy `yaml:"proxies,omitempty" json:"proxies"`
}

// https://github.com/MetaCubeX/mihomo/blob/Meta/config/config.go RawConfig
type Subscription struct {
	Port                    proxy.IntOrString `yaml:"port,omitempty" json:"port"`
	SocksPort               proxy.IntOrString `yaml:"socks-port,omitempty" json:"socks-port"`
	RedirPort               proxy.IntOrString `yaml:"redir-port,omitempty" json:"redir-port"`
	TProxyPort              proxy.IntOrString `yaml:"tproxy-port,omitempty" json:"tproxy-port"`
	MixedPort               proxy.IntOrString `yaml:"mixed-port,omitempty" json:"mixed-port"`
	ShadowSocksConfig       string            `yaml:"ss-config,omitempty" json:"ss-config"`
	VmessConfig             string            `yaml:"vmess-config,omitempty" json:"vmess-config"`
	InboundTfo              bool              `yaml:"inbound-tfo,omitempty" json:"inbound-tfo"`
	InboundMPTCP            bool              `yaml:"inbound-mptcp,omitempty" json:"inbound-mptcp"`
	Authentication          []string          `yaml:"authentication,omitempty" json:"authentication"`
	SkipAuthPrefixes        []netip.Prefix    `yaml:"skip-auth-prefixes,omitempty" json:"skip-auth-prefixes"`
	LanAllowedIPs           []netip.Prefix    `yaml:"lan-allowed-ips,omitempty" json:"lan-allowed-ips"`
	LanDisAllowedIPs        []netip.Prefix    `yaml:"lan-disallowed-ips,omitempty" json:"lan-disallowed-ips"`
	AllowLan                bool              `yaml:"allow-lan,omitempty" json:"allow-lan"`
	BindAddress             string            `yaml:"bind-address,omitempty" json:"bind-address"`
	Mode                    string            `yaml:"mode,omitempty" json:"mode"`
	UnifiedDelay            bool              `yaml:"unified-delay,omitempty" json:"unified-delay"`
	LogLevel                string            `yaml:"log-level,omitempty" json:"log-level"`
	IPv6                    bool              `yaml:"ipv6,omitempty" json:"ipv6"`
	ExternalController      string            `yaml:"external-controller,omitempty" json:"external-controller"`
	ExternalControllerPipe  string            `yaml:"external-controller-pipe,omitempty" json:"external-controller-pipe"`
	ExternalControllerUnix  string            `yaml:"external-controller-unix,omitempty" json:"external-controller-unix"`
	ExternalControllerTLS   string            `yaml:"external-controller-tls,omitempty" json:"external-controller-tls"`
	ExternalControllerCors  C.RawCors         `yaml:"external-controller-cors,omitempty" json:"external-controller-cors"`
	ExternalUI              string            `yaml:"external-ui,omitempty" json:"external-ui"`
	ExternalUIURL           string            `yaml:"external-ui-url,omitempty" json:"external-ui-url"`
	ExternalUIName          string            `yaml:"external-ui-name,omitempty" json:"external-ui-name"`
	ExternalDohServer       string            `yaml:"external-doh-server,omitempty" json:"external-doh-server"`
	Secret                  string            `yaml:"secret,omitempty" json:"secret"`
	Interface               string            `yaml:"interface-name,omitempty" json:"interface-name"`
	RoutingMark             int               `yaml:"routing-mark,omitempty" json:"routing-mark"`
	Tunnels                 []LC.Tunnel       `yaml:"tunnels,omitempty" json:"tunnels"`
	GeoAutoUpdate           bool              `yaml:"geo-auto-update,omitempty" json:"geo-auto-update"`
	GeoUpdateInterval       int               `yaml:"geo-update-interval,omitempty" json:"geo-update-interval"`
	GeodataMode             bool              `yaml:"geodata-mode,omitempty" json:"geodata-mode"`
	GeodataLoader           string            `yaml:"geodata-loader,omitempty" json:"geodata-loader"`
	GeositeMatcher          string            `yaml:"geosite-matcher,omitempty" json:"geosite-matcher"`
	TCPConcurrent           bool              `yaml:"tcp-concurrent,omitempty" json:"tcp-concurrent"`
	FindProcessMode         string            `yaml:"find-process-mode,omitempty" json:"find-process-mode"`
	GlobalClientFingerprint string            `yaml:"global-client-fingerprint,omitempty" json:"global-client-fingerprint"`
	GlobalUA                string            `yaml:"global-ua,omitempty" json:"global-ua"`
	ETagSupport             bool              `yaml:"etag-support,omitempty" json:"etag-support"`
	KeepAliveIdle           int               `yaml:"keep-alive-idle,omitempty" json:"keep-alive-idle"`
	KeepAliveInterval       int               `yaml:"keep-alive-interval,omitempty" json:"keep-alive-interval"`
	DisableKeepAlive        bool              `yaml:"disable-keep-alive,omitempty" json:"disable-keep-alive"`

	ProxyProvider map[string]map[string]any `yaml:"proxy-providers,omitempty" json:"proxy-providers"`
	RuleProvider  map[string]RuleProvider   `yaml:"rule-providers,omitempty" json:"rule-providers"`
	Proxy         []proxy.Proxy             `yaml:"proxies,omitempty" json:"proxies"`
	ProxyGroup    []ProxyGroup              `yaml:"proxy-groups,omitempty" json:"proxy-groups"`
	Rule          []string                  `yaml:"rules,omitempty" json:"rule"`
	SubRules      map[string][]string       `yaml:"sub-rules,omitempty" json:"sub-rules"`
	Listeners     []map[string]any          `yaml:"listeners,omitempty" json:"listeners"`
	Hosts         map[string]any            `yaml:"hosts,omitempty" json:"hosts"`
	DNS           RawDNS                    `yaml:"dns,omitempty" json:"dns"`
	NTP           RawNTP                    `yaml:"ntp,omitempty" json:"ntp"`
	Tun           RawTun                    `yaml:"tun,omitempty" json:"tun"`
	TuicServer    RawTuicServer             `yaml:"tuic-server,omitempty" json:"tuic-server"`
	IPTables      RawIPTables               `yaml:"iptables,omitempty" json:"iptables"`
	Experimental  RawExperimental           `yaml:"experimental,omitempty" json:"experimental"`
	Profile       RawProfile                `yaml:"profile,omitempty" json:"profile"`
	GeoXUrl       RawGeoXUrl                `yaml:"geox-url,omitempty" json:"geox-url"`
	Sniffer       RawSniffer                `yaml:"sniffer,omitempty" json:"sniffer"`
	TLS           RawTLS                    `yaml:"tls,omitempty" json:"tls"`

	ClashForAndroid C.RawClashForAndroid `yaml:"clash-for-android,omitempty" json:"clash-for-android"`
}

type RawDNS struct {
	Enable                       bool                                `yaml:"enable,omitempty" json:"enable"`
	PreferH3                     bool                                `yaml:"prefer-h3,omitempty" json:"prefer-h3"`
	IPv6                         bool                                `yaml:"ipv6,omitempty" json:"ipv6"`
	IPv6Timeout                  uint                                `yaml:"ipv6-timeout,omitempty" json:"ipv6-timeout"`
	UseHosts                     bool                                `yaml:"use-hosts,omitempty" json:"use-hosts"`
	UseSystemHosts               bool                                `yaml:"use-system-hosts,omitempty" json:"use-system-hosts"`
	RespectRules                 bool                                `yaml:"respect-rules,omitempty" json:"respect-rules"`
	NameServer                   []string                            `yaml:"nameserver,omitempty" json:"nameserver"`
	Fallback                     []string                            `yaml:"fallback,omitempty" json:"fallback"`
	FallbackFilter               C.RawFallbackFilter                 `yaml:"fallback-filter,omitempty" json:"fallback-filter"`
	Listen                       string                              `yaml:"listen,omitempty" json:"listen"`
	EnhancedMode                 CC.DNSMode                          `yaml:"enhanced-mode,omitempty" json:"enhanced-mode"`
	FakeIPRange                  string                              `yaml:"fake-ip-range,omitempty" json:"fake-ip-range"`
	FakeIPFilter                 []string                            `yaml:"fake-ip-filter,omitempty" json:"fake-ip-filter"`
	FakeIPFilterMode             CC.FilterMode                       `yaml:"fake-ip-filter-mode,omitempty" json:"fake-ip-filter-mode"`
	DefaultNameserver            []string                            `yaml:"default-nameserver,omitempty" json:"default-nameserver"`
	CacheAlgorithm               string                              `yaml:"cache-algorithm,omitempty" json:"cache-algorithm"`
	NameServerPolicy             *orderedmap.OrderedMap[string, any] `yaml:"nameserver-policy,omitempty" json:"nameserver-policy"`
	ProxyServerNameserver        []string                            `yaml:"proxy-server-nameserver,omitempty" json:"proxy-server-nameserver"`
	DirectNameServer             []string                            `yaml:"direct-nameserver,omitempty" json:"direct-nameserver"`
	DirectNameServerFollowPolicy bool                                `yaml:"direct-nameserver-follow-policy,omitempty" json:"direct-nameserver-follow-policy"`
}

type RawNTP struct {
	Enable        bool   `yaml:"enable,omitempty" json:"enable"`
	Server        string `yaml:"server,omitempty" json:"server"`
	Port          int    `yaml:"port,omitempty" json:"port"`
	Interval      int    `yaml:"interval,omitempty" json:"interval"`
	DialerProxy   string `yaml:"dialer-proxy,omitempty" json:"dialer-proxy"`
	WriteToSystem bool   `yaml:"write-to-system,omitempty" json:"write-to-system"`
}

type RawTun struct {
	Enable              bool        `yaml:"enable,omitempty" json:"enable"`
	Device              string      `yaml:"device,omitempty" json:"device"`
	Stack               CC.TUNStack `yaml:"stack,omitempty" json:"stack"`
	DNSHijack           []string    `yaml:"dns-hijack,omitempty" json:"dns-hijack"`
	AutoRoute           bool        `yaml:"auto-route,omitempty" json:"auto-route"`
	AutoDetectInterface bool        `yaml:"auto-detect-interface,omitempty"`

	MTU        uint32 `yaml:"mtu,omitempty" json:"mtu,omitempty"`
	GSO        bool   `yaml:"gso,omitempty" json:"gso,omitempty"`
	GSOMaxSize uint32 `yaml:"gso-max-size,omitempty" json:"gso-max-size,omitempty"`
	//Inet4Address           []netip.Prefix `yaml:"inet4-address,omitempty" json:"inet4-address,omitempty"`
	Inet6Address           []netip.Prefix `yaml:"inet6-address,omitempty" json:"inet6-address,omitempty"`
	IPRoute2TableIndex     int            `yaml:"iproute2-table-index,omitempty" json:"iproute2-table-index,omitempty"`
	IPRoute2RuleIndex      int            `yaml:"iproute2-rule-index,omitempty" json:"iproute2-rule-index,omitempty"`
	AutoRedirect           bool           `yaml:"auto-redirect,omitempty" json:"auto-redirect,omitempty"`
	AutoRedirectInputMark  uint32         `yaml:"auto-redirect-input-mark,omitempty" json:"auto-redirect-input-mark,omitempty"`
	AutoRedirectOutputMark uint32         `yaml:"auto-redirect-output-mark,omitempty" json:"auto-redirect-output-mark,omitempty"`
	StrictRoute            bool           `yaml:"strict-route,omitempty" json:"strict-route,omitempty"`
	RouteAddress           []netip.Prefix `yaml:"route-address,omitempty" json:"route-address,omitempty"`
	RouteAddressSet        []string       `yaml:"route-address-set,omitempty" json:"route-address-set,omitempty"`
	RouteExcludeAddress    []netip.Prefix `yaml:"route-exclude-address,omitempty" json:"route-exclude-address,omitempty"`
	RouteExcludeAddressSet []string       `yaml:"route-exclude-address-set,omitempty" json:"route-exclude-address-set,omitempty"`
	IncludeInterface       []string       `yaml:"include-interface,omitempty" json:"include-interface,omitempty"`
	ExcludeInterface       []string       `yaml:"exclude-interface,omitempty" json:"exclude-interface,omitempty"`
	IncludeUID             []uint32       `yaml:"include-uid,omitempty" json:"include-uid,omitempty"`
	IncludeUIDRange        []string       `yaml:"include-uid-range,omitempty" json:"include-uid-range,omitempty"`
	ExcludeUID             []uint32       `yaml:"exclude-uid,omitempty" json:"exclude-uid,omitempty"`
	ExcludeUIDRange        []string       `yaml:"exclude-uid-range,omitempty" json:"exclude-uid-range,omitempty"`
	ExcludeSrcPort         []uint16       `yaml:"exclude-src-port,omitempty" json:"exclude-src-port,omitempty"`
	ExcludeSrcPortRange    []string       `yaml:"exclude-src-port-range,omitempty" json:"exclude-src-port-range,omitempty"`
	ExcludeDstPort         []uint16       `yaml:"exclude-dst-port,omitempty" json:"exclude-dst-port,omitempty"`
	ExcludeDstPortRange    []string       `yaml:"exclude-dst-port-range,omitempty" json:"exclude-dst-port-range,omitempty"`
	IncludeAndroidUser     []int          `yaml:"include-android-user,omitempty" json:"include-android-user,omitempty"`
	IncludePackage         []string       `yaml:"include-package,omitempty" json:"include-package,omitempty"`
	ExcludePackage         []string       `yaml:"exclude-package,omitempty" json:"exclude-package,omitempty"`
	EndpointIndependentNat bool           `yaml:"endpoint-independent-nat,omitempty" json:"endpoint-independent-nat,omitempty"`
	UDPTimeout             int64          `yaml:"udp-timeout,omitempty" json:"udp-timeout,omitempty"`
	FileDescriptor         int            `yaml:"file-descriptor,omitempty" json:"file-descriptor"`

	Inet4RouteAddress        []netip.Prefix `yaml:"inet4-route-address,omitempty" json:"inet4-route-address,omitempty"`
	Inet6RouteAddress        []netip.Prefix `yaml:"inet6-route-address,omitempty" json:"inet6-route-address,omitempty"`
	Inet4RouteExcludeAddress []netip.Prefix `yaml:"inet4-route-exclude-address,omitempty" json:"inet4-route-exclude-address,omitempty"`
	Inet6RouteExcludeAddress []netip.Prefix `yaml:"inet6-route-exclude-address,omitempty" json:"inet6-route-exclude-address,omitempty"`
}

type RawTuicServer struct {
	Enable                bool              `yaml:"enable,omitempty" json:"enable"`
	Listen                string            `yaml:"listen,omitempty" json:"listen"`
	Token                 []string          `yaml:"token,omitempty" json:"token"`
	Users                 map[string]string `yaml:"users,omitempty" json:"users,omitempty"`
	Certificate           string            `yaml:"certificate,omitempty" json:"certificate"`
	PrivateKey            string            `yaml:"private-key,omitempty" json:"private-key"`
	CongestionController  string            `yaml:"congestion-controller,omitempty" json:"congestion-controller,omitempty"`
	MaxIdleTime           int               `yaml:"max-idle-time,omitempty" json:"max-idle-time,omitempty"`
	AuthenticationTimeout int               `yaml:"authentication-timeout,omitempty" json:"authentication-timeout,omitempty"`
	ALPN                  []string          `yaml:"alpn,omitempty" json:"alpn,omitempty"`
	MaxUdpRelayPacketSize int               `yaml:"max-udp-relay-packet-size,omitempty" json:"max-udp-relay-packet-size,omitempty"`
	CWND                  int               `yaml:"cwnd,omitempty" json:"cwnd,omitempty"`
}

type RawIPTables struct {
	Enable           bool     `yaml:"enable,omitempty" json:"enable"`
	InboundInterface string   `yaml:"inbound-interface,omitempty" json:"inbound-interface"`
	Bypass           []string `yaml:"bypass,omitempty" json:"bypass"`
	DnsRedirect      bool     `yaml:"dns-redirect,omitempty" json:"dns-redirect"`
}

type RawExperimental struct {
	Fingerprints     []string `yaml:"fingerprints,omitempty"`
	QUICGoDisableGSO bool     `yaml:"quic-go-disable-gso,omitempty"`
	QUICGoDisableECN bool     `yaml:"quic-go-disable-ecn,omitempty"`
	IP4PEnable       bool     `yaml:"dialer-ip4p-convert,omitempty"`
}

type RawProfile struct {
	StoreSelected bool `yaml:"store-selected,omitempty" json:"store-selected"`
	StoreFakeIP   bool `yaml:"store-fake-ip,omitempty" json:"store-fake-ip"`
}

type RawGeoXUrl struct {
	GeoIp   string `yaml:"geoip,omitempty" json:"geoip"`
	Mmdb    string `yaml:"mmdb,omitempty" json:"mmdb"`
	ASN     string `yaml:"asn,omitempty" json:"asn"`
	GeoSite string `yaml:"geosite,omitempty" json:"geosite"`
}

type RawSniffer struct {
	Enable          bool     `yaml:"enable,omitempty" json:"enable"`
	OverrideDest    bool     `yaml:"override-destination,omitempty" json:"override-destination"`
	Sniffing        []string `yaml:"sniffing,omitempty" json:"sniffing"`
	ForceDomain     []string `yaml:"force-domain,omitempty" json:"force-domain"`
	SkipSrcAddress  []string `yaml:"skip-src-address,omitempty" json:"skip-src-address"`
	SkipDstAddress  []string `yaml:"skip-dst-address,omitempty" json:"skip-dst-address"`
	SkipDomain      []string `yaml:"skip-domain,omitempty" json:"skip-domain"`
	Ports           []string `yaml:"port-whitelist,omitempty" json:"port-whitelist"`
	ForceDnsMapping bool     `yaml:"force-dns-mapping,omitempty" json:"force-dns-mapping"`
	ParsePureIp     bool     `yaml:"parse-pure-ip,omitempty" json:"parse-pure-ip"`

	Sniff map[string]C.RawSniffingConfig `yaml:"sniff,omitempty" json:"sniff"`
}

type RawTLS struct {
	Certificate     string   `yaml:"certificate,omitempty" json:"certificate"`
	PrivateKey      string   `yaml:"private-key,omitempty" json:"private-key"`
	EchKey          string   `yaml:"ech-key,omitempty" json:"ech-key"`
	CustomTrustCert []string `yaml:"custom-certifactes,omitempty" json:"custom-certifactes"`
}
