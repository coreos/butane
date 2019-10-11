// Copyright 2019 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.)

package config

import (
	"errors"
	"fmt"

	"github.com/coreos/fcct/config/common"
	"github.com/coreos/fcct/config/v1_0"
	"github.com/coreos/fcct/config/v1_1_exp"

	"github.com/coreos/go-semver/semver"
	"github.com/coreos/vcontext/report"
	"gopkg.in/yaml.v3"
)

var (
	ErrNoVariant      = errors.New("Error parsing variant. Variant must be specified")
	ErrInvalidVersion = errors.New("Error parsing version. Version must be a valid semver")

	registry = map[string]translator{
		"fcos+1.0.0":              v1_0.TranslateBytes,
		"fcos+1.1.0-experimental": v1_1_exp.TranslateBytes,
	}
)

func getTranslator(variant string, version semver.Version) (translator, error) {
	t, ok := registry[fmt.Sprintf("%s+%s", variant, version.String())]
	if !ok {
		return nil, fmt.Errorf("No translator exists for variant %s with version %s", variant, version.String())
	}
	return t, nil
}

// translators take a raw config and translate it to a raw Ignition config. The report returned should include any
// errors, warnings, etc and may or may not be fatal. If report is fatal, or other errors are encountered while translating
// translators should return an error.
type translator func([]byte, common.TranslateOptions) ([]byte, report.Report, error)

// Translate wraps all of the actual translate functions in a switch that determines the correct one to call.
// Translate returns an error if the report had fatal errors or if other errors occured during translation.
func Translate(input []byte, options common.TranslateOptions) ([]byte, report.Report, error) {
	// first determine version. This will ignore most fields, so don't use strict
	ver := common.Common{}
	if err := yaml.Unmarshal(input, &ver); err != nil {
		return nil, report.Report{}, fmt.Errorf("Error unmarshaling yaml: %v", err)
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
