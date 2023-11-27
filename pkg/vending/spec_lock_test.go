package vending

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecLockAdd(t *testing.T) {
	sut := NewSpecLock(nil)
	assert.Empty(t, sut.Deps)

	depLock := NewDependencyLock("some-url", "some-commit")
	sut.AddDependencyLock(depLock)

	assert.Equal(t, []*DependencyLock{depLock}, sut.Deps)
}

func TestSpecLockAddUpdates(t *testing.T) {
	depLock := NewDependencyLock("some-url", "some-commit")
	sut := NewSpecLock(nil)
	sut.Deps = append(sut.Deps, depLock)

	assert.Equal(t, "some-commit", sut.Deps[0].Commit)

	other := NewDependencyLock("some-url", "other-commit")
	sut.AddDependencyLock(other)

	assert.Equal(t, "other-commit", sut.Deps[0].Commit)
}

func TestSpecLock_BumpsVersion_WhenSpecIsOlder(t *testing.T) {
	sut := NewSpecLock(&TestPreset{})
	sut.Version = "v0.0.1"

	sut.applyPreset(sut.preset)

	assert.Equal(t, VERSION, sut.Version)
}

func TestSpecLock_DoesNotBumpVersion_WhenSpecIsNewer(t *testing.T) {
	sut := NewSpecLock(&TestPreset{})
	sut.Version = "v999.0.0"

	sut.applyPreset(sut.preset)

	assert.Equal(t, "v999.0.0", sut.Version)
}

func TestSpecLock_SaveAndLoad(t *testing.T) {
	depLock := NewDependencyLock("some-url", "some-commit")
	expected := NewSpecLock(testPreset)
	expected.AddDependencyLock(depLock)

	err := expected.Save()
	assert.NoError(t, err)

	actual := NewSpecLock(testPreset)
	assert.NotEqual(t, expected, actual)

	actual.Load()
	assert.Equal(t, expected, actual)
}

func TestSpecLock_SaveOutput(t *testing.T) {
	depLock := NewDependencyLock("some-url", "some-commit")
	sut := NewSpecLock(testPreset)
	sut.AddDependencyLock(depLock)

	err := sut.Save()
	assert.NoError(t, err)

	expected := fmt.Sprintf(`version: %s
deps:
  - url: some-url
    commit: some-commit
`, VERSION)

	actual, err := os.ReadFile(testPreset.GetSpecLockFilename())
	assert.NoError(t, err)
	assert.Equal(t, expected, string(actual))
}
