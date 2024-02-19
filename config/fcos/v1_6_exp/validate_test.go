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
	"fmt"
	"testing"

	baseutil "github.com/coreos/butane/base/util"
	base "github.com/coreos/butane/base/v0_6_exp"
	"github.com/coreos/butane/config/common"

	"github.com/coreos/ignition/v2/config/shared/errors"
	"github.com/coreos/ignition/v2/config/util"
	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
	"github.com/stretchr/testify/assert"
)

// TestReportCorrelation tests that errors are correctly correlated to their source lines
func TestReportCorrelation(t *testing.T) {
	tests := []struct {
		in      string
		message string
		line    int64
	}{
		// Butane unused key check
		{
			`storage:
                           files:
                           - path: /z
                             q: z`,
			"Unused key q",
			4,
		},
		// Butane YAML validation error
		{
			`storage:
                           files:
                           - path: /z
                             contents:
                               source: https://example.com
                               inline: z`,
			common.ErrTooManyResourceSources.Error(),
			5,
		},
		// Butane YAML validation warning
		{
			`storage:
                           files:
                           - path: /z
                             mode: 444`,
			common.ErrDecimalMode.Error(),
			4,
		},
		// Butane translation error
		{
			`storage:
                           files:
                           - path: /z
                             contents:
                               local: z`,
			common.ErrNoFilesDir.Error(),
			5,
		},
		// Ignition validation error, leaf node
		{
			`storage:
                           files:
                           - path: z`,
			errors.ErrPathRelative.Error(),
			3,
		},
		// Ignition validation error, partition
		{
			`storage:
                           disks:
                           - device: /dev/z
                             wipe_table: true
                             partitions:
                               - start_mib: 5`,
			errors.ErrNeedLabelOrNumber.Error(),
			6,
		},
		// Ignition validation error, partition list
		{
			`storage:
                           disks:
                           - device: /dev/z
                             wipe_table: true
                             partitions:
                               - number: 1
                                 should_exist: false
                               - label: z`,
			errors.ErrZeroesWithShouldNotExist.Error(),
			6,
		},
		// Ignition duplicate key check, paths
		{
			`storage:
                           files:
                           - path: /z
                           - path: /z`,
			errors.ErrDuplicate.Error(),
			4,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("validate %d", i), func(t *testing.T) {
			_, r, _ := ToIgn3_5Bytes([]byte(test.in), common.TranslateBytesOptions{})
			assert.Len(t, r.Entries, 1, "unexpected report length")
			assert.Equal(t, test.message, r.Entries[0].Message, "bad error")
			assert.NotNil(t, r.Entries[0].Marker.StartP, "marker start is nil")
			assert.Equal(t, test.line, r.Entries[0].Marker.StartP.Line, "incorrect error line")
		})
	}
}

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
		t.Run(fmt.Sprintf("validate %d", i), func(t *testing.T) {
			actual := test.in.Validate(path.New("yaml"))
			baseutil.VerifyReport(t, test.in, actual)
			expected := report.Report{}
			expected.AddOnError(test.errPath, test.out)
			assert.Equal(t, expected, actual, "bad validation report")
		})
	}
}

func TestValidateGrubUser(t *testing.T) {
	tests := []struct {
		in      GrubUser
		out     error
		errPath path.ContextPath
	}{
		// valid user
		{
			in: GrubUser{
				Name:         "name",
				PasswordHash: util.StrToPtr("pkcs5-pass"),
			},
			out:     nil,
			errPath: path.New("yaml"),
		},
		// username is not specified
		{
			in: GrubUser{
				Name:         "",
				PasswordHash: util.StrToPtr("pkcs5-pass"),
			},
			out:     common.ErrGrubUserNameNotSpecified,
			errPath: path.New("yaml", "name"),
		},
		// password is not specified
		{
			in: GrubUser{
				Name: "name",
			},
			out:     common.ErrGrubPasswordNotSpecified,
			errPath: path.New("yaml", "password_hash"),
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("validate %d", i), func(t *testing.T) {
			actual := test.in.Validate(path.New("yaml"))
			baseutil.VerifyReport(t, test.in, actual)
			expected := report.Report{}
			expected.AddOnError(test.errPath, test.out)
			assert.Equal(t, expected, actual, "bad report")
		})
	}
}

func TestValidateMountPoints(t *testing.T) {
	tests := []struct {
		in      Config
		out     error
		errPath path.ContextPath
	}{
		// valid config (has prefix "/etc" or "/var")
		{
			in: Config{
				Config: base.Config{
					Storage: base.Storage{
						Filesystems: []base.Filesystem{
							{
								Path:          util.StrToPtr("/etc/foo"),
								WithMountUnit: util.BoolToPtr(true),
							},
							{
								Path:          util.StrToPtr("/var"),
								WithMountUnit: util.BoolToPtr(true),
							},
							{
								Path:          util.StrToPtr("/invalid/path"),
								WithMountUnit: util.BoolToPtr(false),
							},
							{
								WithMountUnit: util.BoolToPtr(true),
							},
							{
								Path:          nil,
								WithMountUnit: util.BoolToPtr(true),
							},
						},
					},
				},
			},
		},
		// invalid config (path name is '/')
		{
			in: Config{
				Config: base.Config{
					Storage: base.Storage{
						Filesystems: []base.Filesystem{
							{
								Path:          util.StrToPtr("/"),
								WithMountUnit: util.BoolToPtr(true),
							},
						},
					},
				},
			},
			out:     common.ErrMountPointForbidden,
			errPath: path.New("yaml", "storage", "filesystems", 0, "path"),
		},
		// invalid config (path is /boot)
		{
			in: Config{
				Config: base.Config{
					Storage: base.Storage{
						Filesystems: []base.Filesystem{
							{
								Path:          util.StrToPtr("/boot"),
								WithMountUnit: util.BoolToPtr(true),
							},
						},
					},
				},
			},

			out:     common.ErrMountPointForbidden,
			errPath: path.New("yaml", "storage", "filesystems", 0, "path"),
		},
		// invalid config (path is invalid, does not contain /etc or /var)
		{
			in: Config{
				Config: base.Config{
					Storage: base.Storage{
						Filesystems: []base.Filesystem{
							{
								Path:          util.StrToPtr("/thisIsABugTest"),
								WithMountUnit: util.BoolToPtr(true),
							},
						},
					},
				},
			},

			out:     common.ErrMountPointForbidden,
			errPath: path.New("yaml", "storage", "filesystems", 0, "path"),
		},
		// invalid config (path is /varnish)
		{
			in: Config{
				Config: base.Config{
					Storage: base.Storage{
						Filesystems: []base.Filesystem{
							{
								Path:          util.StrToPtr("/varnish"),
								WithMountUnit: util.BoolToPtr(true),
							},
						},
					},
				},
			},

			out:     common.ErrMountPointForbidden,
			errPath: path.New("yaml", "storage", "filesystems", 0, "path"),
		},
		// invalid config (path is /foo/var)
		{
			in: Config{
				Config: base.Config{
					Storage: base.Storage{
						Filesystems: []base.Filesystem{
							{
								Path:          util.StrToPtr("/foo/var"),
								WithMountUnit: util.BoolToPtr(true),
							},
						},
					},
				},
			},

			out:     common.ErrMountPointForbidden,
			errPath: path.New("yaml", "storage", "filesystems", 0, "path"),
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("validate %d", i), func(t *testing.T) {
			actual := test.in.Validate(path.New("yaml"))
			baseutil.VerifyReport(t, test.in, actual)
			expected := report.Report{}
			expected.AddOnError(test.errPath, test.out)
			assert.Equal(t, expected, actual, "invalid report")
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		in      Config
		out     error
		errPath path.ContextPath
	}{
		// valid config (wipe_table is true)
		{
			in: Config{
				Config: base.Config{
					Storage: base.Storage{
						Disks: []base.Disk{
							{
								Device:    "/dev/vda",
								WipeTable: util.BoolToPtr(true),
								Partitions: []base.Partition{
									{
										Label: util.StrToPtr("foo"),
									},
								},
							},
						},
					},
				},
			},
		},
		// valid config (disk is /dev/disk/by-id/coreos-boot-disk)
		{
			in: Config{
				Config: base.Config{
					Storage: base.Storage{
						Disks: []base.Disk{
							{
								Device:    rootDevice,
								WipeTable: util.BoolToPtr(false),
								Partitions: []base.Partition{
									{
										Label: util.StrToPtr("bar"),
									},
								},
							},
						},
					},
				},
			},
		},
		// invalid config (wipe_table is nil)
		{
			in: Config{
				Config: base.Config{
					Storage: base.Storage{
						Disks: []base.Disk{
							{
								Device: "/dev/vda",
								Partitions: []base.Partition{
									{
										Label: util.StrToPtr("foo"),
									},
								},
							},
						},
					},
				},
			},
			out:     common.ErrReuseByLabel,
			errPath: path.New("yaml", "storage", "disks", 0, "partitions", 0, "number"),
		},
		// invalid config (wipe_table is false with a partition numbered 0)
		{
			in: Config{
				Config: base.Config{
					Storage: base.Storage{
						Disks: []base.Disk{
							{
								Device:    "/dev/vda",
								WipeTable: util.BoolToPtr(false),
								Partitions: []base.Partition{
									{
										Label: util.StrToPtr("foo"),
									},
									{
										Label:  util.StrToPtr("bar"),
										Number: 2,
									},
								},
							},
						},
					},
				},
			},
			out:     common.ErrReuseByLabel,
			errPath: path.New("yaml", "storage", "disks", 0, "partitions", 0, "number"),
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("validate %d", i), func(t *testing.T) {
			actual := test.in.Validate(path.New("yaml"))
			baseutil.VerifyReport(t, test.in, actual)
			expected := report.Report{}
			expected.AddOnWarn(test.errPath, test.out)
			assert.Equal(t, expected, actual, "invalid report")
		})
	}
}
