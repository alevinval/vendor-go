package cache

import (
	"io/fs"
	"os"
	"testing"

	"github.com/alevinval/vendor-go/pkg/vending"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var _ CacheFS = (*fsMock)(nil)
var _ vending.Preset = (*TestPreset)(nil)

func TestManager_Reset(t *testing.T) {
	fsMock := &fsMock{}
	fsMock.
		On("RemoveAll", ".test-cache-dir").Return(nil).
		On("MkdirAll", ".test-cache-dir/locks", os.ModePerm).Return(nil).
		On("MkdirAll", ".test-cache-dir/repos", os.ModePerm).Return(nil)

	sut := NewManager(&TestPreset{})
	sut.fs = fsMock

	sut.Reset()

	fsMock.AssertExpectations(t)
}

func TestManager_Ensure(t *testing.T) {
	fsMock := &fsMock{}
	fsMock.
		On("MkdirAll", ".test-cache-dir/locks", os.ModePerm).Return(nil).
		On("MkdirAll", ".test-cache-dir/repos", os.ModePerm).Return(nil)

	sut := NewManager(&TestPreset{})
	sut.fs = fsMock

	sut.Ensure()

	fsMock.AssertExpectations(t)
}

func TestManager_LockCache(t *testing.T) {
	preset := &TestPreset{}
	defer func() {
		os.RemoveAll(preset.GetCacheDir())
	}()

	sut := NewManager(&TestPreset{})

	lock, err := sut.LockCache()

	assert.NoError(t, err)
	assert.NotNil(t, lock)
}

func TestManager_GetRepositoryPath(t *testing.T) {
	dep := vending.NewDependency("some-url", "some-branch")
	sut := NewManager(&TestPreset{})

	actual := sut.GetRepositoryPath(dep)

	assert.Equal(t,
		".test-cache-dir/repos/736f6d652d75726ce3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		actual,
	)
}

func TestManager_GetRepositoryLockPath(t *testing.T) {
	dep := vending.NewDependency("some-url", "some-branch")
	sut := NewManager(&TestPreset{})

	actual := sut.GetRepositoryLockPath(dep)

	assert.Equal(t,
		".test-cache-dir/locks/736f6d652d75726ce3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		actual,
	)
}

type TestPreset struct {
	vending.DefaultPreset
}

func (tp *TestPreset) GetCacheDir() string {
	return ".test-cache-dir"
}

type fsMock struct {
	mock.Mock
}

func (m *fsMock) MkdirAll(path string, perm fs.FileMode) error {
	args := m.Called(path, perm)
	return args.Error(0)
}

func (m *fsMock) RemoveAll(path string) error {
	args := m.Called(path)
	return args.Error(0)
}
