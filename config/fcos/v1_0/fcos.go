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
	base_0_1 "github.com/coreos/fcct/base/v0_1"
	"github.com/coreos/fcct/config/common"
	"github.com/coreos/fcct/config/util"

	"github.com/coreos/ignition/v2/config/v3_0/types"
	"github.com/coreos/vcontext/report"
)

type Config struct {
	base_0_1.Config `yaml:",inline"`
}

// ToIgn3_0 translates the config to an Ignition config.  It returns a
// report of any errors or warnings in the source and resultant config.  If
// the report has fatal errors or it encounters other problems translating,
// an error is returned.
func (c Config) ToIgn3_0(options common.TranslateOptions) (types.Config, report.Report, error) {
	cfg, r, err := util.Translate(c, "ToIgn3_0Unvalidated", options)
	return cfg.(types.Config), r, err
}

// TranslateBytes translates from a v1.0 fcc to a v3.0.0 Ignition config. It returns a report of any errors or
// warnings in the source and resultant config. If the report has fatal errors or it encounters other problems
// translating, an error is returned.
func TranslateBytes(input []byte, options common.TranslateBytesOptions) ([]byte, report.Report, error) {
	return util.TranslateBytes(input, &Config{}, "ToIgn3_0", options)
}
