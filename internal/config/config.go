package config

import (
	"fmt"
	"os"
	"sync"
)

type Instance struct {
	// ProjectID is the app-engine project to upload to.
	ProjectID string

	// Path is the filepath to this repository.
	//	Path string

	// Check is the very basic password you can use to set custom URLs.
	Check string

	// This is the default URL that requests to this app will redirect to.
	DefaultURL string

	// AddPageKey is the add URL form page.
	AdminPath string
}

var (
	instance *Instance
	once     sync.Once
)

func load() error {
	var s []string
	instance = &Instance{}
	getEnv := func(name string) string {
		v := os.Getenv(name)
		if v == "" {
			s = append(s, name)
		}
		return v
	}

	instance.ProjectID = getEnv("PROJECT_ID")
	instance.Check = getEnv("CHECK")
	instance.DefaultURL = getEnv("DEFAULT_URL")
	instance.AdminPath = getEnv("ADMIN_PATH")

	if len(s) > 0 {
		return fmt.Errorf("Required environment variables were unset: %v", s)
	}

	return nil
}

func Get() *Instance {
	once.Do(func() {
		load()
	})
	return instance
}
