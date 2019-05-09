package common

import (
	"bytes"
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type TranslateOptions struct {
	Pretty bool
	Strict bool
}

type Common struct {
	Version string `yaml:"version"`
	Variant string `yaml:"variant"`
}

// Misc helpers
func Unmarshal(data []byte, to interface{}, strict bool) error {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(strict)
	return dec.Decode(to)
}

func Marshal(from interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(from, "", "  ")
	}
	return json.Marshal(from)
}
