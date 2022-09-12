package core

import (
	"fmt"
	"strings"
)

var presets = map[string]Preset{
	"skyscanner": &SkyscannerPreset{},
}

type Preset interface {
	GetExtensions() []string
	GetTargets(dep *Dependency) []string
}

type SkyscannerPreset struct{}

func (sp *SkyscannerPreset) GetExtensions() []string {
	return []string{"proto"}
}

func (sp *SkyscannerPreset) GetTargets(dep *Dependency) []string {
	return []string{fmt.Sprintf("proto/%s", extractGithubGroup(dep))}
}

func extractGithubGroup(dep *Dependency) string {
	repo := strings.Split(dep.URL, ":")[1]
	group := strings.Split(repo, "/")[0]
	return strings.ReplaceAll(group, "-", "")
}
