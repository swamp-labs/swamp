package assertion

import (
	"testing"
)

func TestValidateWithRegex_True(t *testing.T) {
	raw := []byte(`{"id":"123"}`)

	boolean, err := validateWithRegex(raw, "(\\d{1,3})")

	if !boolean {
		t.Error("validateWithRegex was incorrect, should be true, result:", boolean)
	}
	if err != nil {
		t.Error("validateWithRegex was incorrect", err)
	}
}

func TestValidateWithRegex_False(t *testing.T) {
	raw := []byte(`{"message":"Hello World"}`)

	boolean, err := validateWithRegex(raw, "(\\d{1,3})")

	if boolean {
		t.Error("validateWithRegex was incorrect, should be false, result:", boolean)
	}
	if err != nil {
		t.Error("validateWithRegex was incorrect", err)
	}
}

func TestGetFromRegex_EmptyExp(t *testing.T) {
	raw := []byte(`{"id":"123"}`)

	result, _ := getFromRegex(raw, "")

	if result != "" {
		t.Error("result should be empty with no regex expression, result:", result)
	}
}

func TestGetFromRegex_NoSubMatch(t *testing.T) {
	raw := []byte(`{"id":"123"}`)

	result, _ := getFromRegex(raw, "123")

	if result != "" {
		t.Error("result should be empty with no submatch expression, result:", result)
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

func TestGetFromJsonPath_EmptyExp(t *testing.T) {
	raw := []byte(`{"message":"Hello World"}`)

	_, _, err := getFromJsonPath(raw, "")
	if err == nil {
		t.Error("error should not be nil with no jsonpath expression", err)
	}

}

func TestGetFromJsonPath_EmptyRes(t *testing.T) {
	raw := []byte(`{"message":"Hello World"}`)

	_, returnedString, err := getFromJsonPath(raw, "$.id")
	if err == nil {
		t.Error("getFromJsonPath was incorrect", err)
	}
	if returnedString != "" {
		t.Error("getFromJsonPath was incorrect, returnedString should be empty", returnedString)
	}
}
