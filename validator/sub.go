package validator

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type SubValidator struct {
	Sub                 string               `form:"sub" binding:""`
	Subs                []string             `form:"-" binding:""`
	Proxy               string               `form:"proxy" binding:""`
	Proxies             []string             `form:"-" binding:""`
	Refresh             bool                 `form:"refresh,default=false" binding:""`
	SkipErrors 			bool 				 `form:"skipErrors,default=false" binding:"""`
	Template            string               `form:"template" binding:""`
	RuleProvider        string               `form:"ruleProvider" binding:""`
	RuleProviders       []RuleProviderStruct `form:"-" binding:""`
	Rule                string               `form:"rule" binding:""`
	Rules               []RuleStruct         `form:"-" binding:""`
	AutoTest            bool                 `form:"autoTest,default=false" binding:""`
	Lazy                bool                 `form:"lazy,default=false" binding:""`
	Sort                string               `form:"sort" binding:""`
	Remove              string               `form:"remove" binding:""`
	Replace             string               `form:"replace" binding:""`
	ReplaceKeys         []string             `form:"-" binding:""`
	ReplaceTo           []string             `form:"-" binding:""`
	NodeListMode        bool                 `form:"nodeList,default=false" binding:""`
	IgnoreCountryGrooup bool                 `form:"ignoreCountryGroup,default=false" binding:""`
	UserAgent           string               `form:"userAgent" binding:""`
}

type RuleProviderStruct struct {
	Behavior string
	Url      string
	Group    string
	Prepend  bool
	Name     string
}

type RuleStruct struct {
	Rule    string
	Prepend bool
}

func ParseQuery(c *gin.Context) (SubValidator, error) {
	var query SubValidator
	if err := c.ShouldBind(&query); err != nil {
		return SubValidator{}, errors.New("参数错误: " + err.Error())
	}
	if query.Sub == "" && query.Proxy == "" {
		return SubValidator{}, errors.New("参数错误: sub 和 proxy 不能同时为空")
	}
	if query.Sub != "" {
		query.Subs = strings.Split(query.Sub, ",")
		for i := range query.Subs {
			if !strings.HasPrefix(query.Subs[i], "http") {
				return SubValidator{}, errors.New("参数错误: sub 格式错误")
			}
			if _, err := url.ParseRequestURI(query.Subs[i]); err != nil {
				return SubValidator{}, errors.New("参数错误: " + err.Error())
			}
		}
	} else {
		query.Subs = nil
	}
	if query.Proxy != "" {
		query.Proxies = strings.Split(query.Proxy, ",")
	} else {
		query.Proxies = nil
	}
	if query.Template != "" {
		if strings.HasPrefix(query.Template, "http") {
			uri, err := url.ParseRequestURI(query.Template)
			if err != nil {
				return SubValidator{}, err
			}
			query.Template = uri.String()
		}
	}
	if query.RuleProvider != "" {
		reg := regexp.MustCompile(`\[(.*?)\]`)
		ruleProviders := reg.FindAllStringSubmatch(query.RuleProvider, -1)
		for i := range ruleProviders {
			length := len(ruleProviders)
			parts := strings.Split(ruleProviders[length-i-1][1], ",")
			if len(parts) < 4 {
				return SubValidator{}, errors.New("参数错误: ruleProvider 格式错误")
			}
			u := parts[1]
			uri, err := url.ParseRequestURI(u)
			if err != nil {
				return SubValidator{}, errors.New("参数错误: " + err.Error())
			}
			u = uri.String()
			if len(parts) == 4 {
				hash := sha256.Sum224([]byte(u))
				parts = append(parts, hex.EncodeToString(hash[:]))
			}
			query.RuleProviders = append(
				query.RuleProviders, RuleProviderStruct{
					Behavior: parts[0],
					Url:      u,
					Group:    parts[2],
					Prepend:  parts[3] == "true",
					Name:     parts[4],
				},
			)
		}

		names := make(map[string]bool)
		for _, ruleProvider := range query.RuleProviders {
			if _, ok := names[ruleProvider.Name]; ok {
				return SubValidator{}, errors.New("参数错误: Rule-Provider 名称重复")
			}
			names[ruleProvider.Name] = true
		}
	} else {
		query.RuleProviders = nil
	}
	if query.Rule != "" {
		reg := regexp.MustCompile(`\[(.*?)\]`)
		rules := reg.FindAllStringSubmatch(query.Rule, -1)
		for i := range rules {
			length := len(rules)
			r := rules[length-1-i][1]
			strings.LastIndex(r, ",")
			parts := [2]string{}
			parts[0] = r[:strings.LastIndex(r, ",")]
			parts[1] = r[strings.LastIndex(r, ",")+1:]
			query.Rules = append(
				query.Rules, RuleStruct{
					Rule:    parts[0],
					Prepend: parts[1] == "true",
				},
			)
		}
	} else {
		query.Rules = nil
	}
	if strings.TrimSpace(query.Replace) != "" {
		reg := regexp.MustCompile(`\[<(.*?)>,<(.*?)>\]`)
		replaces := reg.FindAllStringSubmatch(query.Replace, -1)
		for i := range replaces {
			length := len(replaces[i])
			if length != 3 {
				return SubValidator{}, errors.New("参数错误: replace 格式错误")
			}
			query.ReplaceKeys = append(query.ReplaceKeys, replaces[i][1])
			query.ReplaceTo = append(query.ReplaceTo, replaces[i][2])
		}
	}
	return query, nil
}
