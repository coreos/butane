package config

import (
	"errors"
	"fmt"

	json "github.com/ajeddeloh/go-json"
	"github.com/go-yaml/yaml"
	"github.com/coreos/go-semver/semver"
)

var (
	ErrNoVariant      = errors.New("Error parsing variant. Variant must be specified")
	ErrInvalidVersion = errors.New("Error parsing version. Version must be a valid semver")

	registry = map[string]translator{
		"fcos+1.0.0": TranslateFcos0_1,
	}
)

func getTranslator(variant string, version semver.Version) (translator, error) {
	t, ok := registry[fmt.Sprintf("%s+%s", variant, version.String())]
	if !ok {
		return nil, fmt.Errorf("No translator exists for variant %s with version %s", variant, version.String())
	}
	return t, nil
}

type TranslateOptions struct {
	Pretty bool
	Strict bool
}

type Common struct {
	Version string `yaml:"version"`
	Variant string `yaml:"variant"`
}

type translator func([]byte, TranslateOptions) ([]byte, error)

// Translate wraps all of the actual translate functions in a switch that determines the correct one to call
func Translate(input []byte, options TranslateOptions) ([]byte, error) {
	// first determine version. This will ignore most fields, so don't use strict
	ver := Common{}
	if err := yaml.Unmarshal(input, &ver); err != nil {
		return nil, fmt.Errorf("Error unmarshaling yaml: %v", err)
	}

	if ver.Variant == "" {
		return nil, ErrNoVariant
	}

	tmp, err := semver.NewVersion(ver.Version)
	if err != nil {
		return nil, ErrInvalidVersion
	}
	version := *tmp

	translator, err := getTranslator(ver.Variant, version)
	if err != nil {
		return nil, err
	}

	return translator(input, options)
}

// Misc helpers
func unmarshal(data []byte, to interface{}, strict bool) error {
	if strict {
		return yaml.UnmarshalStrict(data, to)
	}
	return yaml.Unmarshal(data, to)
}

func marshal(from interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(from, "", "  ")
	}
	return json.Marshal(from)
}
