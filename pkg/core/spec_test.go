package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecAdd(t *testing.T) {
	sut := NewSpec()
	assert.Empty(t, sut.Deps)

	dep := NewDependency("some-url", "some-branch")
	sut.Add(dep)

	assert.Equal(t, []*Dependency{dep}, sut.Deps)
}

func TestSpecAddUpdates(t *testing.T) {
	dep := NewDependency("some-url", "some-branch")
	sut := NewSpec()
	sut.Deps = append(sut.Deps, dep)

	assert.Equal(t, []string{}, sut.Deps[0].Targets)

	other := NewDependency("some-url", "other-branch")
	other.Targets = []string{"other-target"}
	sut.Add(other)

	assert.Equal(t, []*Dependency{dep}, sut.Deps)
	assert.Equal(t, []string{"other-target"}, sut.Deps[0].Targets)
}

func TestSpecToYaml(t *testing.T) {
	sut := NewSpec()
	assert.Equal(t, "version: 0.0.1\nvendor_dir: vendor/\ndeps: []\n", string(sut.ToYaml()))
}
