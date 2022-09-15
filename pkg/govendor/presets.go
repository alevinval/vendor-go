package govendor

type Preset interface {
	GetSpecFilename() string
	GetSpecLockFilename() string
	GetExtensions() []string
	GetTargets(dep *Dependency) []string
	GetIgnores(dep *Dependency) []string
}
