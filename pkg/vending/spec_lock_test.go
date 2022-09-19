package vending

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecLockAdd(t *testing.T) {
	sut := NewSpecLock(nil)
	assert.Empty(t, sut.Deps)

	lock := NewDependencyLock("some-url", "some-commit")
	sut.AddDependencyLock(lock)

	assert.Equal(t, []*DependencyLock{lock}, sut.Deps)
}

func TestSpecLockAddUpdates(t *testing.T) {
	lock := NewDependencyLock("some-url", "some-commit")
	sut := NewSpecLock(nil)
	sut.Deps = append(sut.Deps, lock)

	assert.Equal(t, "some-commit", sut.Deps[0].Commit)

	other := NewDependencyLock("some-url", "other-commit")
	sut.AddDependencyLock(other)

	assert.Equal(t, "other-commit", sut.Deps[0].Commit)
}
