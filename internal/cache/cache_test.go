package cache

import (
	"io/fs"
	"os"
	"testing"

	"github.com/alevinval/vendor-go/pkg/vending"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testCachePath = ".test-cache-dir"

var _ fsOps = (*fsMock)(nil)

func TestCache_Reset(t *testing.T) {
	fsMock := &fsMock{}
	fsMock.
		On("RemoveAll", ".test-cache-dir").Return(nil).
		On("MkdirAll", ".test-cache-dir/locks", os.ModePerm).Return(nil).
		On("MkdirAll", ".test-cache-dir/repos", os.ModePerm).Return(nil)

	sut := New(testCachePath)
	sut.fs = fsMock

	sut.Reset()

	fsMock.AssertExpectations(t)
}

func TestCache_Init(t *testing.T) {
	fsMock := &fsMock{}
	fsMock.
		On("MkdirAll", ".test-cache-dir/locks", os.ModePerm).Return(nil).
		On("MkdirAll", ".test-cache-dir/repos", os.ModePerm).Return(nil)

	sut := New(testCachePath)
	sut.fs = fsMock

	sut.Init()

	fsMock.AssertExpectations(t)
}

func TestCache_Clean(t *testing.T) {
	fsMock := &fsMock{}
	fsMock.
		On("RemoveAll", testCachePath).Return(nil).
		On("MkdirAll", ".test-cache-dir/locks", os.ModePerm).Return(nil).
		On("MkdirAll", ".test-cache-dir/repos", os.ModePerm).Return(nil)

	sut := New(testCachePath)
	sut.fs = fsMock

	sut.Reset()
}

func TestCache_Lock(t *testing.T) {
	defer func() {
		os.RemoveAll(testCachePath)
	}()

	sut := New(testCachePath)

	lock, err := sut.Lock()

	assert.NoError(t, err)
	assert.NotNil(t, lock)
}

func TestCache_GetRepositoryPath(t *testing.T) {
	dep := vending.NewDependency("some-url", "some-branch")
	sut := New(testCachePath)

	actual := sut.getRepositoryPath(dep)

	assert.Equal(t,
		".test-cache-dir/repos/f28318b204791d282d65cc09bba5389e8b9c7406",
		actual,
	)
}

func TestCache_GetRepositoryLockPath(t *testing.T) {
	dep := vending.NewDependency("some-url", "some-branch")
	sut := New(testCachePath)

	actual := sut.getRepositoryLockPath(dep)

	assert.Equal(t,
		".test-cache-dir/locks/f28318b204791d282d65cc09bba5389e8b9c7406",
		actual,
	)
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
