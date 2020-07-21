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

package v1_0

import (
	"reflect"

	base_0_1 "github.com/coreos/fcct/base/v0_1"
	"github.com/coreos/fcct/config/common"
	fcos_0_1 "github.com/coreos/fcct/distro/fcos/v0_1"
	"github.com/coreos/fcct/translate"

	"github.com/coreos/ignition/v2/config/v3_0"
	"github.com/coreos/ignition/v2/config/v3_0/types"
	ignvalidate "github.com/coreos/ignition/v2/config/validate"
	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
	"github.com/coreos/vcontext/validate"
)

type Config struct {
	common.Common   `yaml:",inline"`
	base_0_1.Config `yaml:",inline"`
	fcos_0_1.Fcos   `yaml:",inline"`
}

func (c Config) Translate(options common.TranslateOptions) (types.Config, translate.TranslationSet, report.Report) {
	cfg, baseTranslations, baseReport := c.Config.ToIgn3_0(options.BaseOptions)
	if baseReport.IsFatal() {
		return types.Config{}, translate.TranslationSet{}, baseReport
	}

	finalcfg, distroTranslations, distroReport := c.Fcos.ToIgn3_0(cfg, options.BaseOptions)
	baseReport.Merge(distroReport)
	if baseReport.IsFatal() {
		return types.Config{}, translate.TranslationSet{}, baseReport
	}

	baseTranslations.Merge(distroTranslations)

	return finalcfg, baseTranslations, baseReport
}

// TranslateBytes translates from a v1.0 fcc to a v3.0.0 Ignition config. It returns a report of any errors or
// warnings in the source and resultant config. If the report has fatal errors or it encounters other problems
// translating, an error is returned.
func TranslateBytes(input []byte, options common.TranslateOptions) ([]byte, report.Report, error) {
	cfg := Config{}

	contextTree, err := common.Unmarshal(input, &cfg, options.Strict)
	if err != nil {
		return nil, report.Report{}, err
	}

	r := validate.Validate(cfg, "yaml")
	unusedKeyCheck := func(v reflect.Value, c path.ContextPath) report.Report {
		return ignvalidate.ValidateUnusedKeys(v, c, contextTree)
	}
	r.Merge(validate.ValidateCustom(cfg, "yaml", unusedKeyCheck))
	r.Correlate(contextTree)
	if r.IsFatal() {
		return nil, r, common.ErrInvalidSourceConfig
	}

	final, translations, translateReport := cfg.Translate(options)
	translateReport.Correlate(contextTree)
	r.Merge(translateReport)
	if r.IsFatal() {
		return nil, r, common.ErrInvalidSourceConfig
	}

	// Check for invalid duplicated keys.
	dupsReport := validate.ValidateCustom(final, "json", ignvalidate.ValidateDups)
	common.TranslateReportPaths(&dupsReport, translations)
	dupsReport.Correlate(contextTree)
	r.Merge(dupsReport)

	// Validate JSON semantics.
	jsonReport := validate.Validate(final, "json")
	common.TranslateReportPaths(&jsonReport, translations)
	jsonReport.Correlate(contextTree)
	r.Merge(jsonReport)

	if r.IsFatal() {
		return nil, r, common.ErrInvalidGeneratedConfig
	}

	outbytes, err := common.Marshal(final, options.Pretty)
	return outbytes, r, err
}

func MergeBytes(fragments [][]byte, options common.TranslateOptions) ([]byte, error) {
	// merged Ignition Config
	final := types.Config{}

	for _, fragment := range fragments {
		cfg := Config{}
		_, err := common.Unmarshal(fragment, &cfg, options.Strict)
		if err != nil {
			return nil, err
		}

		fragmentIgn, _, _ := cfg.Translate(options)
		v3_0.Merge(final, fragmentIgn)
	}
	outbytes, err := common.Marshal(final, options.Pretty)
	return outbytes, err
}
