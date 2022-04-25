package assertion

import (
	"testing"
)

func TestValidateWithRegexTrue(t *testing.T) {
	raw := []byte(`{"id":"123"}`)

	boolean, err := validateWithRegex(raw, "(\\d{1,3})")

	if !boolean {
		t.Error("validateWithRegex was incorrect, should be true, result:", boolean)
	}
	if err != nil {
		t.Error("validateWithRegex was incorrect", err)
	}
}

func TestValidateWithRegexFalse(t *testing.T) {
	raw := []byte(`{"message":"Hello World"}`)

	boolean, err := validateWithRegex(raw, "(\\d{1,3})")

	if boolean {
		t.Error("validateWithRegex was incorrect, should be false, result:", boolean)
	}
	if err != nil {
		t.Error("validateWithRegex was incorrect", err)
	}
}

func TestGetFromRegexEmptyExp(t *testing.T) {
	raw := []byte(`{"id":"123"}`)

	result, _ := getFromRegex(raw, "")

	if result != "" {
		t.Error("result should be empty with no regex BodyAssertion, result:", result)
	}
}

func TestGetFromRegexNoSubMatch(t *testing.T) {
	raw := []byte(`{"id":"123"}`)

	result, _ := getFromRegex(raw, "123")

	if result != "" {
		t.Error("result should be empty with no submatch BodyAssertion, result:", result)
	}
}

func TestGetFromJsonPath(t *testing.T) {
	raw := []byte(`{"message":"Hello World"}`)

	boolean, returnedString, err := getFromJsonPath(raw, "$.message")
	if err != nil {
		t.Error("getFromJsonPath was incorrect", err)
	}

	if !boolean || returnedString != "Hello World" {
		t.Errorf("getFromJsonPath was incorrect, got: %s, want: %s.", returnedString, "Hello World")
	}

}

func TestGetFromJsonPathEmptyExp(t *testing.T) {
	raw := []byte(`{"message":"Hello World"}`)

	_, _, err := getFromJsonPath(raw, "")
	if err == nil {
		t.Error("error should not be nil with no jsonpath BodyAssertion", err)
	}

}

func TestGetFromJsonPathEmptyRes(t *testing.T) {
	raw := []byte(`{"message":"Hello World"}`)

	_, returnedString, err := getFromJsonPath(raw, "$.id")
	if err == nil {
		t.Error("getFromJsonPath was incorrect", err)
	}
	if returnedString != "" {
		t.Error("getFromJsonPath was incorrect, returnedString should be empty", returnedString)
	}
}
