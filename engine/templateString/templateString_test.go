package templateString

import (
	"testing"
)

func TestTemplateString_ToStringWithOneKey(t *testing.T) {
	ts := TemplateString{Keys: []string{"id"}, Format: "%s"}
	c := make(map[string]string)
	c["id"] = "id-value"

	s, err := ts.ToString(c)

	if err != nil {
		t.Error("ToString was incorrect", err)
	}
	if s != c["id"] {
		t.Error("String should be equal to id-value, result:", s)
	}
}

func TestTemplateString_ToStringWithTwoKeys(t *testing.T) {
	ts := TemplateString{Keys: []string{"id", "token"}, Format: "%s:%s"}
	c := make(map[string]string)
	c["id"] = "id-value"
	c["token"] = "token-value"

	s, err := ts.ToString(c)

	if err != nil {
		t.Error("ToString was incorrect", err)
	}
	if s != c["id"]+":"+c["token"] {
		t.Error("String should be equal to id-value:id-token, result:", s)
	}
}

func TestTemplateString_ToStringNoMatchingKey(t *testing.T) {
	ts := TemplateString{Keys: []string{"token"}, Format: "%s"}
	c := make(map[string]string)
	c["id"] = "id-value"

	s, err := ts.ToString(c)

	if err != nil {
		t.Error("ToString was incorrect", err)
	}
	if s != "" {
		t.Error("String should be empty, result:", s)
	}
}
