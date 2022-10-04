package templateString

import (
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
