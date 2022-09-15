package govendor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDependencyUpdate(t *testing.T) {
	dep := NewDependency("some-url", "some-branch")

	assert.Empty(t, dep.Targets)
	assert.Empty(t, dep.Ignores)
	assert.Empty(t, dep.Extensions)

	other := NewDependency("some-url", "some-branch")
	other.Targets = []string{"some-targets"}
	other.Ignores = []string{"some-ignores"}
	other.Extensions = []string{"some-extensions"}

	dep.Update(other)

	assert.Equal(t, other.Targets, dep.Targets)
	assert.Equal(t, other.Ignores, dep.Ignores)
	assert.Equal(t, other.Extensions, dep.Extensions)
}

func TestDependencyID(t *testing.T) {
	dep := NewDependency("some-url", "some-branch")
	assert.Equal(t, "736f6d652d75726ce3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", dep.ID())
}
