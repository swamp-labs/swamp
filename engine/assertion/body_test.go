package assertion

import (
	"testing"
)

func TestValidateBodyTrue(t *testing.T) {
	var a BodyAssertion = &jsonPath{"$.id", "id"}
	raw := []byte(`{"id":"123"}`)

	boolean, err := a.validate(raw, map[string]string{})
	if !boolean {
		t.Error("validateBody was incorrect, should be true, result:", boolean)
	}
	if err != nil {
		t.Error("validateBody was incorrect", err)
	}
}

func TestValidateBodyFalse(t *testing.T) {
	var a BodyAssertion = &jsonPath{"$.id", "id"}
	raw := []byte(`{"code":"123"}`)

	_, err := a.validate(raw, map[string]string{})

	if err == nil {
		t.Error("validateBody was incorrect", err)
	}
}
