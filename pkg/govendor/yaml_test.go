package govendor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecYamlSerialization(t *testing.T) {
	dep := NewDependency("some-url", "some-branch")
	dep.Filters.
		AddExtension("dep-extension").
		AddTarget("dep-target").
		AddIgnore("dep-ignore")
	spec := NewSpec(&TestPreset{})
	spec.AddDependency(dep)

	data, err := toYaml(spec)
	assert.NoError(t, err)

	expected := `version: 0.2.0
preset: test-preset
vendor_dir: test-vendor-dir
extensions:
  - preset-extension
targets:
  - preset-target
ignores:
  - preset-ignore
deps:
  - url: some-url
    branch: some-branch
    extensions:
      - dep-extension
      - preset-extension-for-some-url
    targets:
      - dep-target
      - preset-target-for-some-url
    ignores:
      - dep-ignore
      - preset-ignore-for-some-url
`

	assert.Equal(t, expected, string(data))
}

func TestSpecLockYamlSerialization(t *testing.T) {
	dep := NewDependencyLock("some-url", "some-commit")
	spec := NewSpecLock(&TestPreset{})
	spec.AddDependencyLock(dep)

	data, err := toYaml(spec)
	assert.NoError(t, err)

	expected := `version: 0.2.0
deps:
  - url: some-url
    commit: some-commit
`

	assert.Equal(t, expected, string(data))
}
