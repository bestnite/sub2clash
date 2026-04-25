package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bestnite/sub2clash/model"
	"gopkg.in/yaml.v3"
)

func withRepoRoot(t *testing.T) {
	t.Helper()
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}

	repoRoot := filepath.Dir(originalWD)
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("change working directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})
}

func TestBuildSubPreservesUnmodeledTemplateSections(t *testing.T) {
	withRepoRoot(t)

	templateName := "test_scheme_a_template.yaml"
	templatePath := filepath.Join(templatesDir, templateName)
	templateContent := `mixed-port: 7890
dns:
  enable: true
  future-field: true
new-section:
  enabled: true
proxies:
proxy-groups:
  - name: 节点选择
    type: select
    proxies:
      - <countries>
      - DIRECT
rules:
  - MATCH,节点选择
`

	if err := os.WriteFile(templatePath, []byte(templateContent), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(templatePath)
	})

	result, err := BuildSub(model.Clash, model.ConvertConfig{
		ClashType: model.Clash,
		Proxies: []string{
			"ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@127.0.0.1:8080#Test Node",
		},
	}, templateName, 0, 0)
	if err != nil {
		t.Fatalf("build subscription: %v", err)
	}

	output, err := yaml.Marshal(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}

	var doc map[string]any
	if err := yaml.Unmarshal(output, &doc); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	dns, ok := doc["dns"].(map[string]any)
	if !ok {
		t.Fatalf("dns section missing: %s", output)
	}
	if dns["future-field"] != true {
		t.Fatalf("dns future-field not preserved: %#v", dns)
	}

	newSection, ok := doc["new-section"].(map[string]any)
	if !ok {
		t.Fatalf("new-section missing: %s", output)
	}
	if newSection["enabled"] != true {
		t.Fatalf("new-section not preserved: %#v", newSection)
	}

	proxies, ok := doc["proxies"].([]any)
	if !ok || len(proxies) != 1 {
		t.Fatalf("expected generated proxies in output: %#v", doc["proxies"])
	}

	rules, ok := doc["rules"].([]any)
	if !ok || len(rules) != 1 || rules[0] != "MATCH,节点选择" {
		t.Fatalf("rules should stay untouched without rule patches: %#v", doc["rules"])
	}
}

func TestBuildSubPreservesTemplateProxyAndGroupFields(t *testing.T) {
	withRepoRoot(t)

	templateName := "test_scheme_a_group_template.yaml"
	templatePath := filepath.Join(templatesDir, templateName)
	templateContent := `proxies:
  - name: Template Proxy
    type: ss
    server: 1.1.1.1
    port: 443
    cipher: aes-256-gcm
    password: password
    future-proxy-field: keep
proxy-groups:
  - name: 节点选择
    type: select
    future-group-field: keep
    proxies:
      - <countries>
      - DIRECT
rules:
  - MATCH,节点选择
`

	if err := os.WriteFile(templatePath, []byte(templateContent), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(templatePath)
	})

	result, err := BuildSub(model.Clash, model.ConvertConfig{
		ClashType: model.Clash,
		Proxies: []string{
			"ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@127.0.0.1:8080#Test Node",
		},
	}, templateName, 0, 0)
	if err != nil {
		t.Fatalf("build subscription: %v", err)
	}

	output, err := yaml.Marshal(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}

	var doc map[string]any
	if err := yaml.Unmarshal(output, &doc); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	proxies, ok := doc["proxies"].([]any)
	if !ok || len(proxies) != 2 {
		t.Fatalf("expected two proxies in output: %#v", doc["proxies"])
	}
	firstProxy, ok := proxies[0].(map[string]any)
	if !ok {
		t.Fatalf("template proxy should remain a mapping: %#v", proxies[0])
	}
	if firstProxy["future-proxy-field"] != "keep" {
		t.Fatalf("template proxy field not preserved: %#v", firstProxy)
	}

	groups, ok := doc["proxy-groups"].([]any)
	if !ok || len(groups) == 0 {
		t.Fatalf("expected proxy groups in output: %#v", doc["proxy-groups"])
	}
	firstGroup, ok := groups[0].(map[string]any)
	if !ok {
		t.Fatalf("template group should remain a mapping: %#v", groups[0])
	}
	if firstGroup["future-group-field"] != "keep" {
		t.Fatalf("template proxy-group field not preserved: %#v", firstGroup)
	}

	groupProxies, ok := firstGroup["proxies"].([]any)
	if !ok || len(groupProxies) == 0 {
		t.Fatalf("template proxy-group proxies missing: %#v", firstGroup["proxies"])
	}
	for _, value := range groupProxies {
		if value == "<countries>" {
			t.Fatalf("placeholder should be resolved in template proxy-group: %#v", groupProxies)
		}
	}
}

func TestBuildSubAddsRulesForRuleProviderWhenTemplateHasNoRules(t *testing.T) {
	withRepoRoot(t)

	templateName := "test_scheme_a_rule_provider_template.yaml"
	templatePath := filepath.Join(templatesDir, templateName)
	templateContent := `proxy-groups:
  - name: 节点选择
    type: select
    proxies:
      - DIRECT
`

	if err := os.WriteFile(templatePath, []byte(templateContent), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(templatePath)
	})

	result, err := BuildSub(model.Clash, model.ConvertConfig{
		ClashType: model.Clash,
		Proxies: []string{
			"ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@127.0.0.1:8080#Test Node",
		},
		RuleProviders: []model.RuleProviderStruct{{
			Name:     "test-provider",
			Group:    "节点选择",
			Behavior: "domain",
			Url:      "https://example.com/rules.yaml",
		}},
	}, templateName, 0, 0)
	if err != nil {
		t.Fatalf("build subscription: %v", err)
	}

	output, err := yaml.Marshal(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}

	var doc map[string]any
	if err := yaml.Unmarshal(output, &doc); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	ruleProviders, ok := doc["rule-providers"].(map[string]any)
	if !ok {
		t.Fatalf("rule-providers missing: %#v", doc["rule-providers"])
	}
	if _, ok := ruleProviders["test-provider"]; !ok {
		t.Fatalf("test-provider missing: %#v", ruleProviders)
	}

	rules, ok := doc["rules"].([]any)
	if !ok || len(rules) != 1 || rules[0] != "RULE-SET,test-provider,节点选择" {
		t.Fatalf("expected generated rule for provider: %#v", doc["rules"])
	}
}

func TestBuildSubDoesNotInjectProxiesFieldIntoUseBasedGroup(t *testing.T) {
	withRepoRoot(t)

	templateName := "test_scheme_a_use_group_template.yaml"
	templatePath := filepath.Join(templatesDir, templateName)
	templateContent := `proxy-groups:
  - name: 节点选择
    type: select
    use:
      - provider-a
rules:
  - MATCH,节点选择
`

	if err := os.WriteFile(templatePath, []byte(templateContent), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(templatePath)
	})

	result, err := BuildSub(model.Clash, model.ConvertConfig{
		ClashType: model.Clash,
		Proxies: []string{
			"ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@127.0.0.1:8080#Test Node",
		},
	}, templateName, 0, 0)
	if err != nil {
		t.Fatalf("build subscription: %v", err)
	}

	output, err := yaml.Marshal(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}

	var doc map[string]any
	if err := yaml.Unmarshal(output, &doc); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	groups := doc["proxy-groups"].([]any)
	firstGroup := groups[0].(map[string]any)
	if _, exists := firstGroup["proxies"]; exists {
		t.Fatalf("use-based group should not gain proxies field: %#v", firstGroup)
	}
	if _, exists := firstGroup["use"]; !exists {
		t.Fatalf("use-based group should preserve use field: %#v", firstGroup)
	}
}

func TestBuildSubPreservesUnknownFieldsOnExistingRuleProvider(t *testing.T) {
	withRepoRoot(t)

	templateName := "test_scheme_a_existing_provider_template.yaml"
	templatePath := filepath.Join(templatesDir, templateName)
	templateContent := `proxy-groups:
  - name: 节点选择
    type: select
    proxies:
      - DIRECT
rule-providers:
  test-provider:
    type: http
    behavior: classical
    url: https://old.example.com/rules.yaml
    path: ./old.yaml
    interval: 10
    future-provider-field: keep
`

	if err := os.WriteFile(templatePath, []byte(templateContent), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(templatePath)
	})

	result, err := BuildSub(model.Clash, model.ConvertConfig{
		ClashType: model.Clash,
		Proxies: []string{
			"ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@127.0.0.1:8080#Test Node",
		},
		RuleProviders: []model.RuleProviderStruct{{
			Name:     "test-provider",
			Group:    "节点选择",
			Behavior: "domain",
			Url:      "https://example.com/rules.yaml",
		}},
	}, templateName, 0, 0)
	if err != nil {
		t.Fatalf("build subscription: %v", err)
	}

	output, err := yaml.Marshal(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}

	var doc map[string]any
	if err := yaml.Unmarshal(output, &doc); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	ruleProviders := doc["rule-providers"].(map[string]any)
	provider := ruleProviders["test-provider"].(map[string]any)
	if provider["future-provider-field"] != "keep" {
		t.Fatalf("existing provider field not preserved: %#v", provider)
	}
	if provider["behavior"] != "domain" {
		t.Fatalf("provider behavior not updated: %#v", provider)
	}
	if provider["url"] != "https://example.com/rules.yaml" {
		t.Fatalf("provider url not updated: %#v", provider)
	}
}

func TestBuildSubSkipsDuplicateCountryGroupNames(t *testing.T) {
	withRepoRoot(t)

	templateName := "test_scheme_a_country_group_template.yaml"
	templatePath := filepath.Join(templatesDir, templateName)
	templateContent := `proxy-groups:
  - name: 其他地区
    type: select
    proxies:
      - DIRECT
rules:
  - MATCH,其他地区
`

	if err := os.WriteFile(templatePath, []byte(templateContent), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(templatePath)
	})

	result, err := BuildSub(model.Clash, model.ConvertConfig{
		ClashType: model.Clash,
		Proxies: []string{
			"ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@127.0.0.1:8080#UnknownCountryNode",
		},
	}, templateName, 0, 0)
	if err != nil {
		t.Fatalf("build subscription: %v", err)
	}

	output, err := yaml.Marshal(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}

	var doc map[string]any
	if err := yaml.Unmarshal(output, &doc); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	groups := doc["proxy-groups"].([]any)
	count := 0
	for _, item := range groups {
		group := item.(map[string]any)
		if group["name"] == "其他地区" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("expected duplicate country group names to be skipped, got %d entries: %#v", count, groups)
	}
}

func TestBuiltSubMarshalNodeListYAMLUsesFinalYAMLTree(t *testing.T) {
	withRepoRoot(t)

	templateName := "test_scheme_a_nodelist_template.yaml"
	templatePath := filepath.Join(templatesDir, templateName)
	templateContent := `proxies:
  - name: Template Proxy
    type: ss
    server: 1.1.1.1
    port: 443
    cipher: aes-256-gcm
    password: password
    future-proxy-field: keep
proxy-groups:
  - name: 节点选择
    type: select
    proxies:
      - DIRECT
`

	if err := os.WriteFile(templatePath, []byte(templateContent), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(templatePath)
	})

	result, err := BuildSub(model.Clash, model.ConvertConfig{
		ClashType: model.Clash,
		Proxies: []string{
			"ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@127.0.0.1:8080#Generated Node",
		},
	}, templateName, 0, 0)
	if err != nil {
		t.Fatalf("build subscription: %v", err)
	}

	output, err := result.MarshalNodeListYAML()
	if err != nil {
		t.Fatalf("marshal node list: %v", err)
	}

	var doc map[string]any
	if err := yaml.Unmarshal(output, &doc); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	proxies, ok := doc["proxies"].([]any)
	if !ok || len(proxies) != 2 {
		t.Fatalf("expected node list to include template and generated proxies: %#v", doc["proxies"])
	}
	firstProxy, ok := proxies[0].(map[string]any)
	if !ok {
		t.Fatalf("template proxy should remain a mapping: %#v", proxies[0])
	}
	if firstProxy["future-proxy-field"] != "keep" {
		t.Fatalf("node list should be built from final yaml tree: %#v", firstProxy)
	}
}
