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

package v1_6_exp

import (
	base "github.com/coreos/butane/base/v0_6_exp"
)

type Config struct {
	base.Config `yaml:",inline"`
	BootDevice  BootDevice `yaml:"boot_device"`
	Grub        Grub       `yaml:"grub"`
	Selinux     Selinux    `yaml:"selinux"`
}

type BootDevice struct {
	Layout *string          `yaml:"layout"`
	Luks   BootDeviceLuks   `yaml:"luks"`
	Mirror BootDeviceMirror `yaml:"mirror"`
}

type BootDeviceLuks struct {
	Discard   *bool       `yaml:"discard"`
	Device    *string     `yaml:"device"`
	Tang      []base.Tang `yaml:"tang"`
	Threshold *int        `yaml:"threshold"`
	Tpm2      *bool       `yaml:"tpm2"`
}

type BootDeviceMirror struct {
	Devices []string `yaml:"devices"`
}

type Grub struct {
	Users []GrubUser `yaml:"users"`
}

type GrubUser struct {
	Name         string  `yaml:"name"`
	PasswordHash *string `yaml:"password_hash"`
}

type Selinux struct {
	Modules []Module `yaml:"modules"`
}

type Module struct {
	Name     string   `yaml:"name"`
	Contents Resource `yaml:"contents"`
}

type Resource struct {
	Compression  *string      `yaml:"compression"`
	HTTPHeaders  HTTPHeaders  `yaml:"http_headers"`
	Source       *string      `yaml:"source"`
	Inline       *string      `yaml:"inline"` // Added, not in ignition spec
	Local        *string      `yaml:"local"`  // Added, not in ignition spec
	Verification Verification `yaml:"verification"`
}

type HTTPHeader struct {
	Name  string  `yaml:"name"`
	Value *string `yaml:"value"`
}

type HTTPHeaders []HTTPHeader

type Verification struct {
	Hash *string `yaml:"hash"`
}
