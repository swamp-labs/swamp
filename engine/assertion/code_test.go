package assertion

import (
	"testing"
)

func TestValidateCodeStatusTrue(t *testing.T) {
	var a assertion
	a.Code = []int{200}

	boolean := a.validateCodeStatus(200)
	if !boolean {
		t.Error("validateCodeStatus was incorrect, should be true, result:", boolean)
	}
}

func TestValidateCodeStatusFalse(t *testing.T) {
	var a assertion
	a.Code = []int{200}

	boolean := a.validateCodeStatus(500)
	if boolean {
		t.Error("validateCodeStatus was incorrect, should be false, result:", boolean)
	}
}

func TestValidateCodeStatusEmpty(t *testing.T) {
	var a assertion
	a.Code = nil

	boolean := a.validateCodeStatus(200)
	if !boolean {
		t.Error("validateCodeStatus was incorrect, should be true, result:", boolean)
	}
}
