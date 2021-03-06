package assertion

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestAssertResponseTrue(t *testing.T) {
	var a assertion
	a.Code = []int{201}
	a.Body = []BodyAssertion{&jsonPath{"$.id", "id"}, &regex{"(\\d{1,3})", "id"}}
	var r http.Response
	r.Body = io.NopCloser(strings.NewReader(`{"id":"123"}`))
	r.StatusCode = 201
	boolean, err := a.AssertResponse(&r, map[string]string{})

	if !boolean {
		t.Error("AssertResponse was incorrect, should be true, result:", boolean)
	}
	if err != nil {
		t.Error("AssertResponseTrue was incorrect", err)
	}
}

func TestAssertResponseCodeFalse(t *testing.T) {
	var a assertion
	a.Code = []int{500}
	a.Body = []BodyAssertion{&jsonPath{"$.id", "id"}}
	var r http.Response
	r.Body = io.NopCloser(strings.NewReader(`{"id":"123"}`))
	r.StatusCode = 201
	boolean, err := a.AssertResponse(&r, map[string]string{})

	if boolean {
		t.Error("AssertResponse was incorrect, should be false, result:", boolean)
	}
	if err != nil {
		t.Error("AssertResponseCodeFalse was incorrect", err)
	}
}

func TestAssertResponseBodyFalse(t *testing.T) {
	var a assertion
	a.Code = []int{500}
	a.Body = []BodyAssertion{&jsonPath{"$.id", "id"}, &regex{"(\\d{1,3})", "id"}}
	var r http.Response
	r.Body = io.NopCloser(strings.NewReader(`{"code":"123"}`))
	r.StatusCode = 201
	boolean, err := a.AssertResponse(&r, map[string]string{})

	if boolean {
		t.Error("AssertResponse was incorrect, should be true, result:", boolean)
	}
	if err != nil {
		t.Error("AssertResponseBodyFalse was incorrect", err)
	}
}
