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
	"github.com/coreos/butane/config/common"
	fcos1_0 "github.com/coreos/butane/config/fcos/v1_0"
	fcos1_1 "github.com/coreos/butane/config/fcos/v1_1"
	fcos1_2 "github.com/coreos/butane/config/fcos/v1_2"
	fcos1_3 "github.com/coreos/butane/config/fcos/v1_3"
	fcos1_4 "github.com/coreos/butane/config/fcos/v1_4"
	fcos1_5 "github.com/coreos/butane/config/fcos/v1_5"
	fcos1_6 "github.com/coreos/butane/config/fcos/v1_6"
	fcos1_7_exp "github.com/coreos/butane/config/fcos/v1_7_exp"
	fiot1_0 "github.com/coreos/butane/config/fiot/v1_0"
	fiot1_1_exp "github.com/coreos/butane/config/fiot/v1_1_exp"
	flatcar1_0 "github.com/coreos/butane/config/flatcar/v1_0"
	flatcar1_1 "github.com/coreos/butane/config/flatcar/v1_1"
	flatcar1_2_exp "github.com/coreos/butane/config/flatcar/v1_2_exp"
	openshift4_10 "github.com/coreos/butane/config/openshift/v4_10"
	openshift4_11 "github.com/coreos/butane/config/openshift/v4_11"
	openshift4_12 "github.com/coreos/butane/config/openshift/v4_12"
	openshift4_13 "github.com/coreos/butane/config/openshift/v4_13"
	openshift4_14 "github.com/coreos/butane/config/openshift/v4_14"
	openshift4_15 "github.com/coreos/butane/config/openshift/v4_15"
	openshift4_16 "github.com/coreos/butane/config/openshift/v4_16"
	openshift4_17 "github.com/coreos/butane/config/openshift/v4_17"
	openshift4_18 "github.com/coreos/butane/config/openshift/v4_18"
	openshift4_19 "github.com/coreos/butane/config/openshift/v4_19"
	openshift4_20_exp "github.com/coreos/butane/config/openshift/v4_20_exp"
	openshift4_8 "github.com/coreos/butane/config/openshift/v4_8"
	openshift4_9 "github.com/coreos/butane/config/openshift/v4_9"
	r4e1_0 "github.com/coreos/butane/config/r4e/v1_0"
	r4e1_1 "github.com/coreos/butane/config/r4e/v1_1"
	r4e1_2_exp "github.com/coreos/butane/config/r4e/v1_2_exp"

	"github.com/coreos/vcontext/report"
)

func init() {
	common.RegisterTranslator("fcos", "1.0.0", fcos1_0.ToIgn3_0Bytes)
	common.RegisterTranslator("fcos", "1.1.0", fcos1_1.ToIgn3_1Bytes)
	common.RegisterTranslator("fcos", "1.2.0", fcos1_2.ToIgn3_2Bytes)
	common.RegisterTranslator("fcos", "1.3.0", fcos1_3.ToIgn3_2Bytes)
	common.RegisterTranslator("fcos", "1.4.0", fcos1_4.ToIgn3_3Bytes)
	common.RegisterTranslator("fcos", "1.5.0", fcos1_5.ToIgn3_4Bytes)
	common.RegisterTranslator("fcos", "1.6.0", fcos1_6.ToIgn3_5Bytes)
	common.RegisterTranslator("fcos", "1.7.0-experimental", fcos1_7_exp.ToIgn3_6Bytes)
	common.RegisterTranslator("flatcar", "1.0.0", flatcar1_0.ToIgn3_3Bytes)
	common.RegisterTranslator("flatcar", "1.1.0", flatcar1_1.ToIgn3_4Bytes)
	common.RegisterTranslator("flatcar", "1.2.0-experimental", flatcar1_2_exp.ToIgn3_6Bytes)
	common.RegisterTranslator("openshift", "4.8.0", openshift4_8.ToConfigBytes)
	common.RegisterTranslator("openshift", "4.9.0", openshift4_9.ToConfigBytes)
	common.RegisterTranslator("openshift", "4.10.0", openshift4_10.ToConfigBytes)
	common.RegisterTranslator("openshift", "4.11.0", openshift4_11.ToConfigBytes)
	common.RegisterTranslator("openshift", "4.12.0", openshift4_12.ToConfigBytes)
	common.RegisterTranslator("openshift", "4.13.0", openshift4_13.ToConfigBytes)
	common.RegisterTranslator("openshift", "4.14.0", openshift4_14.ToConfigBytes)
	common.RegisterTranslator("openshift", "4.15.0", openshift4_15.ToConfigBytes)
	common.RegisterTranslator("openshift", "4.16.0", openshift4_16.ToConfigBytes)
	common.RegisterTranslator("openshift", "4.17.0", openshift4_17.ToConfigBytes)
	common.RegisterTranslator("openshift", "4.18.0", openshift4_18.ToConfigBytes)
	common.RegisterTranslator("openshift", "4.19.0", openshift4_19.ToConfigBytes)
	common.RegisterTranslator("openshift", "4.20.0-experimental", openshift4_20_exp.ToConfigBytes)
	common.RegisterTranslator("r4e", "1.0.0", r4e1_0.ToIgn3_3Bytes)
	common.RegisterTranslator("r4e", "1.1.0", r4e1_1.ToIgn3_4Bytes)
	common.RegisterTranslator("r4e", "1.2.0-experimental", r4e1_2_exp.ToIgn3_6Bytes)
	common.RegisterTranslator("fiot", "1.0.0", fiot1_0.ToIgn3_4Bytes)
	common.RegisterTranslator("fiot", "1.1.0-experimental", fiot1_1_exp.ToIgn3_6Bytes)
	common.RegisterTranslator("rhcos", "0.1.0", unsupportedRhcosVariant)
}

func unsupportedRhcosVariant(input []byte, options common.TranslateBytesOptions) ([]byte, report.Report, error) {
	return nil, report.Report{}, common.ErrRhcosVariantUnsupported
}
