package assertion

import (
	"log"
	"testing"
)

func TestValidateHeadersTrue(t *testing.T) {
	var a Assertion
	a.Code = []int{200}

	boolean := a.validateCodeStatus(200)
	if !boolean {
		t.Error("validateCodeStatus was incorrect, should be true, result:", boolean)
	}
}

func TestValidateHeadersFalse(t *testing.T) {
	var a Assertion
	a.Code = []int{200}

	boolean := a.validateCodeStatus(500)
	if boolean {
		t.Error("validateCodeStatus was incorrect, should be false, result:", boolean)
	}
}

func TestValidateHeadersEmpty(t *testing.T) {
	var a Assertion
	a.Code = nil

	boolean := a.validateCodeStatus(200)
	log.Println(boolean)
	if !boolean {
		t.Error("validateCodeStatus was incorrect, should be true, result:", boolean)
	}
}
