package vending

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// SpecLock holds relevant information related to the specification of what
// specific versions need to be pinned when vendoring.
//
// This model directly maps to the serialized YAML of the spec lock file.
type SpecLock struct {
	Version string            `yaml:"version"`
	Deps    []*DependencyLock `yaml:"deps"`
	preset  Preset            `yaml:"-"`
}

// LoadSpecLock loads a SpecLock given a Preset.
func LoadSpecLock(preset Preset) (*SpecLock, error) {
	preset = checkPreset(preset, false)

	filename := preset.GetSpecLockFilename()
	data, err := os.ReadFile(filename)
	if err != nil {
		return NewSpecLock(preset), nil
	}

	spec := NewSpecLock(preset)
	err = yaml.Unmarshal(data, spec)
	if err != nil {
		return nil, fmt.Errorf("cannot read: %w", err)
	}

	return spec, nil
}

// NewSpecLock allocates a new SpecLock instance with a default initialization.
func NewSpecLock(preset Preset) *SpecLock {
	preset = checkPreset(preset, false)
	return &SpecLock{
		Version: VERSION,
		Deps:    []*DependencyLock{},
		preset:  preset,
	}
}

// AddDependencyLock adds a DependencyLock to the list of locked dependencies.
func (s *SpecLock) AddDependencyLock(lock *DependencyLock) {
	existing, ok := s.FindByURL(lock.URL)
	if ok {
		existing.Commit = lock.Commit
	} else {
		s.Deps = append(s.Deps, lock)
	}
}

// FindByURL finds a DependencyLock by URL.
func (s *SpecLock) FindByURL(url string) (*DependencyLock, bool) {
	for _, dep := range s.Deps {
		if strings.EqualFold(dep.URL, url) {
			return dep, true
		}
	}
	return nil, false
}

// Save converts the SpecLock to YAML, and writes the data in the spec lock
// file, as specified by the Preset.
func (s *SpecLock) Save() error {
	filename := s.preset.GetSpecLockFilename()
	data, err := toYaml(s)
	if err != nil {
		return fmt.Errorf("cannot save:  %w", err)
	}

	err = os.WriteFile(filename, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot save: %w", err)
	}

	return nil
}
