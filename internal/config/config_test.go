package config

import (
	"fmt"
	"os"
	"testing"
)

// TODO: We should test error cases here as well.
func TestLoad(t *testing.T) {
	v := []string{"PROJECT_ID", "CHECK", "DEFAULT_URL", "ADMIN_PATH"}
	for i, s := range v {
		os.Setenv(s, fmt.Sprintf("%d", i))
	}
	if err := load(); err != nil {
		t.Error("Fail:", err)
	}
	for _, s := range v {
		os.Unsetenv(s)
	}

	if instance.ProjectID != "0" {
		t.Error("ProjectID has incorrect value", instance.ProjectID)
	}
	if instance.Check != "1" {
		t.Error("Check has incorrect value")
	}
	if instance.DefaultURL != "2" {
		t.Error("DefaultURL has incorrect value")
	}
	if instance.AdminPath != "3" {
		t.Error("AdminPath has incorrect value")
	}
}
