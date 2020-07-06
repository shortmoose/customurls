package config

type Instance struct {
	// ProjectID is the app-engine project to upload to.
	ProjectID string

	// This is the default URL that requests to this app will redirect to.
	DefaultURL string
}
