package model

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/bestnite/sub2clash/utils"
	"github.com/gin-gonic/gin"
)

type ConvertConfig struct {
	ClashType           ClashType            `json:"clashType" binding:"required"`
	Subs                []string             `json:"subscriptions" binding:""`
	Proxies             []string             `json:"proxies" binding:""`
	Refresh             bool                 `json:"refresh" binding:""`
	Template            string               `json:"template" binding:""`
	RuleProviders       []RuleProviderStruct `json:"ruleProviders" binding:""`
	Rules               []RuleStruct         `json:"rules" binding:""`
	AutoTest            bool                 `json:"autoTest" binding:""`
	Lazy                bool                 `json:"lazy" binding:""`
	Sort                string               `json:"sort" binding:""`
	Remove              string               `json:"remove" binding:""`
	Replace             map[string]string    `json:"replace" binding:""`
	NodeListMode        bool                 `json:"nodeList" binding:""`
	IgnoreCountryGrooup bool                 `json:"ignoreCountryGroup" binding:""`
	UserAgent           string               `json:"userAgent" binding:""`
	UseUDP              bool                 `json:"useUDP" binding:""`
}

type RuleProviderStruct struct {
	Behavior string `json:"behavior" binding:""`
	Url      string `json:"url" binding:""`
	Group    string `json:"group" binding:""`
	Prepend  bool   `json:"prepend" binding:""`
	Name     string `json:"name" binding:""`
}

type RuleStruct struct {
	Rule    string `json:"rule" binding:""`
	Prepend bool   `json:"prepend" binding:""`
}

func ParseConvertQuery(c *gin.Context) (ConvertConfig, error) {
	config := c.Param("config")
	queryBytes, err := utils.DecodeBase64(config, true)
	if err != nil {
		return ConvertConfig{}, errors.New("参数错误: " + err.Error())
	}
	var query ConvertConfig
	err = json.Unmarshal([]byte(queryBytes), &query)
	if err != nil {
		return ConvertConfig{}, errors.New("参数错误: " + err.Error())
	}
	if len(query.Subs) == 0 && len(query.Proxies) == 0 {
		return ConvertConfig{}, errors.New("参数错误: sub 和 proxy 不能同时为空")
	}
	if len(query.Subs) > 0 {
		for i := range query.Subs {
			if !strings.HasPrefix(query.Subs[i], "http") {
				return ConvertConfig{}, errors.New("参数错误: sub 格式错误")
			}
			if _, err := url.ParseRequestURI(query.Subs[i]); err != nil {
				return ConvertConfig{}, errors.New("参数错误: " + err.Error())
			}
		}
	} else {
		query.Subs = nil
	}
	if query.Template != "" {
		if strings.HasPrefix(query.Template, "http") {
			uri, err := url.ParseRequestURI(query.Template)
			if err != nil {
				return ConvertConfig{}, err
			}
			query.Template = uri.String()
		}
	}
	if len(query.RuleProviders) > 0 {
		names := make(map[string]bool)
		for _, ruleProvider := range query.RuleProviders {
			if _, ok := names[ruleProvider.Name]; ok {
				return ConvertConfig{}, errors.New("参数错误: Rule-Provider 名称重复")
			}
			names[ruleProvider.Name] = true
		}
	} else {
		query.RuleProviders = nil
	}
	return query, nil
}
