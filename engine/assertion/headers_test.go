package assertion

import (
	"testing"
)

func TestValidateHeadersTrue(t *testing.T) {
	var a assertion
	a.Code = []int{200}

	boolean := a.validateCodeStatus(200)
	if !boolean {
		t.Error("validateCodeStatus was incorrect, should be true, result:", boolean)
	}
}

func TestValidateHeadersFalse(t *testing.T) {
	var a assertion
	a.Code = []int{200}

	boolean := a.validateCodeStatus(500)
	if boolean {
		t.Error("validateCodeStatus was incorrect, should be false, result:", boolean)
	}
}

func TestValidateHeadersEmpty(t *testing.T) {
	var a assertion
	a.Code = nil

	boolean := a.validateCodeStatus(200)
	if !boolean {
		t.Error("validateCodeStatus was incorrect, should be true, result:", boolean)
	}
}
