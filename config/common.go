package config

import (
	"errors"
	"fmt"

	yaml "github.com/ajeddeloh/yaml"
	"github.com/coreos/go-semver/semver"
)

var (
	ErrNoVariant      = errors.New("Error parsing variant. Variant must be specified")
	ErrInvalidVersion = errors.New("Error parsing version. Version must be a valid semver")

	registry = map[string]translator{}
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
	}

	return translator(input, options)
}
