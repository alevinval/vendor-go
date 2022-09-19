package vending

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ Preset = (*TestPreset)(nil)

type TestPreset struct {
	force bool
}

func (tp *TestPreset) GetPresetName() string {
	return "test-preset"
}

func (*TestPreset) GetVendorDir() string {
	return "test-vendor-dir"
}

func (tp *TestPreset) GetSpecFilename() string {
	return "some-spec-filename"
}

func (tp *TestPreset) GetSpecLockFilename() string {
	return "some-spec-lock-filename"
}

func (tp *TestPreset) GetFilters() *Filters {
	return NewFilters().
		AddExtension("preset-extension").
		AddTarget("preset-target").
		AddIgnore("preset-ignore")
}

func (tp *TestPreset) ForceFilters() bool {
	return tp.force
}

func (tp *TestPreset) GetCacheDir() string {
	return ".test-cache-dir"
}

func (tp *TestPreset) GetFiltersForDependency(dep *Dependency) *Filters {
	extension := fmt.Sprintf("preset-extension-for-%s", dep.URL)
	target := fmt.Sprintf("preset-target-for-%s", dep.URL)
	ignore := fmt.Sprintf("preset-ignore-for-%s", dep.URL)
	return NewFilters().
		AddExtension(extension).
		AddTarget(target).
		AddIgnore(ignore)

}

func TestNewSpec_LoadsPreset(t *testing.T) {
	spec := NewSpec(&TestPreset{})
	dep := NewDependency("some-url", "some-branch")
	spec.AddDependency(dep)

	expectedSpecFilters := NewFilters().
		AddExtension("preset-extension").
		AddTarget("preset-target").
		AddIgnore("preset-ignore")

	expectedDependencyFilters := NewFilters().
		AddExtension("preset-extension-for-some-url").
		AddTarget("preset-target-for-some-url").
		AddIgnore("preset-ignore-for-some-url")

	assert.Equal(t, &TestPreset{}, spec.preset)
	assert.Equal(t, "test-preset", spec.PresetName)
	assert.Equal(t, expectedSpecFilters, spec.Filters)
	assert.Equal(t, expectedDependencyFilters, spec.Deps[0].Filters)
}

func TestSpecAdd_AddsDeps(t *testing.T) {
	sut := NewSpec(&TestPreset{})
	assert.Empty(t, sut.Deps)

	dep := NewDependency("some-url", "some-branch")
	sut.AddDependency(dep)

	assert.Equal(t, []*Dependency{dep}, sut.Deps)
}

func TestSpecAdd_WhenDepAlreadyPresent_UpdatesExistingDep(t *testing.T) {
	dep := NewDependency("some-url", "some-branch")
	dep.Filters = NewFilters().
		AddExtension("to-be-replaced").
		AddTarget("to-be-replaced").
		AddIgnore("to-be-replaced")

	sut := NewSpec(nil)
	sut.Deps = append(sut.Deps, dep)

	other := NewDependency("some-url", "other-branch")
	other.Filters = NewFilters().
		AddExtension("other-extension").
		AddTarget("other-target").
		AddIgnore("other-ignore")
	sut.AddDependency(other)

	assert.Equal(t, []*Dependency{dep}, sut.Deps)
	assert.Equal(t, other.Filters, sut.Deps[0].Filters)
}

func TestSpec_WhenForceFilters_OverridesFilters(t *testing.T) {
	preset := &TestPreset{true}

	dep := NewDependency("some-url", "some-branch")
	dep.Filters = NewFilters().
		AddExtension("to-be-replaced").
		AddTarget("to-be-replaced").
		AddIgnore("to-be-replaced")

	sut := NewSpec(nil)
	sut.Filters = NewFilters().
		AddExtension("to-be-replaced").
		AddTarget("to-be-replaced").
		AddIgnore("to-be-replaced")
	sut.AddDependency(dep)

	sut.applyPreset(preset)

	assert.Equal(t, preset.GetFilters(), sut.Filters)
	assert.Equal(t, preset.GetFiltersForDependency(dep), sut.Deps[0].Filters)
}

func TestSpec_BumpsVersion_WhenSpecIsOlder(t *testing.T) {
	spec := NewSpec(nil)
	spec.Version = "v0.0.1"

	spec.applyPreset(&TestPreset{})

	assert.Equal(t, VERSION, spec.Version)
}

func TestSpec_DoesNotBumpVersion_WhenSpecIsNewer(t *testing.T) {
	spec := NewSpec(nil)
	spec.Version = "v999.0.0"

	spec.applyPreset(&TestPreset{})

	assert.Equal(t, "v999.0.0", spec.Version)
}
