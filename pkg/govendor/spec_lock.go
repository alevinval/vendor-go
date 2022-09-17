package govendor

import (
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

	fileName := preset.GetSpecLockFilename()
	data, err := os.ReadFile(fileName)
	if err != nil {
		return NewSpecLock(preset), nil
	}

	spec := NewSpecLock(preset)
	err = yaml.Unmarshal(data, spec)
	if err != nil {
		logger.Errorf("cannot read %s: %s", fileName, err)
		return nil, err
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

func (s *SpecLock) Add(lock *DependencyLock) {
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
	data, err := toYaml(s)
	if err != nil {
		return err
	}

	return os.WriteFile(s.preset.GetSpecLockFilename(), data, os.ModePerm)
}
