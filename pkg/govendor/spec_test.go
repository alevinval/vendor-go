package govendor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ Preset = (*TestPreset)(nil)

type TestPreset struct {
	*DefaultPreset
}

func (tp *TestPreset) GetExtensions() []string {
	return []string{"md", "java"}
}

func (tp *TestPreset) GetTargets(dep *Dependency) []string {
	target := fmt.Sprintf("preset-target-for-%s", dep.URL)
	return []string{target}
}

func (tp *TestPreset) GetIgnores(dep *Dependency) []string {
	ignore := fmt.Sprintf("preset-ignore-for-%s", dep.URL)
	return []string{ignore}
}

func TestNewSpec_LoadsPreset(t *testing.T) {
	spec := NewSpec(&TestPreset{})

	assert.Equal(t, ".vendor.yml", spec.Preset.GetSpecFilename())
	assert.Equal(t, ".vendor-lock.yml", spec.Preset.GetSpecLockFilename())
	assert.Equal(t, []string{"java", "md"}, spec.Extensions)
	assert.Equal(t, []string{}, spec.Targets)
	assert.Equal(t, []string{}, spec.Ignores)
}

func TestSpecAdd_AddsDeps(t *testing.T) {
	sut := NewSpec(nil)
	assert.Empty(t, sut.Deps)

	dep := NewDependency("some-url", "some-branch")
	sut.Add(dep)

	assert.Equal(t, []*Dependency{dep}, sut.Deps)
}

func TestSpecAdd_WhenDepAlreadyPresent_UpdatesExistingDep(t *testing.T) {
	dep := NewDependency("some-url", "some-branch")
	dep.Targets = []string{"to-be-overwritten-target"}
	dep.Ignores = []string{"to-be-overwritten-ignore"}

	sut := NewSpec(&TestPreset{})
	sut.Deps = append(sut.Deps, dep)

	other := NewDependency("some-url", "other-branch")
	other.Targets = []string{"other-target"}
	other.Ignores = []string{"other-ignore"}
	sut.Add(other)

	assert.Equal(t, []*Dependency{dep}, sut.Deps)
	assert.Equal(t, []string{"other-target", "preset-target-for-some-url"}, sut.Deps[0].Targets)
	assert.Equal(t, []string{"other-ignore", "preset-ignore-for-some-url"}, sut.Deps[0].Ignores)
}
