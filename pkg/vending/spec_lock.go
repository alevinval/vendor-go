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

// Load SpecLock from the filesystem.
func (s *SpecLock) Load() error {
	preset := s.preset

	filename := preset.GetSpecLockFilename()
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	err = yaml.Unmarshal(data, s)
	if err != nil {
		return fmt.Errorf("cannot unmarshal: %w", err)
	}

	s.applyPreset(preset)
	return nil
}

// Save converts the SpecLock to YAML, and writes the data in the spec lock
// file, as specified by the Preset.
func (s *SpecLock) Save() error {
	filename := s.preset.GetSpecLockFilename()
	data, err := toYaml(s)
	if err != nil {
		return fmt.Errorf("cannot convert to yaml:  %w", err)
	}

	err = os.WriteFile(filename, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot write file: %w", err)
	}

	return nil
}

func (s *SpecLock) applyPreset(preset Preset) {
	s.preset = preset

	if s.Version < VERSION {
		s.Version = VERSION
	}
}
