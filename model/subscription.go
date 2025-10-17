package model

import (
	"net/netip"

	"github.com/bestnite/sub2clash/model/proxy"
	C "github.com/metacubex/mihomo/config"
	LC "github.com/metacubex/mihomo/listener/config"
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
	DNS           C.RawDNS                  `yaml:"dns,omitempty" json:"dns"`
	NTP           C.RawNTP                  `yaml:"ntp,omitempty" json:"ntp"`
	Tun           C.RawTun                  `yaml:"tun,omitempty" json:"tun"`
	TuicServer    C.RawTuicServer           `yaml:"tuic-server,omitempty" json:"tuic-server"`
	IPTables      C.RawIPTables             `yaml:"iptables,omitempty" json:"iptables"`
	Experimental  C.RawExperimental         `yaml:"experimental,omitempty" json:"experimental"`
	Profile       C.RawProfile              `yaml:"profile,omitempty" json:"profile"`
	GeoXUrl       C.RawGeoXUrl              `yaml:"geox-url,omitempty" json:"geox-url"`
	Sniffer       C.RawSniffer              `yaml:"sniffer,omitempty" json:"sniffer"`
	TLS           C.RawTLS                  `yaml:"tls,omitempty" json:"tls"`

	ClashForAndroid C.RawClashForAndroid `yaml:"clash-for-android,omitempty" json:"clash-for-android"`
}
