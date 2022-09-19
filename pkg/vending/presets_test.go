package vending

import (
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
