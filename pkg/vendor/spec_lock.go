package vendor

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type SpecLock struct {
	Version string            `yaml:"version"`
	Deps    []*DependencyLock `yaml:"deps"`
	preset  Preset            `yaml:"-"`
}

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

func NewSpecLock(preset Preset) *SpecLock {
	preset = checkPreset(preset, false)
	return &SpecLock{
		Version: VERSION,
		Deps:    []*DependencyLock{},
		preset:  preset,
	}
}

func (s *SpecLock) AddDependencyLock(lock *DependencyLock) {
	existing, ok := s.Find(lock.URL)
	if ok {
		existing.Commit = lock.Commit
	} else {
		s.Deps = append(s.Deps, lock)
	}
}

func (s *SpecLock) Find(url string) (*DependencyLock, bool) {
	for _, dep := range s.Deps {
		if strings.EqualFold(dep.URL, url) {
			return dep, true
		}
	}
	return nil, false
}

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
