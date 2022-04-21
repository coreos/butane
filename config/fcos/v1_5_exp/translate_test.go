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

package v1_5_exp

import (
	"fmt"
	"testing"

	baseutil "github.com/coreos/butane/base/util"
	base "github.com/coreos/butane/base/v0_5_exp"
	"github.com/coreos/butane/config/common"
	"github.com/coreos/butane/translate"

	"github.com/coreos/ignition/v2/config/util"
	"github.com/coreos/ignition/v2/config/v3_4_experimental/types"
	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
	"github.com/stretchr/testify/assert"
)

// Most of this is covered by the Ignition translator generic tests, so just test the custom bits

// TestTranslateBootDevice tests translating the Butane config boot_device section.
func TestTranslateBootDevice(t *testing.T) {
	tests := []struct {
		in         Config
		out        types.Config
		exceptions []translate.Translation
		report     report.Report
	}{
		// empty config
		{
			Config{},
			types.Config{
				Ignition: types.Ignition{
					Version: "3.4.0-experimental",
				},
			},
			[]translate.Translation{
				{path.New("yaml", "version"), path.New("json", "ignition", "version")},
			},
			report.Report{},
		},
		// partition number for the `root` label is incorrect
		{
			Config{
				Config: base.Config{
					Storage: base.Storage{
						Disks: []base.Disk{
							{
								Device: "/dev/vda",
								Partitions: []base.Partition{
									{
										Label:   util.StrToPtr("root"),
										SizeMiB: util.IntToPtr(12000),
										Resize:  util.BoolToPtr(true),
									},
									{
										Label:   util.StrToPtr("var-home"),
										SizeMiB: util.IntToPtr(10240),
									},
								},
							},
						},
						Filesystems: []base.Filesystem{
							{
								Device:         "/dev/disk/by-partlabel/var-home",
								Format:         util.StrToPtr("xfs"),
								Path:           util.StrToPtr("/var/home"),
								Label:          util.StrToPtr("var-home"),
								WipeFilesystem: util.BoolToPtr(false),
							},
						},
					},
				},
			},
			types.Config{
				Ignition: types.Ignition{
					Version: "3.4.0-experimental",
				},
				Storage: types.Storage{
					Disks: []types.Disk{
						{
							Device: "/dev/vda",
							Partitions: []types.Partition{
								{
									Label:   util.StrToPtr("root"),
									SizeMiB: util.IntToPtr(12000),
									Resize:  util.BoolToPtr(true),
								},
								{
									Label:   util.StrToPtr("var-home"),
									SizeMiB: util.IntToPtr(10240),
								},
							},
						},
					},
					Filesystems: []types.Filesystem{
						{
							Device:         "/dev/disk/by-partlabel/var-home",
							Format:         util.StrToPtr("xfs"),
							Path:           util.StrToPtr("/var/home"),
							Label:          util.StrToPtr("var-home"),
							WipeFilesystem: util.BoolToPtr(false),
						},
					},
				},
			},
			[]translate.Translation{
				{path.New("yaml", "version"), path.New("json", "ignition", "version")},
				{path.New("yaml", "storage", "disks", 0, "partitions", 0, "label"), path.New("json", "storage", "disks", 0, "partitions", 0, "label")},
				{path.New("yaml", "storage", "disks", 0, "partitions", 0, "size_mib"), path.New("json", "storage", "disks", 0, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "storage", "disks", 0, "partitions", 0, "resize"), path.New("json", "storage", "disks", 0, "partitions", 0, "resize")},
				{path.New("yaml", "storage", "disks", 0, "partitions", 1, "label"), path.New("json", "storage", "disks", 0, "partitions", 1, "label")},
				{path.New("yaml", "storage", "disks", 0, "partitions", 1, "size_mib"), path.New("json", "storage", "disks", 0, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "storage", "disks", 0, "partitions", 0), path.New("json", "storage", "disks", 0, "partitions", 0)},
				{path.New("yaml", "storage", "disks", 0), path.New("json", "storage", "disks", 0)},
				{path.New("yaml", "storage", "filesystems", 0, "device"), path.New("json", "storage", "filesystems", 0, "device")},
				{path.New("yaml", "storage", "filesystems", 0, "format"), path.New("json", "storage", "filesystems", 0, "format")},
				{path.New("yaml", "storage", "filesystems", 0, "path"), path.New("json", "storage", "filesystems", 0, "path")},
				{path.New("yaml", "storage", "filesystems", 0, "label"), path.New("json", "storage", "filesystems", 0, "label")},
				{path.New("yaml", "storage", "filesystems", 0, "wipe_filesystem"), path.New("json", "storage", "filesystems", 0, "wipeFilesystem")},
				{path.New("yaml", "storage", "filesystems", 0), path.New("json", "storage", "filesystems", 0)},
				{path.New("yaml", "storage", "filesystems"), path.New("json", "storage", "filesystems")},
				{path.New("yaml", "storage"), path.New("json", "storage")},
			},
			report.Report{
				Entries: []report.Entry{
					{
						Kind:    report.Warn,
						Message: common.ErrWrongPartitionNumber.Error(),
						Context: path.New("json", "storage", "disks", 0, "partitions", 0, "label"),
					},
				},
			},
		},
		// LUKS, x86_64
		{
			Config{
				BootDevice: BootDevice{
					Luks: BootDeviceLuks{
						Tang: []base.Tang{{
							URL:        "https://example.com/",
							Thumbprint: util.StrToPtr("z"),
						}},
						Threshold: util.IntToPtr(2),
						Tpm2:      util.BoolToPtr(true),
					},
				},
			},
			types.Config{
				Ignition: types.Ignition{
					Version: "3.4.0-experimental",
				},
				Storage: types.Storage{
					Luks: []types.Luks{
						{
							Clevis: types.Clevis{
								Tang: []types.Tang{{
									URL:        "https://example.com/",
									Thumbprint: util.StrToPtr("z"),
								}},
								Threshold: util.IntToPtr(2),
								Tpm2:      util.BoolToPtr(true),
							},
							Device:     util.StrToPtr("/dev/disk/by-partlabel/root"),
							Label:      util.StrToPtr("luks-root"),
							Name:       "root",
							WipeVolume: util.BoolToPtr(true),
						},
					},
					Filesystems: []types.Filesystem{
						{
							Device:         "/dev/mapper/root",
							Format:         util.StrToPtr("xfs"),
							Label:          util.StrToPtr("root"),
							WipeFilesystem: util.BoolToPtr(true),
						},
					},
				},
			},
			[]translate.Translation{
				{path.New("yaml", "version"), path.New("json", "ignition", "version")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "url"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "url")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "thumbprint"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "thumbprint")},
				{path.New("yaml", "boot_device", "luks", "tang", 0), path.New("json", "storage", "luks", 0, "clevis", "tang", 0)},
				{path.New("yaml", "boot_device", "luks", "tang"), path.New("json", "storage", "luks", 0, "clevis", "tang")},
				{path.New("yaml", "boot_device", "luks", "threshold"), path.New("json", "storage", "luks", 0, "clevis", "threshold")},
				{path.New("yaml", "boot_device", "luks", "tpm2"), path.New("json", "storage", "luks", 0, "clevis", "tpm2")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "clevis")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "device")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "label")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "name")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "wipeVolume")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0)},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 0, "device")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 0, "format")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 0, "label")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 0, "wipeFilesystem")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 0)},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems")},
				{path.New("yaml", "boot_device"), path.New("json", "storage")},
			},
			report.Report{},
		},
		// 3-disk mirror, x86_64
		{
			Config{
				BootDevice: BootDevice{
					Mirror: BootDeviceMirror{
						Devices: []string{"/dev/vda", "/dev/vdb", "/dev/vdc"},
					},
				},
			},
			types.Config{
				Ignition: types.Ignition{
					Version: "3.4.0-experimental",
				},
				Storage: types.Storage{
					Disks: []types.Disk{
						{
							Device: "/dev/vda",
							Partitions: []types.Partition{
								{
									Label:    util.StrToPtr("bios-1"),
									SizeMiB:  util.IntToPtr(biosV1SizeMiB),
									TypeGUID: util.StrToPtr(biosTypeGuid),
								},
								{
									Label:    util.StrToPtr("esp-1"),
									SizeMiB:  util.IntToPtr(espV1SizeMiB),
									TypeGUID: util.StrToPtr(espTypeGuid),
								},
								{
									Label:   util.StrToPtr("boot-1"),
									SizeMiB: util.IntToPtr(bootV1SizeMiB),
								},
								{
									Label: util.StrToPtr("root-1"),
								},
							},
							WipeTable: util.BoolToPtr(true),
						},
						{
							Device: "/dev/vdb",
							Partitions: []types.Partition{
								{
									Label:    util.StrToPtr("bios-2"),
									SizeMiB:  util.IntToPtr(biosV1SizeMiB),
									TypeGUID: util.StrToPtr(biosTypeGuid),
								},
								{
									Label:    util.StrToPtr("esp-2"),
									SizeMiB:  util.IntToPtr(espV1SizeMiB),
									TypeGUID: util.StrToPtr(espTypeGuid),
								},
								{
									Label:   util.StrToPtr("boot-2"),
									SizeMiB: util.IntToPtr(bootV1SizeMiB),
								},
								{
									Label: util.StrToPtr("root-2"),
								},
							},
							WipeTable: util.BoolToPtr(true),
						},
						{
							Device: "/dev/vdc",
							Partitions: []types.Partition{
								{
									Label:    util.StrToPtr("bios-3"),
									SizeMiB:  util.IntToPtr(biosV1SizeMiB),
									TypeGUID: util.StrToPtr(biosTypeGuid),
								},
								{
									Label:    util.StrToPtr("esp-3"),
									SizeMiB:  util.IntToPtr(espV1SizeMiB),
									TypeGUID: util.StrToPtr(espTypeGuid),
								},
								{
									Label:   util.StrToPtr("boot-3"),
									SizeMiB: util.IntToPtr(bootV1SizeMiB),
								},
								{
									Label: util.StrToPtr("root-3"),
								},
							},
							WipeTable: util.BoolToPtr(true),
						},
					},
					Raid: []types.Raid{
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/boot-1",
								"/dev/disk/by-partlabel/boot-2",
								"/dev/disk/by-partlabel/boot-3",
							},
							Level:   util.StrToPtr("raid1"),
							Name:    "md-boot",
							Options: []types.RaidOption{"--metadata=1.0"},
						},
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/root-1",
								"/dev/disk/by-partlabel/root-2",
								"/dev/disk/by-partlabel/root-3",
							},
							Level: util.StrToPtr("raid1"),
							Name:  "md-root",
						},
					},
					Filesystems: []types.Filesystem{
						{
							Device:         "/dev/disk/by-partlabel/esp-1",
							Format:         util.StrToPtr("vfat"),
							Label:          util.StrToPtr("esp-1"),
							WipeFilesystem: util.BoolToPtr(true),
						}, {
							Device:         "/dev/disk/by-partlabel/esp-2",
							Format:         util.StrToPtr("vfat"),
							Label:          util.StrToPtr("esp-2"),
							WipeFilesystem: util.BoolToPtr(true),
						}, {
							Device:         "/dev/disk/by-partlabel/esp-3",
							Format:         util.StrToPtr("vfat"),
							Label:          util.StrToPtr("esp-3"),
							WipeFilesystem: util.BoolToPtr(true),
						}, {
							Device:         "/dev/md/md-boot",
							Format:         util.StrToPtr("ext4"),
							Label:          util.StrToPtr("boot"),
							WipeFilesystem: util.BoolToPtr(true),
						}, {
							Device:         "/dev/md/md-root",
							Format:         util.StrToPtr("xfs"),
							Label:          util.StrToPtr("root"),
							WipeFilesystem: util.BoolToPtr(true),
						},
					},
				},
			},
			[]translate.Translation{
				{path.New("yaml", "version"), path.New("json", "ignition", "version")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 3, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 3)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0, "format")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 3, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 3)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1, "format")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 2)},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 3, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 3)},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2)},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "filesystems", 2, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "filesystems", 2, "format")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "filesystems", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "filesystems", 2, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "filesystems", 2)},
				{path.New("yaml", "boot_device", "mirror", "devices"), path.New("json", "storage", "disks")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 2)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "options", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "options")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 2)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 3, "device")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 3, "format")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 3, "label")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 3, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 3)},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 4, "device")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 4, "format")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 4, "label")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 4, "wipeFilesystem")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 4)},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems")},
				{path.New("yaml", "boot_device"), path.New("json", "storage")},
			},
			report.Report{},
		},
		// 3-disk mirror + LUKS, x86_64
		{
			Config{
				BootDevice: BootDevice{
					Luks: BootDeviceLuks{
						Tang: []base.Tang{{
							URL:        "https://example.com/",
							Thumbprint: util.StrToPtr("z"),
						}},
						Threshold: util.IntToPtr(2),
						Tpm2:      util.BoolToPtr(true),
					},
					Mirror: BootDeviceMirror{
						Devices: []string{"/dev/vda", "/dev/vdb", "/dev/vdc"},
					},
				},
			},
			types.Config{
				Ignition: types.Ignition{
					Version: "3.4.0-experimental",
				},
				Storage: types.Storage{
					Disks: []types.Disk{
						{
							Device: "/dev/vda",
							Partitions: []types.Partition{
								{
									Label:    util.StrToPtr("bios-1"),
									SizeMiB:  util.IntToPtr(biosV1SizeMiB),
									TypeGUID: util.StrToPtr(biosTypeGuid),
								},
								{
									Label:    util.StrToPtr("esp-1"),
									SizeMiB:  util.IntToPtr(espV1SizeMiB),
									TypeGUID: util.StrToPtr(espTypeGuid),
								},
								{
									Label:   util.StrToPtr("boot-1"),
									SizeMiB: util.IntToPtr(bootV1SizeMiB),
								},
								{
									Label: util.StrToPtr("root-1"),
								},
							},
							WipeTable: util.BoolToPtr(true),
						},
						{
							Device: "/dev/vdb",
							Partitions: []types.Partition{
								{
									Label:    util.StrToPtr("bios-2"),
									SizeMiB:  util.IntToPtr(biosV1SizeMiB),
									TypeGUID: util.StrToPtr(biosTypeGuid),
								},
								{
									Label:    util.StrToPtr("esp-2"),
									SizeMiB:  util.IntToPtr(espV1SizeMiB),
									TypeGUID: util.StrToPtr(espTypeGuid),
								},
								{
									Label:   util.StrToPtr("boot-2"),
									SizeMiB: util.IntToPtr(bootV1SizeMiB),
								},
								{
									Label: util.StrToPtr("root-2"),
								},
							},
							WipeTable: util.BoolToPtr(true),
						},
						{
							Device: "/dev/vdc",
							Partitions: []types.Partition{
								{
									Label:    util.StrToPtr("bios-3"),
									SizeMiB:  util.IntToPtr(biosV1SizeMiB),
									TypeGUID: util.StrToPtr(biosTypeGuid),
								},
								{
									Label:    util.StrToPtr("esp-3"),
									SizeMiB:  util.IntToPtr(espV1SizeMiB),
									TypeGUID: util.StrToPtr(espTypeGuid),
								},
								{
									Label:   util.StrToPtr("boot-3"),
									SizeMiB: util.IntToPtr(bootV1SizeMiB),
								},
								{
									Label: util.StrToPtr("root-3"),
								},
							},
							WipeTable: util.BoolToPtr(true),
						},
					},
					Raid: []types.Raid{
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/boot-1",
								"/dev/disk/by-partlabel/boot-2",
								"/dev/disk/by-partlabel/boot-3",
							},
							Level:   util.StrToPtr("raid1"),
							Name:    "md-boot",
							Options: []types.RaidOption{"--metadata=1.0"},
						},
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/root-1",
								"/dev/disk/by-partlabel/root-2",
								"/dev/disk/by-partlabel/root-3",
							},
							Level: util.StrToPtr("raid1"),
							Name:  "md-root",
						},
					},
					Luks: []types.Luks{
						{
							Clevis: types.Clevis{
								Tang: []types.Tang{{
									URL:        "https://example.com/",
									Thumbprint: util.StrToPtr("z"),
								}},
								Threshold: util.IntToPtr(2),
								Tpm2:      util.BoolToPtr(true),
							},
							Device:     util.StrToPtr("/dev/md/md-root"),
							Label:      util.StrToPtr("luks-root"),
							Name:       "root",
							WipeVolume: util.BoolToPtr(true),
						},
					},
					Filesystems: []types.Filesystem{
						{
							Device:         "/dev/disk/by-partlabel/esp-1",
							Format:         util.StrToPtr("vfat"),
							Label:          util.StrToPtr("esp-1"),
							WipeFilesystem: util.BoolToPtr(true),
						}, {
							Device:         "/dev/disk/by-partlabel/esp-2",
							Format:         util.StrToPtr("vfat"),
							Label:          util.StrToPtr("esp-2"),
							WipeFilesystem: util.BoolToPtr(true),
						}, {
							Device:         "/dev/disk/by-partlabel/esp-3",
							Format:         util.StrToPtr("vfat"),
							Label:          util.StrToPtr("esp-3"),
							WipeFilesystem: util.BoolToPtr(true),
						}, {
							Device:         "/dev/md/md-boot",
							Format:         util.StrToPtr("ext4"),
							Label:          util.StrToPtr("boot"),
							WipeFilesystem: util.BoolToPtr(true),
						}, {
							Device:         "/dev/mapper/root",
							Format:         util.StrToPtr("xfs"),
							Label:          util.StrToPtr("root"),
							WipeFilesystem: util.BoolToPtr(true),
						},
					},
				},
			},
			[]translate.Translation{
				{path.New("yaml", "version"), path.New("json", "ignition", "version")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 3, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 3)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0, "format")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 3, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 3)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1, "format")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 2)},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 3, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 3)},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2)},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "filesystems", 2, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "filesystems", 2, "format")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "filesystems", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "filesystems", 2, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "filesystems", 2)},
				{path.New("yaml", "boot_device", "mirror", "devices"), path.New("json", "storage", "disks")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 2)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "options", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "options")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 2)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "url"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "url")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "thumbprint"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "thumbprint")},
				{path.New("yaml", "boot_device", "luks", "tang", 0), path.New("json", "storage", "luks", 0, "clevis", "tang", 0)},
				{path.New("yaml", "boot_device", "luks", "tang"), path.New("json", "storage", "luks", 0, "clevis", "tang")},
				{path.New("yaml", "boot_device", "luks", "threshold"), path.New("json", "storage", "luks", 0, "clevis", "threshold")},
				{path.New("yaml", "boot_device", "luks", "tpm2"), path.New("json", "storage", "luks", 0, "clevis", "tpm2")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "clevis")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "device")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "label")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "name")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "wipeVolume")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0)},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 3, "device")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 3, "format")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 3, "label")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 3, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 3)},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 4, "device")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 4, "format")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 4, "label")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 4, "wipeFilesystem")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 4)},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems")},
				{path.New("yaml", "boot_device"), path.New("json", "storage")},
			},
			report.Report{},
		},
		// 2-disk mirror + LUKS, aarch64
		{
			Config{
				BootDevice: BootDevice{
					Layout: util.StrToPtr("aarch64"),
					Luks: BootDeviceLuks{
						Tang: []base.Tang{{
							URL:        "https://example.com/",
							Thumbprint: util.StrToPtr("z"),
						}},
						Threshold: util.IntToPtr(2),
						Tpm2:      util.BoolToPtr(true),
					},
					Mirror: BootDeviceMirror{
						Devices: []string{"/dev/vda", "/dev/vdb"},
					},
				},
			},
			types.Config{
				Ignition: types.Ignition{
					Version: "3.4.0-experimental",
				},
				Storage: types.Storage{
					Disks: []types.Disk{
						{
							Device: "/dev/vda",
							Partitions: []types.Partition{
								{
									Label:    util.StrToPtr("reserved-1"),
									SizeMiB:  util.IntToPtr(reservedV1SizeMiB),
									TypeGUID: util.StrToPtr(reservedTypeGuid),
								},
								{
									Label:    util.StrToPtr("esp-1"),
									SizeMiB:  util.IntToPtr(espV1SizeMiB),
									TypeGUID: util.StrToPtr(espTypeGuid),
								},
								{
									Label:   util.StrToPtr("boot-1"),
									SizeMiB: util.IntToPtr(bootV1SizeMiB),
								},
								{
									Label: util.StrToPtr("root-1"),
								},
							},
							WipeTable: util.BoolToPtr(true),
						},
						{
							Device: "/dev/vdb",
							Partitions: []types.Partition{
								{
									Label:    util.StrToPtr("reserved-2"),
									SizeMiB:  util.IntToPtr(reservedV1SizeMiB),
									TypeGUID: util.StrToPtr(reservedTypeGuid),
								},
								{
									Label:    util.StrToPtr("esp-2"),
									SizeMiB:  util.IntToPtr(espV1SizeMiB),
									TypeGUID: util.StrToPtr(espTypeGuid),
								},
								{
									Label:   util.StrToPtr("boot-2"),
									SizeMiB: util.IntToPtr(bootV1SizeMiB),
								},
								{
									Label: util.StrToPtr("root-2"),
								},
							},
							WipeTable: util.BoolToPtr(true),
						},
					},
					Raid: []types.Raid{
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/boot-1",
								"/dev/disk/by-partlabel/boot-2",
							},
							Level:   util.StrToPtr("raid1"),
							Name:    "md-boot",
							Options: []types.RaidOption{"--metadata=1.0"},
						},
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/root-1",
								"/dev/disk/by-partlabel/root-2",
							},
							Level: util.StrToPtr("raid1"),
							Name:  "md-root",
						},
					},
					Luks: []types.Luks{
						{
							Clevis: types.Clevis{
								Tang: []types.Tang{{
									URL:        "https://example.com/",
									Thumbprint: util.StrToPtr("z"),
								}},
								Threshold: util.IntToPtr(2),
								Tpm2:      util.BoolToPtr(true),
							},
							Device:     util.StrToPtr("/dev/md/md-root"),
							Label:      util.StrToPtr("luks-root"),
							Name:       "root",
							WipeVolume: util.BoolToPtr(true),
						},
					},
					Filesystems: []types.Filesystem{
						{
							Device:         "/dev/disk/by-partlabel/esp-1",
							Format:         util.StrToPtr("vfat"),
							Label:          util.StrToPtr("esp-1"),
							WipeFilesystem: util.BoolToPtr(true),
						}, {
							Device:         "/dev/disk/by-partlabel/esp-2",
							Format:         util.StrToPtr("vfat"),
							Label:          util.StrToPtr("esp-2"),
							WipeFilesystem: util.BoolToPtr(true),
						}, {
							Device:         "/dev/md/md-boot",
							Format:         util.StrToPtr("ext4"),
							Label:          util.StrToPtr("boot"),
							WipeFilesystem: util.BoolToPtr(true),
						}, {
							Device:         "/dev/mapper/root",
							Format:         util.StrToPtr("xfs"),
							Label:          util.StrToPtr("root"),
							WipeFilesystem: util.BoolToPtr(true),
						},
					},
				},
			},
			[]translate.Translation{
				{path.New("yaml", "version"), path.New("json", "ignition", "version")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 3, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 3)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0, "format")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 3, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 3)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1, "format")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices"), path.New("json", "storage", "disks")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "options", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "options")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "url"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "url")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "thumbprint"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "thumbprint")},
				{path.New("yaml", "boot_device", "luks", "tang", 0), path.New("json", "storage", "luks", 0, "clevis", "tang", 0)},
				{path.New("yaml", "boot_device", "luks", "tang"), path.New("json", "storage", "luks", 0, "clevis", "tang")},
				{path.New("yaml", "boot_device", "luks", "threshold"), path.New("json", "storage", "luks", 0, "clevis", "threshold")},
				{path.New("yaml", "boot_device", "luks", "tpm2"), path.New("json", "storage", "luks", 0, "clevis", "tpm2")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "clevis")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "device")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "label")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "name")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "wipeVolume")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0)},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 2, "device")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 2, "format")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 2, "label")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 2, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 2)},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 3, "device")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 3, "format")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 3, "label")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 3, "wipeFilesystem")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 3)},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems")},
				{path.New("yaml", "boot_device"), path.New("json", "storage")},
			},
			report.Report{},
		},
		// 2-disk mirror + LUKS, ppc64le
		{
			Config{
				BootDevice: BootDevice{
					Layout: util.StrToPtr("ppc64le"),
					Luks: BootDeviceLuks{
						Tang: []base.Tang{{
							URL:        "https://example.com/",
							Thumbprint: util.StrToPtr("z"),
						}},
						Threshold: util.IntToPtr(2),
						Tpm2:      util.BoolToPtr(true),
					},
					Mirror: BootDeviceMirror{
						Devices: []string{"/dev/vda", "/dev/vdb"},
					},
				},
			},
			types.Config{
				Ignition: types.Ignition{
					Version: "3.4.0-experimental",
				},
				Storage: types.Storage{
					Disks: []types.Disk{
						{
							Device: "/dev/vda",
							Partitions: []types.Partition{
								{
									Label:    util.StrToPtr("prep-1"),
									SizeMiB:  util.IntToPtr(prepV1SizeMiB),
									TypeGUID: util.StrToPtr(prepTypeGuid),
								},
								{
									Label:    util.StrToPtr("reserved-1"),
									SizeMiB:  util.IntToPtr(reservedV1SizeMiB),
									TypeGUID: util.StrToPtr(reservedTypeGuid),
								},
								{
									Label:   util.StrToPtr("boot-1"),
									SizeMiB: util.IntToPtr(bootV1SizeMiB),
								},
								{
									Label: util.StrToPtr("root-1"),
								},
							},
							WipeTable: util.BoolToPtr(true),
						},
						{
							Device: "/dev/vdb",
							Partitions: []types.Partition{
								{
									Label:    util.StrToPtr("prep-2"),
									SizeMiB:  util.IntToPtr(prepV1SizeMiB),
									TypeGUID: util.StrToPtr(prepTypeGuid),
								},
								{
									Label:    util.StrToPtr("reserved-2"),
									SizeMiB:  util.IntToPtr(reservedV1SizeMiB),
									TypeGUID: util.StrToPtr(reservedTypeGuid),
								},
								{
									Label:   util.StrToPtr("boot-2"),
									SizeMiB: util.IntToPtr(bootV1SizeMiB),
								},
								{
									Label: util.StrToPtr("root-2"),
								},
							},
							WipeTable: util.BoolToPtr(true),
						},
					},
					Raid: []types.Raid{
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/boot-1",
								"/dev/disk/by-partlabel/boot-2",
							},
							Level:   util.StrToPtr("raid1"),
							Name:    "md-boot",
							Options: []types.RaidOption{"--metadata=1.0"},
						},
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/root-1",
								"/dev/disk/by-partlabel/root-2",
							},
							Level: util.StrToPtr("raid1"),
							Name:  "md-root",
						},
					},
					Luks: []types.Luks{
						{
							Clevis: types.Clevis{
								Tang: []types.Tang{{
									URL:        "https://example.com/",
									Thumbprint: util.StrToPtr("z"),
								}},
								Threshold: util.IntToPtr(2),
								Tpm2:      util.BoolToPtr(true),
							},
							Device:     util.StrToPtr("/dev/md/md-root"),
							Label:      util.StrToPtr("luks-root"),
							Name:       "root",
							WipeVolume: util.BoolToPtr(true),
						},
					},
					Filesystems: []types.Filesystem{
						{
							Device:         "/dev/md/md-boot",
							Format:         util.StrToPtr("ext4"),
							Label:          util.StrToPtr("boot"),
							WipeFilesystem: util.BoolToPtr(true),
						}, {
							Device:         "/dev/mapper/root",
							Format:         util.StrToPtr("xfs"),
							Label:          util.StrToPtr("root"),
							WipeFilesystem: util.BoolToPtr(true),
						},
					},
				},
			},
			[]translate.Translation{
				{path.New("yaml", "version"), path.New("json", "ignition", "version")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 3, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 3)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 3, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 3)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices"), path.New("json", "storage", "disks")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "options", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "options")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "url"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "url")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "thumbprint"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "thumbprint")},
				{path.New("yaml", "boot_device", "luks", "tang", 0), path.New("json", "storage", "luks", 0, "clevis", "tang", 0)},
				{path.New("yaml", "boot_device", "luks", "tang"), path.New("json", "storage", "luks", 0, "clevis", "tang")},
				{path.New("yaml", "boot_device", "luks", "threshold"), path.New("json", "storage", "luks", 0, "clevis", "threshold")},
				{path.New("yaml", "boot_device", "luks", "tpm2"), path.New("json", "storage", "luks", 0, "clevis", "tpm2")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "clevis")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "device")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "label")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "name")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "wipeVolume")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0)},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "device")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "format")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "label")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0)},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 1, "device")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 1, "format")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 1, "label")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 1, "wipeFilesystem")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 1)},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems")},
				{path.New("yaml", "boot_device"), path.New("json", "storage")},
			},
			report.Report{},
		},
		// 2-disk mirror + LUKS with overridden root partition size
		// and filesystem type, x86_64
		{
			Config{
				Config: base.Config{
					Storage: base.Storage{
						Disks: []base.Disk{
							{
								Device: "/dev/vda",
								Partitions: []base.Partition{
									{
										Label:   util.StrToPtr("root-1"),
										SizeMiB: util.IntToPtr(8192),
									},
								},
							},
							{
								Device: "/dev/vdb",
								Partitions: []base.Partition{
									{
										Label:   util.StrToPtr("root-2"),
										SizeMiB: util.IntToPtr(8192),
									},
								},
							},
						},
						Filesystems: []base.Filesystem{
							{
								Device: "/dev/mapper/root",
								Format: util.StrToPtr("ext4"),
							},
						},
					},
				},
				BootDevice: BootDevice{
					Luks: BootDeviceLuks{
						Tang: []base.Tang{{
							URL:        "https://example.com/",
							Thumbprint: util.StrToPtr("z"),
						}},
						Threshold: util.IntToPtr(2),
						Tpm2:      util.BoolToPtr(true),
					},
					Mirror: BootDeviceMirror{
						Devices: []string{"/dev/vda", "/dev/vdb"},
					},
				},
			},
			types.Config{
				Ignition: types.Ignition{
					Version: "3.4.0-experimental",
				},
				Storage: types.Storage{
					Disks: []types.Disk{
						{
							Device: "/dev/vda",
							Partitions: []types.Partition{
								{
									Label:    util.StrToPtr("bios-1"),
									SizeMiB:  util.IntToPtr(biosV1SizeMiB),
									TypeGUID: util.StrToPtr(biosTypeGuid),
								},
								{
									Label:    util.StrToPtr("esp-1"),
									SizeMiB:  util.IntToPtr(espV1SizeMiB),
									TypeGUID: util.StrToPtr(espTypeGuid),
								},
								{
									Label:   util.StrToPtr("boot-1"),
									SizeMiB: util.IntToPtr(bootV1SizeMiB),
								},
								{
									Label:   util.StrToPtr("root-1"),
									SizeMiB: util.IntToPtr(8192),
								},
							},
							WipeTable: util.BoolToPtr(true),
						},
						{
							Device: "/dev/vdb",
							Partitions: []types.Partition{
								{
									Label:    util.StrToPtr("bios-2"),
									SizeMiB:  util.IntToPtr(biosV1SizeMiB),
									TypeGUID: util.StrToPtr(biosTypeGuid),
								},
								{
									Label:    util.StrToPtr("esp-2"),
									SizeMiB:  util.IntToPtr(espV1SizeMiB),
									TypeGUID: util.StrToPtr(espTypeGuid),
								},
								{
									Label:   util.StrToPtr("boot-2"),
									SizeMiB: util.IntToPtr(bootV1SizeMiB),
								},
								{
									Label:   util.StrToPtr("root-2"),
									SizeMiB: util.IntToPtr(8192),
								},
							},
							WipeTable: util.BoolToPtr(true),
						},
					},
					Raid: []types.Raid{
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/boot-1",
								"/dev/disk/by-partlabel/boot-2",
							},
							Level:   util.StrToPtr("raid1"),
							Name:    "md-boot",
							Options: []types.RaidOption{"--metadata=1.0"},
						},
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/root-1",
								"/dev/disk/by-partlabel/root-2",
							},
							Level: util.StrToPtr("raid1"),
							Name:  "md-root",
						},
					},
					Luks: []types.Luks{
						{
							Clevis: types.Clevis{
								Tang: []types.Tang{{
									URL:        "https://example.com/",
									Thumbprint: util.StrToPtr("z"),
								}},
								Threshold: util.IntToPtr(2),
								Tpm2:      util.BoolToPtr(true),
							},
							Device:     util.StrToPtr("/dev/md/md-root"),
							Label:      util.StrToPtr("luks-root"),
							Name:       "root",
							WipeVolume: util.BoolToPtr(true),
						},
					},
					Filesystems: []types.Filesystem{
						{
							Device:         "/dev/disk/by-partlabel/esp-1",
							Format:         util.StrToPtr("vfat"),
							Label:          util.StrToPtr("esp-1"),
							WipeFilesystem: util.BoolToPtr(true),
						}, {
							Device:         "/dev/disk/by-partlabel/esp-2",
							Format:         util.StrToPtr("vfat"),
							Label:          util.StrToPtr("esp-2"),
							WipeFilesystem: util.BoolToPtr(true),
						}, {
							Device:         "/dev/md/md-boot",
							Format:         util.StrToPtr("ext4"),
							Label:          util.StrToPtr("boot"),
							WipeFilesystem: util.BoolToPtr(true),
						}, {
							Device:         "/dev/mapper/root",
							Format:         util.StrToPtr("ext4"),
							Label:          util.StrToPtr("root"),
							WipeFilesystem: util.BoolToPtr(true),
						},
					},
				},
			},
			[]translate.Translation{
				{path.New("yaml", "version"), path.New("json", "ignition", "version")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2)},
				{path.New("yaml", "storage", "disks", 0, "partitions", 0, "label"), path.New("json", "storage", "disks", 0, "partitions", 3, "label")},
				{path.New("yaml", "storage", "disks", 0, "partitions", 0, "size_mib"), path.New("json", "storage", "disks", 0, "partitions", 3, "sizeMiB")},
				{path.New("yaml", "storage", "disks", 0, "partitions", 0), path.New("json", "storage", "disks", 0, "partitions", 3)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "wipeTable")},
				{path.New("yaml", "storage", "disks", 0), path.New("json", "storage", "disks", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0, "format")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "filesystems", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2)},
				{path.New("yaml", "storage", "disks", 1, "partitions", 0, "label"), path.New("json", "storage", "disks", 1, "partitions", 3, "label")},
				{path.New("yaml", "storage", "disks", 1, "partitions", 0, "size_mib"), path.New("json", "storage", "disks", 1, "partitions", 3, "sizeMiB")},
				{path.New("yaml", "storage", "disks", 1, "partitions", 0), path.New("json", "storage", "disks", 1, "partitions", 3)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "wipeTable")},
				{path.New("yaml", "storage", "disks", 1), path.New("json", "storage", "disks", 1)},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1, "format")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "filesystems", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "options", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "options")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "url"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "url")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "thumbprint"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "thumbprint")},
				{path.New("yaml", "boot_device", "luks", "tang", 0), path.New("json", "storage", "luks", 0, "clevis", "tang", 0)},
				{path.New("yaml", "boot_device", "luks", "tang"), path.New("json", "storage", "luks", 0, "clevis", "tang")},
				{path.New("yaml", "boot_device", "luks", "threshold"), path.New("json", "storage", "luks", 0, "clevis", "threshold")},
				{path.New("yaml", "boot_device", "luks", "tpm2"), path.New("json", "storage", "luks", 0, "clevis", "tpm2")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "clevis")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "device")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "label")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "name")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "wipeVolume")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0)},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 2, "device")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 2, "format")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 2, "label")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 2, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 2)},
				{path.New("yaml", "storage", "filesystems", 0, "device"), path.New("json", "storage", "filesystems", 3, "device")},
				{path.New("yaml", "storage", "filesystems", 0, "format"), path.New("json", "storage", "filesystems", 3, "format")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 3, "label")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 3, "wipeFilesystem")},
				{path.New("yaml", "storage", "filesystems", 0), path.New("json", "storage", "filesystems", 3)},
				{path.New("yaml", "storage", "filesystems"), path.New("json", "storage", "filesystems")},
				{path.New("yaml", "storage"), path.New("json", "storage")},
			},
			report.Report{},
		},
	}

	// The partition sizes of existing layouts must never change, but
	// we use the constants in tests for clarity.  Ensure no one has
	// changed them.
	assert.Equal(t, reservedV1SizeMiB, 1)
	assert.Equal(t, biosV1SizeMiB, 1)
	assert.Equal(t, prepV1SizeMiB, 4)
	assert.Equal(t, espV1SizeMiB, 127)
	assert.Equal(t, bootV1SizeMiB, 384)

	for i, test := range tests {
		t.Run(fmt.Sprintf("translate %d", i), func(t *testing.T) {
			actual, translations, r := test.in.ToIgn3_4Unvalidated(common.TranslateOptions{})
			assert.Equal(t, test.out, actual, "translation mismatch")
			assert.Equal(t, test.report, r, "report mismatch")
			baseutil.VerifyTranslations(t, translations, test.exceptions)
			assert.NoError(t, translations.DebugVerifyCoverage(actual), "incomplete TranslationSet coverage")
		})
	}
}

// TestTranslateExtensions tests translating the Butane config extensions section.
func TestTranslateExtensions(t *testing.T) {
	tests := []struct {
		in         Config
		out        types.Config
		exceptions []translate.Translation
		report     report.Report
	}{
		// config with two extensions/packages
		{
			Config{
				Extensions: []Extension{
					{
						Name: "strace",
					},
					{
						Name: "zsh",
					},
				},
			},
			types.Config{
				Ignition: types.Ignition{
					Version: "3.4.0-experimental",
				},
				Storage: types.Storage{
					Files: []types.File{
						{
							Node: types.Node{
								Path: "/etc/rpm-ostree/origin.d/extensions-e2ecf66.yaml",
							},
							FileEmbedded1: types.FileEmbedded1{
								Contents: types.Resource{
									Source:      util.StrToPtr("data:;base64,IyBHZW5lcmF0ZWQgYnkgQnV0YW5lCgpwYWNrYWdlczoKICAgIC0gc3RyYWNlCiAgICAtIHpzaAo="),
									Compression: util.StrToPtr(""),
								},
								Mode: util.IntToPtr(420),
							},
						},
					},
				},
			},
			[]translate.Translation{
				{path.New("yaml", "version"), path.New("json", "ignition", "version")},
				{path.New("yaml", "extensions"), path.New("json", "storage")},
				{path.New("yaml", "extensions"), path.New("json", "storage", "files")},
				{path.New("yaml", "extensions"), path.New("json", "storage", "files", 0)},
				{path.New("yaml", "extensions"), path.New("json", "storage", "files", 0, "path")},
				{path.New("yaml", "extensions"), path.New("json", "storage", "files", 0, "mode")},
				{path.New("yaml", "extensions"), path.New("json", "storage", "files", 0, "contents")},
				{path.New("yaml", "extensions"), path.New("json", "storage", "files", 0, "contents", "source")},
				{path.New("yaml", "extensions"), path.New("json", "storage", "files", 0, "contents", "compression")},
			},
			report.Report{},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("translate %d", i), func(t *testing.T) {
			actual, translations, r := test.in.ToIgn3_4Unvalidated(common.TranslateOptions{})
			assert.Equal(t, test.out, actual, "translation mismatch")
			assert.Equal(t, test.report, r, "report mismatch")
			baseutil.VerifyTranslations(t, translations, test.exceptions)
			assert.NoError(t, translations.DebugVerifyCoverage(actual), "incomplete TranslationSet coverage")
		})
	}
}
