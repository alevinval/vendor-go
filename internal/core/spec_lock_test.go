package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecLockAdd(t *testing.T) {
	sut := NewSpecLock()
	assert.Empty(t, sut.Deps)

	lock := NewDependencyLock("some-url", "some-commit")
	sut.Add(lock)

	assert.Equal(t, []*DependencyLock{lock}, sut.Deps)
}

func TestSpecLockAddUpdates(t *testing.T) {
	lock := NewDependencyLock("some-url", "some-commit")
	sut := NewSpecLock()
	sut.Deps = append(sut.Deps, lock)

	assert.Equal(t, "some-commit", sut.Deps[0].Commit)

	other := NewDependencyLock("some-url", "other-commit")
	sut.Add(other)

	assert.Equal(t, "other-commit", sut.Deps[0].Commit)
}

func TestSpecLockToYaml(t *testing.T) {
	sut := NewSpecLock()
	assert.Equal(t, "version: 0.0.1\ndeps: []\n", string(sut.ToYaml()))
}
