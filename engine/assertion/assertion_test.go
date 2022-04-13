package assertion

import (
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestAssertResponse(t *testing.T) {
	var a Assertion
	a.Code = []int{201}
	a.Body = []expression{&JsonPath{"$.id", "id"}, &Regex{"(\\d{1,3})", "id"}}
	var r http.Response
	r.Body = io.NopCloser(strings.NewReader(`{"id":"123"}`))
	r.StatusCode = 201
	boolean, err := AssertResponse(a, &r, map[string]string{})

	if !boolean {
		t.Error("validate was incorrect, should be true, result:", boolean)
	}
	if err != nil {
		t.Error("validate was incorrect", err)
	}
}

func TestValidateBody(t *testing.T) {
	var a expression = &JsonPath{"$.id", "id"}
	raw := []byte(`{"id":"123"}`)

	boolean, err := a.validate(raw, map[string]string{})
	log.Println(boolean)
	if !boolean {
		t.Error("validate was incorrect, should be true, result:", boolean)
	}
	if err != nil {
		t.Error("validate was incorrect", err)
	}
}
