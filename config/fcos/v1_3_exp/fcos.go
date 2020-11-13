// Copyright 2020 Red Hat, Inc
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

package v1_3_exp

import (
	base_0_4 "github.com/coreos/fcct/base/v0_4_exp"
	"github.com/coreos/fcct/config/common"
	"github.com/coreos/fcct/config/util"
	fcos_0_1 "github.com/coreos/fcct/distro/fcos/v0_1"
	"github.com/coreos/fcct/translate"

	"github.com/coreos/ignition/v2/config/v3_3_experimental/types"
	"github.com/coreos/vcontext/report"
)

type Config struct {
	common.Common   `yaml:",inline"`
	base_0_4.Config `yaml:",inline"`
	fcos_0_1.Fcos   `yaml:",inline"`
}

func (c Config) Translate(options common.TranslateOptions) (types.Config, translate.TranslationSet, report.Report) {
	cfg, baseTranslations, baseReport := c.Config.ToIgn3_3(options.BaseOptions)
	if baseReport.IsFatal() {
		return types.Config{}, translate.TranslationSet{}, baseReport
	}

	finalcfg, distroTranslations, distroReport := c.Fcos.ToIgn3_3(cfg, options.BaseOptions)
	baseReport.Merge(distroReport)
	if baseReport.IsFatal() {
		return types.Config{}, translate.TranslationSet{}, baseReport
	}

	baseTranslations.Merge(distroTranslations)

	return finalcfg, baseTranslations, baseReport
}

// TranslateBytes translates from a v1.3 fcc to a v3.3.0 Ignition config. It returns a report of any errors or
// warnings in the source and resultant config. If the report has fatal errors or it encounters other problems
// translating, an error is returned.
func TranslateBytes(input []byte, options common.TranslateOptions) ([]byte, report.Report, error) {
	return util.TranslateBytes(input, &Config{}, options)
}
