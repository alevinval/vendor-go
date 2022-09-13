package core

type Preset interface {
	GetExtensions() []string
	GetTargets(dep *Dependency) []string
	GetIgnores(dep *Dependency) []string
}
