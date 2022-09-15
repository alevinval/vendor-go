package govendor

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type SpecLock struct {
	Version string            `yaml:"version"`
	Deps    []*DependencyLock `yaml:"deps"`
	Preset  Preset            `yaml:"-"`
}

func LoadSpecLock(preset Preset) (*SpecLock, error) {
	fileName := preset.GetSpecLockFilename()
	data, err := os.ReadFile(fileName)
	if err != nil {
		logger.Errorf("cannot read %s: %s", fileName, err)
		return nil, err
	}

	spec := &SpecLock{}
	err = yaml.Unmarshal(data, spec)
	if err != nil {
		logger.Errorf("cannot read %s: %s", fileName, err)
		return nil, err
	}

	spec.Preset = preset
	return spec, nil
}

func NewSpecLock(preset Preset) *SpecLock {
	return &SpecLock{
		Version: VERSION,
		Deps:    []*DependencyLock{},
		Preset:  preset,
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

	return os.WriteFile(s.Preset.GetSpecLockFilename(), data, os.ModePerm)
}
