package govendor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ Preset = (*TestPreset)(nil)

type TestPreset struct{}

func (tp *TestPreset) GetPresetName() string {
	return "test-preset"
}

func (tp *TestPreset) GetSpecFilename() string {
	return "some-spec-filename"
}

func (tp *TestPreset) GetSpecLockFilename() string {
	return "some-spec-lock-filename"
}

func (tp *TestPreset) GetExtensions() []string {
	return []string{"preset-extension"}
}

func (tp *TestPreset) GetTargets() []string {
	return []string{"preset-target"}
}

func (tp *TestPreset) GetIgnores() []string {
	return []string{"preset-ignore"}
}

func (tp *TestPreset) GetDepExtensions(dep *Dependency) []string {
	target := fmt.Sprintf("preset-extension-for-%s", dep.URL)
	return []string{target}
}

func (tp *TestPreset) GetDepTargets(dep *Dependency) []string {
	target := fmt.Sprintf("preset-target-for-%s", dep.URL)
	return []string{target}
}

func (tp *TestPreset) GetDepIgnores(dep *Dependency) []string {
	ignore := fmt.Sprintf("preset-ignore-for-%s", dep.URL)
	return []string{ignore}
}

func TestNewSpec_LoadsPreset(t *testing.T) {
	spec := NewSpec(&TestPreset{})

	assert.Equal(t, &TestPreset{}, spec.preset)
	assert.Equal(t, []string{"preset-extension"}, spec.Extensions)
	assert.Equal(t, []string{"preset-target"}, spec.Targets)
	assert.Equal(t, []string{"preset-ignore"}, spec.Ignores)
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
