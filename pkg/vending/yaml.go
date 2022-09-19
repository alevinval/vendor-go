package vending

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

func toYaml(obj interface{}) ([]byte, error) {
	b := &bytes.Buffer{}
	encoder := yaml.NewEncoder(b)
	encoder.SetIndent(2)
	err := encoder.Encode(obj)
	if err != nil {
		return nil, err
	}
	encoder.Close()
	return b.Bytes(), nil
}
