package vendor

import (
	"crypto/sha256"
	"encoding/hex"
)

type Dependency struct {
	URL    string `yaml:"url"`
	Branch string `yaml:"branch"`

	Extensions []string `yaml:"extensions,omitempty"`
	Targets    []string `yaml:"targets,omitempty"`
	Ignores    []string `yaml:"ignores,omitempty"`
}

type DependencyLock struct {
	URL    string `yaml:"url"`
	Commit string `yaml:"commit"`
}

func NewDependency(url, branch string) *Dependency {
	return &Dependency{
		URL:        url,
		Branch:     branch,
		Extensions: []string{},
		Targets:    []string{},
		Ignores:    []string{},
	}
}

func NewDependencyLock(url, commit string) *DependencyLock {
	return &DependencyLock{
		URL:    url,
		Commit: commit,
	}
}

func (d *Dependency) ID() string {
	sha := sha256.New()
	data := sha.Sum([]byte(d.URL))
	return hex.EncodeToString(data)
}

func (d *Dependency) Update(other *Dependency) {
	d.URL = other.URL
	d.Branch = other.Branch
	d.Extensions = other.Extensions
	d.Targets = other.Targets
	d.Ignores = other.Ignores
}
