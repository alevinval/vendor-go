package control

import (
	"fmt"
	"os"

	"github.com/alevinval/vendor-go/internal/cache"
	"github.com/alevinval/vendor-go/internal/installer"
	"github.com/alevinval/vendor-go/pkg/log"
	"github.com/alevinval/vendor-go/pkg/vending"
	"github.com/fatih/color"
)

// Option is used to apply customizations to the Controller.
type Option = func(c *Controller)

func WithPreset(preset vending.Preset) Option {
	return func(c *Controller) {
		c.preset = preset
	}
}

// Controller knows how to execute the main business logic of the tool.
type Controller struct {
	preset vending.Preset
	cache  *cache.Cache
}

// New allocates a command controller based on the provided options.
// The orchestrator contains the core logic of the operations that the vending
// tool supports. This can be used to automate vending operations without having
// to rely on running cobra commands.
func New(opts ...Option) *Controller {
	c := &Controller{
		preset: &vending.DefaultPreset{},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.cache = cache.New(c.preset.GetCacheDir())

	return c
}

// Init initializes the vending tool for the current directory. This creates a
// default spec in the filesystem.
func (c *Controller) Init() error {
	_, err := os.ReadFile(c.preset.GetSpecFilename())
	if err == nil {
		return fmt.Errorf("%q already exists? %w", c.preset.GetSpecFilename(), err)
	}

	spec := vending.NewSpec(c.preset)

	err = spec.Save()
	if err != nil {
		return fmt.Errorf("cannot save spec: %w", err)
	}

	log.S().Infof("%s has been created ✅", c.preset.GetSpecFilename())
	return nil
}

// Install vendors the dependencies at the version specified by the lockfile.
// When no lockfile is present, it locks the dependencies at the latest
// reference of the branch that the spec defines for each dependency.
func (c *Controller) Install() error {
	lock, err := c.cache.Lock()
	if err != nil {
		return fmt.Errorf("cannot lock cache: %w", err)
	}
	defer lock.Release()

	spec := vending.NewSpec(c.preset)
	if err := spec.Load(); err != nil {
		return fmt.Errorf("cannot load spec: %w", err)
	}

	specLock := vending.NewSpecLock(c.preset)
	if err := specLock.Load(); err != nil {
		return fmt.Errorf("cannot load speclock: %w", err)
	}

	cacheDir := c.preset.GetCacheDir()
	log.S().Infof("repository cache located at %s",
		color.MagentaString(cacheDir),
	)

	ins := installer.New(c.cache, spec, specLock)
	if err := ins.Install(); err != nil {
		return fmt.Errorf("cannot install: %w", err)
	}

	if err := spec.Save(); err != nil {
		return fmt.Errorf("cannot save spec: %w", err)
	}

	if err := specLock.Save(); err != nil {
		return fmt.Errorf("cannot save speclock: %w", err)
	}

	log.S().Infof("install success ✅")
	return nil
}

// Update vendors the dependencies at the latest reference from the specified
// branch, this updates the lockfile with the locked references for each
// dependency.
func (c *Controller) Update() error {
	lock, err := c.cache.Lock()
	if err != nil {
		return fmt.Errorf("cannot lock cache: %w", err)
	}
	defer lock.Release()

	spec := vending.NewSpec(c.preset)
	if err := spec.Load(); err != nil {
		return fmt.Errorf("cannot load spec: %w", err)
	}

	specLock := vending.NewSpecLock(c.preset)
	if err := specLock.Load(); err != nil {
		return fmt.Errorf("cannot load speclock: %w", err)
	}

	cacheDir := c.preset.GetCacheDir()
	log.S().Infof("repository cache located at %s",
		color.MagentaString(cacheDir),
	)

	ins := installer.New(c.cache, spec, specLock)
	if err := ins.Update(); err != nil {
		return fmt.Errorf("cannot update: %w", err)
	}

	if err := spec.Save(); err != nil {
		return fmt.Errorf("cannot save spec: %w", err)
	}

	if err := specLock.Save(); err != nil {
		return fmt.Errorf("cannot save speclock: %w", err)
	}

	log.S().Infof("update success ✅")
	return nil
}

// AddDependency adds a new dependency into the spec file.
func (c *Controller) AddDependency(url, branch string, filters *vending.Filters) error {
	spec := vending.NewSpec(c.preset)
	if err := spec.Load(); err != nil {
		return fmt.Errorf("cannot load spec: %w", err)
	}

	dep := vending.NewDependency(url, branch)
	dep.Filters.ApplyFilters(filters)

	spec.AddDependency(dep)

	if err := spec.Save(); err != nil {
		return fmt.Errorf("cannot save spec: %w", err)
	}

	log.S().Infof("added dependency %s@%s ✅",
		color.CyanString(url),
		color.YellowString(branch),
	)
	return nil
}

// CleanCache performs a reset of the repository cache, once cleaned, the
// repositories of the dependencies will have to be cloned again.
func (c *Controller) CleanCache() error {
	if err := c.cache.Reset(); err != nil {
		return fmt.Errorf("cannot clean cache: %w", err)
	}

	log.S().Infof("cache has been cleaned ✅")
	return nil
}
