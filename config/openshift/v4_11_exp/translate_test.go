// Copyright 2021 Red Hat, Inc
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

package v4_11_exp

import (
	"fmt"
	"testing"

	baseutil "github.com/coreos/butane/base/util"
	base "github.com/coreos/butane/base/v0_5_exp"
	"github.com/coreos/butane/config/common"
	fcos "github.com/coreos/butane/config/fcos/v1_5_exp"
	"github.com/coreos/butane/config/openshift/v4_11_exp/result"
	"github.com/coreos/butane/translate"

	"github.com/coreos/ignition/v2/config/util"
	"github.com/coreos/ignition/v2/config/v3_4_experimental/types"
	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
	"github.com/stretchr/testify/assert"
)

// TestElidedFieldWarning tests that we warn when transpiling fields to an
// Ignition config that can't be represented in an Ignition config.
func TestElidedFieldWarning(t *testing.T) {
	in := Config{
		Metadata: Metadata{
			Name: "z",
		},
		OpenShift: OpenShift{
			KernelArguments: []string{"a", "b"},
			FIPS:            util.BoolToPtr(true),
			KernelType:      util.StrToPtr("realtime"),
		},
	}

	var expected report.Report
	expected.AddOnWarn(path.New("yaml", "openshift", "kernel_arguments"), common.ErrFieldElided)
	expected.AddOnWarn(path.New("yaml", "openshift", "fips"), common.ErrFieldElided)
	expected.AddOnWarn(path.New("yaml", "openshift", "kernel_type"), common.ErrFieldElided)

	_, _, r := in.ToIgn3_4Unvalidated(common.TranslateOptions{})
	assert.Equal(t, expected, r, "report mismatch")
}

func TestTranslateConfig(t *testing.T) {
	tests := []struct {
		in         Config
		out        result.MachineConfig
		exceptions []translate.Translation
	}{
		// empty-ish config
		{
			Config{
				Metadata: Metadata{
					Name: "z",
					Labels: map[string]string{
						ROLE_LABEL_KEY: "z",
					},
				},
			},
			result.MachineConfig{
				ApiVersion: result.MC_API_VERSION,
				Kind:       result.MC_KIND,
				Metadata: result.Metadata{
					Name: "z",
					Labels: map[string]string{
						ROLE_LABEL_KEY: "z",
					},
				},
				Spec: result.Spec{
					Config: types.Config{
						Ignition: types.Ignition{
							Version: "3.4.0-experimental",
						},
					},
				},
			},
			[]translate.Translation{
				{path.New("yaml", "version"), path.New("json", "apiVersion")},
				{path.New("yaml", "version"), path.New("json", "kind")},
				{path.New("yaml", "version"), path.New("json", "spec")},
				{path.New("yaml"), path.New("json", "spec", "config")},
				{path.New("yaml", "ignition"), path.New("json", "spec", "config", "ignition")},
				{path.New("yaml", "version"), path.New("json", "spec", "config", "ignition", "version")},
			},
		},
		// FIPS
		{
			Config{
				Metadata: Metadata{
					Name: "z",
					Labels: map[string]string{
						ROLE_LABEL_KEY: "z",
					},
				},
				OpenShift: OpenShift{
					FIPS: util.BoolToPtr(true),
				},
				Config: fcos.Config{
					Config: base.Config{
						Storage: base.Storage{
							Luks: []base.Luks{
								{
									Name: "a",
								},
								{
									Name:    "b",
									Options: []base.LuksOption{"b", "b"},
								},
								{
									Name:    "c",
									Options: []base.LuksOption{"c", "--cipher", "c"},
								},
								{
									Name:    "d",
									Options: []base.LuksOption{"--cipher=z"},
								},
								{
									Name:    "e",
									Options: []base.LuksOption{"-c", "z"},
								},
								{
									Name:    "f",
									Options: []base.LuksOption{"--ciphertext"},
								},
							},
						},
					},
					BootDevice: fcos.BootDevice{
						Luks: fcos.BootDeviceLuks{
							Tpm2: util.BoolToPtr(true),
						},
					},
				},
			},
			result.MachineConfig{
				ApiVersion: result.MC_API_VERSION,
				Kind:       result.MC_KIND,
				Metadata: result.Metadata{
					Name: "z",
					Labels: map[string]string{
						ROLE_LABEL_KEY: "z",
					},
				},
				Spec: result.Spec{
					Config: types.Config{
						Ignition: types.Ignition{
							Version: "3.4.0-experimental",
						},
						Storage: types.Storage{
							Filesystems: []types.Filesystem{
								{
									Device:         "/dev/mapper/root",
									Format:         util.StrToPtr("xfs"),
									Label:          util.StrToPtr("root"),
									WipeFilesystem: util.BoolToPtr(true),
								},
							},
							Luks: []types.Luks{
								{
									Name:       "root",
									Device:     util.StrToPtr("/dev/disk/by-partlabel/root"),
									Label:      util.StrToPtr("luks-root"),
									WipeVolume: util.BoolToPtr(true),
									Options:    []types.LuksOption{fipsCipherOption, fipsCipherArgument},
									Clevis: types.Clevis{
										Tpm2: util.BoolToPtr(true),
									},
								},
								{
									Name:    "a",
									Options: []types.LuksOption{fipsCipherOption, fipsCipherArgument},
								},
								{
									Name:    "b",
									Options: []types.LuksOption{"b", "b", fipsCipherOption, fipsCipherArgument},
								},
								{
									Name:    "c",
									Options: []types.LuksOption{"c", "--cipher", "c"},
								},
								{
									Name:    "d",
									Options: []types.LuksOption{"--cipher=z"},
								},
								{
									Name:    "e",
									Options: []types.LuksOption{"-c", "z"},
								},
								{
									Name:    "f",
									Options: []types.LuksOption{"--ciphertext", fipsCipherOption, fipsCipherArgument},
								},
							},
						},
					},
					FIPS: util.BoolToPtr(true),
				},
			},
			[]translate.Translation{
				{path.New("yaml", "version"), path.New("json", "apiVersion")},
				{path.New("yaml", "version"), path.New("json", "kind")},
				{path.New("yaml", "version"), path.New("json", "spec")},
				{path.New("yaml"), path.New("json", "spec", "config")},
				{path.New("yaml", "ignition"), path.New("json", "spec", "config", "ignition")},
				{path.New("yaml", "version"), path.New("json", "spec", "config", "ignition", "version")},
				{path.New("yaml", "boot_device", "luks", "tpm2"), path.New("json", "spec", "config", "storage", "luks", 0, "clevis", "tpm2")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "spec", "config", "storage", "luks", 0, "clevis")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "spec", "config", "storage", "luks", 0, "device")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "spec", "config", "storage", "luks", 0, "label")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "spec", "config", "storage", "luks", 0, "name")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "spec", "config", "storage", "luks", 0, "wipeVolume")},
				{path.New("yaml", "openshift", "fips"), path.New("json", "spec", "config", "storage", "luks", 0, "options", 0)},
				{path.New("yaml", "openshift", "fips"), path.New("json", "spec", "config", "storage", "luks", 0, "options", 1)},
				{path.New("yaml", "openshift", "fips"), path.New("json", "spec", "config", "storage", "luks", 0, "options")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "spec", "config", "storage", "luks", 0)},
				{path.New("yaml", "storage", "luks", 0, "name"), path.New("json", "spec", "config", "storage", "luks", 1, "name")},
				{path.New("yaml", "openshift", "fips"), path.New("json", "spec", "config", "storage", "luks", 1, "options", 0)},
				{path.New("yaml", "openshift", "fips"), path.New("json", "spec", "config", "storage", "luks", 1, "options", 1)},
				{path.New("yaml", "openshift", "fips"), path.New("json", "spec", "config", "storage", "luks", 1, "options")},
				{path.New("yaml", "storage", "luks", 0), path.New("json", "spec", "config", "storage", "luks", 1)},
				{path.New("yaml", "storage", "luks", 1, "name"), path.New("json", "spec", "config", "storage", "luks", 2, "name")},
				{path.New("yaml", "storage", "luks", 1, "options", 0), path.New("json", "spec", "config", "storage", "luks", 2, "options", 0)},
				{path.New("yaml", "storage", "luks", 1, "options", 1), path.New("json", "spec", "config", "storage", "luks", 2, "options", 1)},
				{path.New("yaml", "openshift", "fips"), path.New("json", "spec", "config", "storage", "luks", 2, "options", 2)},
				{path.New("yaml", "openshift", "fips"), path.New("json", "spec", "config", "storage", "luks", 2, "options", 3)},
				{path.New("yaml", "storage", "luks", 1, "options"), path.New("json", "spec", "config", "storage", "luks", 2, "options")},
				{path.New("yaml", "storage", "luks", 1), path.New("json", "spec", "config", "storage", "luks", 2)},
				{path.New("yaml", "storage", "luks", 2, "name"), path.New("json", "spec", "config", "storage", "luks", 3, "name")},
				{path.New("yaml", "storage", "luks", 2, "options", 0), path.New("json", "spec", "config", "storage", "luks", 3, "options", 0)},
				{path.New("yaml", "storage", "luks", 2, "options", 1), path.New("json", "spec", "config", "storage", "luks", 3, "options", 1)},
				{path.New("yaml", "storage", "luks", 2, "options", 2), path.New("json", "spec", "config", "storage", "luks", 3, "options", 2)},
				{path.New("yaml", "storage", "luks", 2, "options"), path.New("json", "spec", "config", "storage", "luks", 3, "options")},
				{path.New("yaml", "storage", "luks", 2), path.New("json", "spec", "config", "storage", "luks", 3)},
				{path.New("yaml", "storage", "luks", 3, "name"), path.New("json", "spec", "config", "storage", "luks", 4, "name")},
				{path.New("yaml", "storage", "luks", 3, "options", 0), path.New("json", "spec", "config", "storage", "luks", 4, "options", 0)},
				{path.New("yaml", "storage", "luks", 3, "options"), path.New("json", "spec", "config", "storage", "luks", 4, "options")},
				{path.New("yaml", "storage", "luks", 3), path.New("json", "spec", "config", "storage", "luks", 4)},
				{path.New("yaml", "storage", "luks", 4, "name"), path.New("json", "spec", "config", "storage", "luks", 5, "name")},
				{path.New("yaml", "storage", "luks", 4, "options", 0), path.New("json", "spec", "config", "storage", "luks", 5, "options", 0)},
				{path.New("yaml", "storage", "luks", 4, "options", 1), path.New("json", "spec", "config", "storage", "luks", 5, "options", 1)},
				{path.New("yaml", "storage", "luks", 4, "options"), path.New("json", "spec", "config", "storage", "luks", 5, "options")},
				{path.New("yaml", "storage", "luks", 4), path.New("json", "spec", "config", "storage", "luks", 5)},
				{path.New("yaml", "storage", "luks", 5, "name"), path.New("json", "spec", "config", "storage", "luks", 6, "name")},
				{path.New("yaml", "storage", "luks", 5, "options", 0), path.New("json", "spec", "config", "storage", "luks", 6, "options", 0)},
				{path.New("yaml", "openshift", "fips"), path.New("json", "spec", "config", "storage", "luks", 6, "options", 1)},
				{path.New("yaml", "openshift", "fips"), path.New("json", "spec", "config", "storage", "luks", 6, "options", 2)},
				{path.New("yaml", "storage", "luks", 5, "options"), path.New("json", "spec", "config", "storage", "luks", 6, "options")},
				{path.New("yaml", "storage", "luks", 5), path.New("json", "spec", "config", "storage", "luks", 6)},
				{path.New("yaml", "storage", "luks"), path.New("json", "spec", "config", "storage", "luks")},
				{path.New("yaml", "boot_device"), path.New("json", "spec", "config", "storage", "filesystems", 0, "device")},
				{path.New("yaml", "boot_device"), path.New("json", "spec", "config", "storage", "filesystems", 0, "format")},
				{path.New("yaml", "boot_device"), path.New("json", "spec", "config", "storage", "filesystems", 0, "label")},
				{path.New("yaml", "boot_device"), path.New("json", "spec", "config", "storage", "filesystems", 0, "wipeFilesystem")},
				{path.New("yaml", "boot_device"), path.New("json", "spec", "config", "storage", "filesystems", 0)},
				{path.New("yaml", "boot_device"), path.New("json", "spec", "config", "storage", "filesystems")},
				{path.New("yaml", "storage"), path.New("json", "spec", "config", "storage")},
				{path.New("yaml", "openshift", "fips"), path.New("json", "spec", "fips")},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("translate %d", i), func(t *testing.T) {
			actual, translations, r := test.in.ToMachineConfig4_11Unvalidated(common.TranslateOptions{})
			assert.Equal(t, test.out, actual, "translation mismatch")
			assert.Equal(t, report.Report{}, r, "non-empty report")
			baseutil.VerifyTranslations(t, translations, test.exceptions)
			assert.NoError(t, translations.DebugVerifyCoverage(actual), "incomplete TranslationSet coverage")
		})
	}
}

// Test post-translation validation of RHCOS/MCO support for Ignition config fields.
func TestValidateSupport(t *testing.T) {
	type entry struct {
		kind report.EntryKind
		err  error
		path path.ContextPath
	}
	tests := []struct {
		in      Config
		entries []entry
	}{
		// empty-ish config
		{
			Config{
				Metadata: Metadata{
					Name: "z",
					Labels: map[string]string{
						ROLE_LABEL_KEY: "z",
					},
				},
			},
			[]entry{},
		},
		// core user with only accepted fields
		{
			Config{
				Metadata: Metadata{
					Name: "z",
					Labels: map[string]string{
						ROLE_LABEL_KEY: "z",
					},
				},
				Config: fcos.Config{
					Config: base.Config{
						Passwd: base.Passwd{
							Users: []base.PasswdUser{
								{
									Name:              "core",
									SSHAuthorizedKeys: []base.SSHAuthorizedKey{"value"},
								},
							},
						},
					},
				},
			},
			[]entry{},
		},
		// valid data URL
		{
			Config{
				Metadata: Metadata{
					Name: "z",
					Labels: map[string]string{
						ROLE_LABEL_KEY: "z",
					},
				},
				Config: fcos.Config{
					Config: base.Config{
						Storage: base.Storage{
							Files: []base.File{
								{
									Path: "/f",
									Contents: base.Resource{
										Source: util.StrToPtr("data:,foo"),
									},
								},
							},
						},
					},
				},
			},
			[]entry{},
		},
		// all the warnings/errors
		{
			Config{
				Metadata: Metadata{
					Name: "z",
					Labels: map[string]string{
						ROLE_LABEL_KEY: "z",
					},
				},
				Config: fcos.Config{
					Config: base.Config{
						Storage: base.Storage{
							Files: []base.File{
								{
									Path: "/f",
								},
								{
									Path: "/g",
									Append: []base.Resource{
										{
											Inline: util.StrToPtr("z"),
										},
									},
								},
								{
									Path: "/h",
									Contents: base.Resource{
										Source: util.StrToPtr("https://example.com/"),
									},
								},
							},
							Filesystems: []base.Filesystem{
								{
									Device: "/dev/vda4",
									Format: util.StrToPtr("btrfs"),
								},
								{
									Device: "/dev/vda5",
									Format: util.StrToPtr("none"),
								},
							},
							Directories: []base.Directory{
								{
									Path: "/d",
								},
							},
							Links: []base.Link{
								{
									Path:   "/l",
									Target: util.StrToPtr("/t"),
								},
							},
						},
						Passwd: base.Passwd{
							Users: []base.PasswdUser{
								{
									Name:  "core",
									Gecos: util.StrToPtr("mercury delay line"),
									Groups: []base.Group{
										"z",
									},
									HomeDir:           util.StrToPtr("/home/drum"),
									NoCreateHome:      util.BoolToPtr(true),
									NoLogInit:         util.BoolToPtr(true),
									NoUserGroup:       util.BoolToPtr(true),
									PasswordHash:      util.StrToPtr("corned beef"),
									PrimaryGroup:      util.StrToPtr("wheel"),
									SSHAuthorizedKeys: []base.SSHAuthorizedKey{"value"},
									Shell:             util.StrToPtr("/bin/tcsh"),
									ShouldExist:       util.BoolToPtr(false),
									System:            util.BoolToPtr(true),
									UID:               util.IntToPtr(42),
								},
								{
									Name: "bovik",
								},
							},
							Groups: []base.PasswdGroup{
								{
									Name: "mock",
								},
							},
						},
						KernelArguments: base.KernelArguments{
							ShouldExist: []base.KernelArgument{
								"foo",
							},
							ShouldNotExist: []base.KernelArgument{
								"bar",
							},
						},
					},
				},
			},
			[]entry{
				{report.Error, common.ErrBtrfsSupport, path.New("yaml", "storage", "filesystems", 0, "format")},
				{report.Error, common.ErrFilesystemNoneSupport, path.New("yaml", "storage", "filesystems", 1, "format")},
				{report.Error, common.ErrDirectorySupport, path.New("yaml", "storage", "directories", 0)},
				{report.Error, common.ErrFileAppendSupport, path.New("yaml", "storage", "files", 1, "append")},
				{report.Error, common.ErrFileSchemeSupport, path.New("yaml", "storage", "files", 2, "contents", "source")},
				{report.Error, common.ErrLinkSupport, path.New("yaml", "storage", "links", 0)},
				{report.Error, common.ErrGroupSupport, path.New("yaml", "passwd", "groups", 0)},
				{report.Error, common.ErrUserFieldSupport, path.New("yaml", "passwd", "users", 0, "gecos")},
				{report.Error, common.ErrUserFieldSupport, path.New("yaml", "passwd", "users", 0, "groups")},
				{report.Error, common.ErrUserFieldSupport, path.New("yaml", "passwd", "users", 0, "home_dir")},
				{report.Error, common.ErrUserFieldSupport, path.New("yaml", "passwd", "users", 0, "no_create_home")},
				{report.Error, common.ErrUserFieldSupport, path.New("yaml", "passwd", "users", 0, "no_log_init")},
				{report.Error, common.ErrUserFieldSupport, path.New("yaml", "passwd", "users", 0, "no_user_group")},
				{report.Error, common.ErrUserFieldSupport, path.New("yaml", "passwd", "users", 0, "password_hash")},
				{report.Error, common.ErrUserFieldSupport, path.New("yaml", "passwd", "users", 0, "primary_group")},
				{report.Error, common.ErrUserFieldSupport, path.New("yaml", "passwd", "users", 0, "shell")},
				{report.Error, common.ErrUserFieldSupport, path.New("yaml", "passwd", "users", 0, "should_exist")},
				{report.Error, common.ErrUserFieldSupport, path.New("yaml", "passwd", "users", 0, "system")},
				{report.Error, common.ErrUserFieldSupport, path.New("yaml", "passwd", "users", 0, "uid")},
				{report.Error, common.ErrUserNameSupport, path.New("yaml", "passwd", "users", 1)},
				{report.Error, common.ErrKernelArgumentSupport, path.New("yaml", "kernel_arguments", "should_exist", 0)},
				{report.Error, common.ErrKernelArgumentSupport, path.New("yaml", "kernel_arguments", "should_not_exist", 0)},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("translate %d", i), func(t *testing.T) {
			var expectedReport report.Report
			for _, entry := range test.entries {
				expectedReport.AddOn(entry.path, entry.err, entry.kind)
			}
			actual, translations, r := test.in.ToMachineConfig4_11Unvalidated(common.TranslateOptions{})
			assert.Equal(t, expectedReport, r, "report mismatch")
			assert.NoError(t, translations.DebugVerifyCoverage(actual), "incomplete TranslationSet coverage")
		})
	}
}
