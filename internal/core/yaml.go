package core

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

func toYaml(obj interface{}) []byte {
	b := &bytes.Buffer{}
	encoder := yaml.NewEncoder(b)
	encoder.SetIndent(2)
	encoder.Encode(obj)
	encoder.Close()
	return b.Bytes()
}
