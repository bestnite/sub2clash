package test

import (
	"testing"

	"github.com/bestnite/sub2clash/model/proxy"
	"gopkg.in/yaml.v3"
)

type testStruct struct {
	A proxy.IntOrString `yaml:"a"`
}

func TestUnmarshal(t *testing.T) {
	yamlData1 := `a: 123`
	res := testStruct{}
	err := yaml.Unmarshal([]byte(yamlData1), &res)
	if err != nil {
		t.Errorf("failed to unmarshal yaml: %v", err)
	}
	if res.A != 123 {
		t.Errorf("expected 123, but got %v", res.A)
	}

	yamlData2 := `a: "123"`
	err = yaml.Unmarshal([]byte(yamlData2), &res)
	if err != nil {
		t.Errorf("failed to unmarshal yaml: %v", err)
	}
	if res.A != 123 {
		t.Errorf("expected 123, but got %v", res.A)
	}
}
