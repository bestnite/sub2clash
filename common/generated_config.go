package common

import (
	P "github.com/bestnite/sub2clash/model/proxy"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// proxyListDoc 只用于解析 YAML 订阅中的 proxies 字段。
// 方案 A/B 下我们不再关心订阅 YAML 里的其他 mihomo 配置项。
type proxyListDoc struct {
	Proxy []P.Proxy `yaml:"proxies,omitempty"`
}

// generatedConfig 是运行期的最小叠加模型：
// 只保留本项目真正会读取、生成或修改的字段。
//
// 这里承载的是“本项目的业务叠加层”，而不是 mihomo 的完整配置模型：
// - Proxy: 解析出的节点，用于过滤、去重、分组等中间处理
// - ProxyGroup: 模板中需要参与占位符展开的组，以及本项目生成的国家组
// - Rule: 模板规则 + 用户追加规则，用于保持 MATCH 规则前插入的语义
type generatedConfig struct {
	Proxy      []P.Proxy        `yaml:"proxies,omitempty"`
	ProxyGroup []generatedGroup `yaml:"proxy-groups,omitempty"`
	Rule       []string         `yaml:"rules,omitempty"`
}

// generatedGroup 表示本项目生成出来的代理组最小模型，
// 它不再镜像 mihomo 的完整 proxy-group 配置结构。
//
// 这里只保留“当前逻辑真正需要读写的字段”：
// - Name / Proxies：用于模板占位符展开与 patch
// - Type / Url / Interval / Tolerance / Lazy：用于输出自动测速国家组
// - Size / IsCountry：仅作为运行期辅助信息，不参与 YAML 输出
type generatedGroup struct {
	Type      string   `yaml:"type,omitempty"`
	Name      string   `yaml:"name,omitempty"`
	Proxies   []string `yaml:"proxies,omitempty"`
	Url       string   `yaml:"url,omitempty"`
	Interval  int      `yaml:"interval,omitempty"`
	Tolerance int      `yaml:"tolerance,omitempty"`
	Lazy      bool     `yaml:"lazy"`
	Size      int      `yaml:"-"`
	IsCountry bool     `yaml:"-"`
}

// generatedRulePatch 表示本项目追加/覆盖的 rule-provider 最小模型。
// 它仅用于把用户请求转换成对 templateDoc 的字段级 patch。
type generatedRulePatch struct {
	Type     string `yaml:"type,omitempty"`
	Behavior string `yaml:"behavior,omitempty"`
	Url      string `yaml:"url,omitempty"`
	Path     string `yaml:"path,omitempty"`
	Interval int    `yaml:"interval,omitempty"`
	Format   string `yaml:"format,omitempty"`
}

type generatedGroupsSortByName []generatedGroup
type generatedGroupsSortBySize []generatedGroup

func (p generatedGroupsSortByName) Len() int {
	return len(p)
}

func (p generatedGroupsSortBySize) Len() int {
	return len(p)
}

func (p generatedGroupsSortByName) Less(i, j int) bool {
	tags := []language.Tag{
		language.English,
		language.Chinese,
	}
	matcher := language.NewMatcher(tags)
	bestMatch, _, _ := matcher.Match(language.Make("zh"))
	c := collate.New(bestMatch)
	return c.CompareString(p[i].Name, p[j].Name) < 0
}

func (p generatedGroupsSortBySize) Less(i, j int) bool {
	if p[i].Size == p[j].Size {
		return p[i].Name < p[j].Name
	}
	return p[i].Size < p[j].Size
}

func (p generatedGroupsSortByName) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p generatedGroupsSortBySize) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
