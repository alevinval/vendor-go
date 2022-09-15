package govendor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultPreset(t *testing.T) {
	sut := DefaultPreset{}

	assert.Equal(t, "default", sut.GetPresetName())
	assert.Equal(t, ".vendor.yml", sut.GetSpecFilename())
	assert.Equal(t, ".vendor-lock.yml", sut.GetSpecLockFilename())
	assert.Empty(t, sut.GetExtensions())
	assert.Empty(t, sut.GetTargets(nil))
	assert.Empty(t, sut.GetIgnores(nil))
}
