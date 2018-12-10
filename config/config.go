package config

// Config represents a grouping of environments, stores and packages
type Config struct {
	Stores    []Store            `yaml:"stores"`
	Processes []Process          `yaml:"processes"`
	Imports   []PackageReference `yaml:"imports"`
}

// Store represents a configuration store
type Store struct {
	Name      string   `yaml:"name"`
	StoreType string   `yaml:"type"`
	Stores    []string `yaml:"stores"`
	Feeds     []Feed   `yaml:"feeds"`

	Params map[string]string `yaml:"params"`
}

// Process represents a process under the given environment
type Process struct {
	Name   string            `yaml:"name"`
	Stores []string          `yaml:"stores"`
	Path   string            `yaml:"path"`
	Args   []string          `yaml:"args"`
	Vars   map[string]string `yaml:"env"`
}

// PackageReference represents a reference to a package
type PackageReference struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// Package represents a versioned artifiact
type Package struct {
	Name    string   `yaml:"name"`
	Version string   `yaml:"version"`
	Targets []Target `yaml:"targets"`
}

// Target repesents an install target
type Target struct {
	Platform     string `yaml:"platform"`
	Architecture string `yaml:"architecture"`
	Name         string `yaml:"name"`

	Tasks []interface{} `yaml:"tasks"`
}

// Feed denotes a package feed
type Feed struct {
	Name string `yaml:"name"`
	URI  string `yaml:"uri"`
}
