package util

import (
	"testing"
)

func TestGetKey(t *testing.T) {
	v := GetKey("abcDEF")
	if v != "abcdef" {
		t.Error("invalid value")
	}

	v = GetKey("abcdef")
	if v != "abcdef" {
		t.Error("invalid value")
	}

	v = GetKey("abc-ef")
	if v != "abc-ef" {
		t.Error("invalid value")
	}

	v = GetKey("abc_def")
	if v != "" {
		t.Error("invalid value")
	}

	v = GetKey("abc?sdf")
	if v != "" {
		t.Error("invalid value")
	}
}
