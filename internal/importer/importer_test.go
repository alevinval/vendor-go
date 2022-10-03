package importer

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/alevinval/vendor-go/internal/git"
	"github.com/alevinval/vendor-go/internal/lock"
	"github.com/alevinval/vendor-go/pkg/log"
	"github.com/alevinval/vendor-go/pkg/vending"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func init() {
	log.Level.SetLevel(zapcore.DebugLevel)
}

const VENDOR_DIR = ".testvendor"
const INPUT_DIR = ".testinput"

func TestImporter_Import_VendorsFiles(t *testing.T) {
	filters := vending.NewFilters().
		AddExtension("txt").
		AddTarget("root.txt", "target").
		AddIgnore("ignored", "target/nested_ignored")

	filepaths := []string{
		"root.txt",
		"ignored/ignored.txt",
		"target/target.txt",
		"target/nested/target.txt",
		"target/nested_ignored/ignored.txt",
	}

	sut := setUp(t, filters, filepaths)
	defer cleanUp(t)

	sut.Import()

	assertExists(t, vendorPath("root.txt"))
	assertExists(t, vendorPath("target/target.txt"))
	assertExists(t, vendorPath("target/nested/target.txt"))

	assertNotExists(t, vendorPath("target/nested_ignored/ignored.txt"))
}

func assertExists(t *testing.T, filepath string) {
	_, err := os.Stat(filepath)
	assert.NoError(t, err, fmt.Sprintf("%q should exist, but it does not", filepath))
}

func assertNotExists(t *testing.T, filepath string) {
	_, err := os.Stat(filepath)
	assert.True(t, os.IsNotExist(err),
		fmt.Sprintf("%q should not exist, but it does", filepath))
}

func vendorPath(filepath string) string {
	return path.Join(VENDOR_DIR, filepath)
}

func setUp(t *testing.T, filters *vending.Filters, filepaths []string) *Importer {
	os.MkdirAll(INPUT_DIR, os.ModePerm)
	for _, filepath := range filepaths {
		filepath := path.Join(INPUT_DIR, filepath)
		os.MkdirAll(path.Dir(filepath), os.ModePerm)
		_, err := os.Create(filepath)
		assert.NoError(t, err)
	}

	spec := vending.NewSpec(nil)
	spec.VendorDir = VENDOR_DIR
	spec.Filters = filters

	lock := lock.New(".testlock")
	dep := vending.NewDependency("some-url", "some-branch")
	repo := git.NewRepository(INPUT_DIR, lock, dep)
	return New(repo, spec, dep)
}

func cleanUp(t *testing.T) {
	os.RemoveAll(VENDOR_DIR)
	os.RemoveAll(INPUT_DIR)
}
