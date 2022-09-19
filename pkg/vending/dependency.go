package vending

import (
	"crypto/sha256"
	"encoding/hex"
)

// Dependency holds relevant information related to a dependency that has to be
// vendored. This model directly maps to the serialized YAML, for dependencies.
type Dependency struct {
	URL     string   `yaml:"url"`
	Branch  string   `yaml:"branch"`
	Filters *Filters `yaml:",inline"`
}

// DependencyLock holds relevant information of a dependency that has been
// locked to a specific commit. This model directly maps to the serialized YAML
// for locked dependencies.
type DependencyLock struct {
	URL    string `yaml:"url"`
	Commit string `yaml:"commit"`
}

// NewDependency allocates a Dependency, with a default Filters instance.
func NewDependency(url, branch string) *Dependency {
	return &Dependency{
		URL:     url,
		Branch:  branch,
		Filters: NewFilters(),
	}
}

// NewDependencyLock allocates a new DependencyLock
func NewDependencyLock(url, commit string) *DependencyLock {
	return &DependencyLock{
		URL:    url,
		Commit: commit,
	}
}

// ID returns the unique identifier for that dependency.
// It returns the sha256 of the URL.
func (d *Dependency) ID() string {
	sha := sha256.New()
	data := sha.Sum([]byte(d.URL))
	return hex.EncodeToString(data)
}

// Update changes the URL, Branch and Filters fields of the dependency by the
// fields of another one. This clones the Filters to ensure there's no shared
// data with the other Dependency.
func (d *Dependency) Update(other *Dependency) {
	d.URL = other.URL
	d.Branch = other.Branch
	d.Filters = other.Filters.Clone()
}
