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

package v1_3

import (
	"testing"

	base "github.com/coreos/fcct/base/v0_3"
	"github.com/coreos/fcct/config/common"

	"github.com/coreos/ignition/v2/config/util"
	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
	"github.com/stretchr/testify/assert"
)

// TestValidateBootDevice tests boot device validation
func TestValidateBootDevice(t *testing.T) {
	tests := []struct {
		in      BootDevice
		out     error
		errPath path.ContextPath
	}{
		// empty config
		{
			BootDevice{},
			nil,
			path.New("yaml"),
		},
		// complete config
		{
			BootDevice{
				Layout: util.StrToPtr("x86_64"),
				Luks: BootDeviceLuks{
					Tang: []base.Tang{{
						URL:        "https://example.com/",
						Thumbprint: util.StrToPtr("x"),
					}},
					Threshold: util.IntToPtr(2),
					Tpm2:      util.BoolToPtr(true),
				},
				Mirror: BootDeviceMirror{
					Devices: []string{"/dev/vda", "/dev/vdb"},
				},
			},
			nil,
			path.New("yaml"),
		},
		// invalid layout
		{
			BootDevice{
				Layout: util.StrToPtr("sparc"),
			},
			common.ErrUnknownBootDeviceLayout,
			path.New("yaml", "layout"),
		},
		// only one mirror device
		{
			BootDevice{
				Mirror: BootDeviceMirror{
					Devices: []string{"/dev/vda"},
				},
			},
			common.ErrTooFewMirrorDevices,
			path.New("yaml", "mirror", "devices"),
		},
	}

	for i, test := range tests {
		actual := test.in.Validate(path.New("yaml"))
		expected := report.Report{}
		expected.AddOnError(test.errPath, test.out)
		assert.Equal(t, expected, actual, "#%d: bad validation report", i)
	}
}
