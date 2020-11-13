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

package v1_1

import (
	base_0_2 "github.com/coreos/fcct/base/v0_2"
	"github.com/coreos/fcct/config/common"
	"github.com/coreos/fcct/config/util"

	"github.com/coreos/vcontext/report"
)

type Config struct {
	base_0_2.Config `yaml:",inline"`
}

// TranslateBytes translates from a v1.1 fcc to a v3.1.0 Ignition config. It returns a report of any errors or
// warnings in the source and resultant config. If the report has fatal errors or it encounters other problems
// translating, an error is returned.
func TranslateBytes(input []byte, options common.TranslateBytesOptions) ([]byte, report.Report, error) {
	return util.TranslateBytes(input, &Config{}, "ToIgn3_1", options)
}
