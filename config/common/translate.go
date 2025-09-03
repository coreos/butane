package common

import (
	"fmt"
	"github.com/coreos/go-semver/semver"
	"github.com/coreos/vcontext/report"
	"gopkg.in/yaml.v3"
)

var (
	registry = map[string]translator{}
)

// Fields that must be included in the root struct of every spec version.
type commonFields struct {
	Version string `yaml:"version"`
	Variant string `yaml:"variant"`
}

// RegisterTranslator registers a translator for the specified variant and
// version to be available for use by TranslateBytes.  This is only needed
// by users implementing their own translators outside the Butane package.
func RegisterTranslator(variant, version string, trans translator) {
	key := fmt.Sprintf("%s+%s", variant, version)
	if _, ok := registry[key]; ok {
		panic("tried to reregister existing translator")
	}
	registry[key] = trans
}

func getTranslator(variant string, version semver.Version) (translator, error) {
	t, ok := registry[fmt.Sprintf("%s+%s", variant, version.String())]
	if !ok {
		return nil, ErrUnknownVersion{
			Variant: variant,
			Version: version,
		}
	}
	return t, nil
}

// translators take a raw config and translate it to a raw Ignition config. The report returned should include any
// errors, warnings, etc. and may or may not be fatal. If report is fatal, or other errors are encountered while translating
// translators should return an error.
type translator func([]byte, TranslateBytesOptions) ([]byte, report.Report, error)

// TranslateBytes wraps all of the individual TranslateBytes functions in a switch that determines the correct one to call.
// TranslateBytes returns an error if the report had fatal errors or if other errors occured during translation.
func TranslateBytes(input []byte, options TranslateBytesOptions) ([]byte, report.Report, error) {
	// first determine version; this will ignore most fields
	ver := commonFields{}
	if err := yaml.Unmarshal(input, &ver); err != nil {
		return nil, report.Report{}, ErrUnmarshal{
			Detail: err.Error(),
		}
	}

	if ver.Variant == "" {
		return nil, report.Report{}, ErrNoVariant
	}

	tmp, err := semver.NewVersion(ver.Version)
	if err != nil {
		return nil, report.Report{}, ErrInvalidVersion
	}
	version := *tmp

	translator, err := getTranslator(ver.Variant, version)
	if err != nil {
		return nil, report.Report{}, err
	}

	return translator(input, options)
}
