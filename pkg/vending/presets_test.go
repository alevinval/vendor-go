package vending

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultPreset(t *testing.T) {
	sut := DefaultPreset{}

	assert.Equal(t, "default", sut.GetPresetName())
	assert.Equal(t, ".vendor.yml", sut.GetSpecFilename())
	assert.Equal(t, ".vendor-lock.yml", sut.GetSpecLockFilename())
	assert.Equal(t, NewFilters(), sut.GetFilters())
	assert.Equal(t, NewFilters(), sut.GetFiltersForDependency(nil))

}

func TestDefaultPreset_GetCacheDir_UsesHome(t *testing.T) {
	os.Setenv("HOME", "some-home-path")

	sut := DefaultPreset{}

	assert.Equal(t, "some-home-path/.cache/vending", sut.GetCacheDir())
}

func TestDefaultPreset_GetCacheDir_DefaultsToTempDir(t *testing.T) {
	os.Clearenv()

	sut := DefaultPreset{}

	assert.Equal(t, "/tmp/.cache/vending", sut.GetCacheDir())
}
