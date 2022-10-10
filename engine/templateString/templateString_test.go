package templateString

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestTemplateStringToStringWithOneKey(t *testing.T) {
	ts := TemplateString{Keys: []string{"id"}, Format: "%s"}
	c := make(map[string]string)
	c["id"] = "value"

	s, _ := ts.ToString(c)

	if s != c["id"] {
		t.Error("String should be equal to value, result:", s)
	}
}

func TestTemplateStringToStringWithTwoKeys(t *testing.T) {
	ts := TemplateString{Keys: []string{"id", "token"}, Format: "%s:%s"}
	c := make(map[string]string)
	c["id"] = "id-value"
	c["token"] = "token-value"

	s, _ := ts.ToString(c)

	if s != c["id"]+":"+c["token"] {
		t.Error("String should be equal to id-value:id-token, result:", s)
	}
}

func TestTemplateStringToStringNoMatchingKey(t *testing.T) {
	ts := TemplateString{Keys: []string{"token"}, Format: "%s"}
	c := make(map[string]string)
	c["id"] = "value"

	s, _ := ts.ToString(c)

	if s != "" {
		t.Error("String should be empty, result:", s)
	}
}

func TestTemplateStringUnmarshalYAML(t *testing.T) {
	getUnmarshall := func(txt string) *yaml.Node {
		nodeStruct := struct {
			node yaml.Node
		}{}
		tpl := fmt.Sprintf("node: %s", txt)
		_ = yaml.Unmarshal([]byte(tpl), &nodeStruct)
		return &nodeStruct.node
	}

	type fields struct {
		Format string
		Keys   []string
	}
	type args struct {
		node *yaml.Node
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "string without template",
			fields: fields{
				Format: "toto",
				Keys:   []string{},
			},
			args: args{
				node: getUnmarshall("toto"),
			},
			wantErr: false,
		},
		{
			name: "string with template",
			fields: fields{
				Format: "toto %s tata",
				Keys:   []string{"myVar"},
			},
			args: args{
				node: getUnmarshall("toto ${var.myVar} tata"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &TemplateString{
				Format: tt.fields.Format,
				Keys:   tt.fields.Keys,
			}
			if err := ts.UnmarshalYAML(tt.args.node); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
