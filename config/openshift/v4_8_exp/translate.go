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

package v4_8_exp

import (
	"github.com/coreos/fcct/config/common"
	"github.com/coreos/fcct/config/openshift/v4_8_exp/result"
	cutil "github.com/coreos/fcct/config/util"
	"github.com/coreos/fcct/translate"

	"github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
)

// ToMachineConfig4_8Unvalidated translates the config to a MachineConfig.  It also
// returns the set of translations it did so paths in the resultant config
// can be tracked back to their source in the source config.  No config
// validation is performed on input or output.
func (c Config) ToMachineConfig4_8Unvalidated(options common.TranslateOptions) (result.MachineConfig, translate.TranslationSet, report.Report) {
	cfg, ts, r := c.Config.ToIgn3_2Unvalidated(options)
	if r.IsFatal() {
		return result.MachineConfig{}, ts, r
	}

	// wrap
	ts = ts.PrefixPaths(path.New("yaml"), path.New("json", "spec", "config"))
	mc := result.MachineConfig{
		ApiVersion: "machineconfiguration.openshift.io/v1",
		Kind:       "MachineConfig",
		Metadata: result.Metadata{
			Name:   c.Metadata.Name,
			Labels: make(map[string]string),
		},
		Spec: result.Spec{
			Config: cfg,
		},
	}
	ts.AddTranslation(path.New("yaml", "version"), path.New("json", "apiVersion"))
	ts.AddTranslation(path.New("yaml", "version"), path.New("json", "kind"))
	ts.AddTranslation(path.New("yaml", "metadata"), path.New("json", "metadata"))
	ts.AddTranslation(path.New("yaml", "metadata", "name"), path.New("json", "metadata", "name"))
	ts.AddTranslation(path.New("yaml", "version"), path.New("json", "spec"))
	ts.AddTranslation(path.New("yaml"), path.New("json", "spec", "config"))
	for k, v := range c.Metadata.Labels {
		mc.Metadata.Labels[k] = v
		ts.AddTranslation(path.New("yaml", "metadata", "labels", k), path.New("json", "metadata", "labels", k))
	}
	if len(mc.Metadata.Labels) > 0 {
		ts.AddTranslation(path.New("yaml", "metadata", "labels"), path.New("json", "metadata", "labels"))
	}

	return mc, ts, r
}

// ToMachineConfig4_8 translates the config to a MachineConfig.  It returns a
// report of any errors or warnings in the source and resultant config.  If
// the report has fatal errors or it encounters other problems translating,
// an error is returned.
func (c Config) ToMachineConfig4_8(options common.TranslateOptions) (result.MachineConfig, report.Report, error) {
	cfg, r, err := cutil.Translate(c, "ToMachineConfig4_8Unvalidated", options)
	return cfg.(result.MachineConfig), r, err
}

// ToIgn3_2Unvalidated translates the config to an Ignition config.  It also
// returns the set of translations it did so paths in the resultant config
// can be tracked back to their source in the source config.  No config
// validation is performed on input or output.
func (c Config) ToIgn3_2Unvalidated(options common.TranslateOptions) (types.Config, translate.TranslationSet, report.Report) {
	mc, ts, r := c.ToMachineConfig4_8Unvalidated(options)
	cfg := mc.Spec.Config
	ts = ts.Descend(path.New("json", "spec", "config"))
	return cfg, ts, r
}

// ToIgn3_2 translates the config to an Ignition config.  It returns a
// report of any errors or warnings in the source and resultant config.  If
// the report has fatal errors or it encounters other problems translating,
// an error is returned.
func (c Config) ToIgn3_2(options common.TranslateOptions) (types.Config, report.Report, error) {
	cfg, r, err := cutil.Translate(c, "ToIgn3_2Unvalidated", options)
	return cfg.(types.Config), r, err
}

// ToConfigBytes translates from a v4.8 occ to a v4.8 MachineConfig or a v3.2.0 Ignition config. It returns a report of any errors or
// warnings in the source and resultant config. If the report has fatal errors or it encounters other problems
// translating, an error is returned.
func ToConfigBytes(input []byte, options common.TranslateBytesOptions) ([]byte, report.Report, error) {
	if options.Raw {
		return cutil.TranslateBytes(input, &Config{}, "ToIgn3_2", options)
	} else {
		return cutil.TranslateBytesYAML(input, &Config{}, "ToMachineConfig4_8", options)
	}
}
