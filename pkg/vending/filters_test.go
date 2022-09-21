package vending

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFilters_IsEmptyByDefault(t *testing.T) {
	sut := NewFilters()

	assert.Empty(t, sut.Extensions)
	assert.Empty(t, sut.Targets)
	assert.Empty(t, sut.Ignores)
}

func TestFilters_AddVariadic_AppendsUnderlyingArray(t *testing.T) {
	sut := NewFilters().
		AddExtension("ext-1", "ext-2").
		AddTarget("target-1", "target-2").
		AddIgnore("ignore-1", "ignore-2")

	assert.Equal(t, []string{"ext-1", "ext-2"}, sut.Extensions)
	assert.Equal(t, []string{"target-1", "target-2"}, sut.Targets)
	assert.Equal(t, []string{"ignore-1", "ignore-2"}, sut.Ignores)
}

func TestFilters_ApplyPreset_UpdatesWithUnion(t *testing.T) {
	sut := NewFilters().
		AddExtension("ext-1").
		AddTarget("target-1").
		AddIgnore("ignore-1")

	sut.ApplyPreset(&TestPreset{})

	assert.Equal(t, []string{"ext-1", "preset-extension"}, sut.Extensions)
	assert.Equal(t, []string{"preset-target", "target-1"}, sut.Targets)
	assert.Equal(t, []string{"ignore-1", "preset-ignore"}, sut.Ignores)
}

func TestFilters_ApplyPresetDep_UpdatesWithUnion(t *testing.T) {
	sut := NewFilters().
		AddExtension("ext-1").
		AddTarget("target-1").
		AddIgnore("ignore-1")

	dep := NewDependency("some-url", "some-branch")
	dep.Filters = NewFilters().
		AddExtension("ext-2").
		AddTarget("target-2").
		AddIgnore("ignore-2")

	sut.ApplyDep(&TestPreset{}, dep)

	assert.Equal(t, []string{"ext-1", "ext-2", "preset-extension-for-some-url"}, sut.Extensions)
	assert.Equal(t, []string{"preset-target-for-some-url", "target-1", "target-2"}, sut.Targets)
	assert.Equal(t, []string{"ignore-1", "ignore-2", "preset-ignore-for-some-url"}, sut.Ignores)
}

func TestFilters_Clone_ReturnsNewInstance(t *testing.T) {
	sut := NewFilters().
		AddExtension("ext-1").
		AddTarget("target-1").
		AddIgnore("ignore-1")

	clone := sut.Clone()
	assert.Equal(t, sut, clone)

	clone.AddExtension("ext-2")
	assert.NotEqual(t, sut, clone)
}

func TestFilters_Clear_ThenClearsLists(t *testing.T) {
	sut := NewFilters().
		AddExtension("ext-1").
		AddTarget("target-1").
		AddIgnore("ignore-1")

	assert.NotEqual(t, NewFilters(), sut)

	sut.Clear()

	assert.Equal(t, NewFilters(), sut)
}
