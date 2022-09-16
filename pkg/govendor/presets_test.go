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
	assert.Empty(t, sut.GetTargets())
	assert.Empty(t, sut.GetTargets())
	assert.Empty(t, sut.GetDepExtensions(nil))
	assert.Empty(t, sut.GetDepTargets(nil))
	assert.Empty(t, sut.GetDepIgnores(nil))
}
