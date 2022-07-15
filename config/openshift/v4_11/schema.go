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

package v4_11

import (
	fcos "github.com/coreos/butane/config/fcos/v1_3"
	"github.com/coreos/butane/config/openshift/v4_8"
)

const ROLE_LABEL_KEY = "machineconfiguration.openshift.io/role"

type Config struct {
	fcos.Config `yaml:",inline"`
	Metadata    Metadata  `yaml:"metadata"`
	OpenShift   OpenShift `yaml:"openshift"`
}

type Metadata v4_8.Metadata

type OpenShift v4_8.OpenShift
