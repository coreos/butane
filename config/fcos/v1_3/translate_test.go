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

	baseutil "github.com/coreos/fcct/base/util"
	base "github.com/coreos/fcct/base/v0_3"
	"github.com/coreos/fcct/config/common"
	"github.com/coreos/fcct/translate"

	"github.com/coreos/ignition/v2/config/util"
	"github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
	"github.com/stretchr/testify/assert"
)

// Most of this is covered by the Ignition translator generic tests, so just test the custom bits

// TestTranslateBootDevice tests translating the FCC boot_device section.
func TestTranslateBootDevice(t *testing.T) {
	tests := []struct {
		in         Config
		out        types.Config
		exceptions []translate.Translation
	}{
		// empty config
		{
			Config{},
			types.Config{
				Ignition: types.Ignition{
					Version: "3.2.0",
				},
			},
			[]translate.Translation{},
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
					Version: "3.2.0",
				},
				Storage: types.Storage{
					Luks: []types.Luks{
						{
							Clevis: &types.Clevis{
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
				{path.New("yaml", "boot_device", "luks", "tang", 0, "url"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "url")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "thumbprint"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "thumbprint")},
				{path.New("yaml", "boot_device", "luks", "threshold"), path.New("json", "storage", "luks", 0, "clevis", "threshold")},
				{path.New("yaml", "boot_device", "luks", "tpm2"), path.New("json", "storage", "luks", 0, "clevis", "tpm2")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "device")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "label")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "name")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "wipeVolume")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 0, "device")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 0, "format")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 0, "label")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 0, "wipeFilesystem")},
			},
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
					Version: "3.2.0",
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
								"/dev/disk/by-partlabel/esp-1",
								"/dev/disk/by-partlabel/esp-2",
								"/dev/disk/by-partlabel/esp-3",
							},
							Level:   "raid1",
							Name:    "md-esp",
							Options: []types.RaidOption{"--metadata=1.0"},
						},
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/boot-1",
								"/dev/disk/by-partlabel/boot-2",
								"/dev/disk/by-partlabel/boot-3",
							},
							Level:   "raid1",
							Name:    "md-boot",
							Options: []types.RaidOption{"--metadata=1.0"},
						},
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/root-1",
								"/dev/disk/by-partlabel/root-2",
								"/dev/disk/by-partlabel/root-3",
							},
							Level: "raid1",
							Name:  "md-root",
						},
					},
					Filesystems: []types.Filesystem{
						{
							Device:         "/dev/md/md-esp",
							Format:         util.StrToPtr("vfat"),
							Label:          util.StrToPtr("EFI-SYSTEM"),
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
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 3, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 3, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 3, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 2)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "options", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 2)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "options", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "devices", 2)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "device")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "format")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "label")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "device")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "format")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "label")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "wipeFilesystem")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 2, "device")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 2, "format")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 2, "label")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 2, "wipeFilesystem")},
			},
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
					Version: "3.2.0",
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
								"/dev/disk/by-partlabel/esp-1",
								"/dev/disk/by-partlabel/esp-2",
								"/dev/disk/by-partlabel/esp-3",
							},
							Level:   "raid1",
							Name:    "md-esp",
							Options: []types.RaidOption{"--metadata=1.0"},
						},
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/boot-1",
								"/dev/disk/by-partlabel/boot-2",
								"/dev/disk/by-partlabel/boot-3",
							},
							Level:   "raid1",
							Name:    "md-boot",
							Options: []types.RaidOption{"--metadata=1.0"},
						},
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/root-1",
								"/dev/disk/by-partlabel/root-2",
								"/dev/disk/by-partlabel/root-3",
							},
							Level: "raid1",
							Name:  "md-root",
						},
					},
					Luks: []types.Luks{
						{
							Clevis: &types.Clevis{
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
							Device:         "/dev/md/md-esp",
							Format:         util.StrToPtr("vfat"),
							Label:          util.StrToPtr("EFI-SYSTEM"),
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
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 3, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 3, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "partitions", 3, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 2), path.New("json", "storage", "disks", 2, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 2)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "options", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 2)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "options", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "devices", 2)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "name")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "url"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "url")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "thumbprint"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "thumbprint")},
				{path.New("yaml", "boot_device", "luks", "threshold"), path.New("json", "storage", "luks", 0, "clevis", "threshold")},
				{path.New("yaml", "boot_device", "luks", "tpm2"), path.New("json", "storage", "luks", 0, "clevis", "tpm2")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "device")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "label")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "name")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "wipeVolume")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "device")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "format")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "label")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "device")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "format")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "label")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "wipeFilesystem")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 2, "device")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 2, "format")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 2, "label")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 2, "wipeFilesystem")},
			},
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
					Version: "3.2.0",
				},
				Storage: types.Storage{
					Disks: []types.Disk{
						{
							Device: "/dev/vda",
							Partitions: []types.Partition{
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
								"/dev/disk/by-partlabel/esp-1",
								"/dev/disk/by-partlabel/esp-2",
							},
							Level:   "raid1",
							Name:    "md-esp",
							Options: []types.RaidOption{"--metadata=1.0"},
						},
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/boot-1",
								"/dev/disk/by-partlabel/boot-2",
							},
							Level:   "raid1",
							Name:    "md-boot",
							Options: []types.RaidOption{"--metadata=1.0"},
						},
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/root-1",
								"/dev/disk/by-partlabel/root-2",
							},
							Level: "raid1",
							Name:  "md-root",
						},
					},
					Luks: []types.Luks{
						{
							Clevis: &types.Clevis{
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
							Device:         "/dev/md/md-esp",
							Format:         util.StrToPtr("vfat"),
							Label:          util.StrToPtr("EFI-SYSTEM"),
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
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "options", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "options", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "name")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "url"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "url")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "thumbprint"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "thumbprint")},
				{path.New("yaml", "boot_device", "luks", "threshold"), path.New("json", "storage", "luks", 0, "clevis", "threshold")},
				{path.New("yaml", "boot_device", "luks", "tpm2"), path.New("json", "storage", "luks", 0, "clevis", "tpm2")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "device")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "label")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "name")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "wipeVolume")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "device")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "format")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "label")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "device")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "format")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "label")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "wipeFilesystem")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 2, "device")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 2, "format")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 2, "label")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 2, "wipeFilesystem")},
			},
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
					Version: "3.2.0",
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
									Label:    util.StrToPtr("prep-2"),
									SizeMiB:  util.IntToPtr(prepV1SizeMiB),
									TypeGUID: util.StrToPtr(prepTypeGuid),
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
								"/dev/disk/by-partlabel/esp-1",
								"/dev/disk/by-partlabel/esp-2",
							},
							Level:   "raid1",
							Name:    "md-esp",
							Options: []types.RaidOption{"--metadata=1.0"},
						},
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/boot-1",
								"/dev/disk/by-partlabel/boot-2",
							},
							Level:   "raid1",
							Name:    "md-boot",
							Options: []types.RaidOption{"--metadata=1.0"},
						},
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/root-1",
								"/dev/disk/by-partlabel/root-2",
							},
							Level: "raid1",
							Name:  "md-root",
						},
					},
					Luks: []types.Luks{
						{
							Clevis: &types.Clevis{
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
							Device:         "/dev/md/md-esp",
							Format:         util.StrToPtr("vfat"),
							Label:          util.StrToPtr("EFI-SYSTEM"),
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
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 3, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "device")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 3, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "options", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "options", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "name")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "url"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "url")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "thumbprint"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "thumbprint")},
				{path.New("yaml", "boot_device", "luks", "threshold"), path.New("json", "storage", "luks", 0, "clevis", "threshold")},
				{path.New("yaml", "boot_device", "luks", "tpm2"), path.New("json", "storage", "luks", 0, "clevis", "tpm2")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "device")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "label")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "name")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "wipeVolume")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "device")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "format")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "label")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "device")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "format")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "label")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "wipeFilesystem")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 2, "device")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 2, "format")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 2, "label")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 2, "wipeFilesystem")},
			},
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
					Version: "3.2.0",
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
								"/dev/disk/by-partlabel/esp-1",
								"/dev/disk/by-partlabel/esp-2",
							},
							Level:   "raid1",
							Name:    "md-esp",
							Options: []types.RaidOption{"--metadata=1.0"},
						},
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/boot-1",
								"/dev/disk/by-partlabel/boot-2",
							},
							Level:   "raid1",
							Name:    "md-boot",
							Options: []types.RaidOption{"--metadata=1.0"},
						},
						{
							Devices: []types.Device{
								"/dev/disk/by-partlabel/root-1",
								"/dev/disk/by-partlabel/root-2",
							},
							Level: "raid1",
							Name:  "md-root",
						},
					},
					Luks: []types.Luks{
						{
							Clevis: &types.Clevis{
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
							Device:         "/dev/md/md-esp",
							Format:         util.StrToPtr("vfat"),
							Label:          util.StrToPtr("EFI-SYSTEM"),
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
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "storage", "disks", 0, "partitions", 0, "label"), path.New("json", "storage", "disks", 0, "partitions", 3, "label")},
				{path.New("yaml", "storage", "disks", 0, "partitions", 0, "size_mib"), path.New("json", "storage", "disks", 0, "partitions", 3, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 0), path.New("json", "storage", "disks", 0, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 0, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 1, "typeGuid")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "label")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "partitions", 2, "sizeMiB")},
				{path.New("yaml", "storage", "disks", 1, "partitions", 0, "label"), path.New("json", "storage", "disks", 1, "partitions", 3, "label")},
				{path.New("yaml", "storage", "disks", 1, "partitions", 0, "size_mib"), path.New("json", "storage", "disks", 1, "partitions", 3, "sizeMiB")},
				{path.New("yaml", "boot_device", "mirror", "devices", 1), path.New("json", "storage", "disks", 1, "wipeTable")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 0, "options", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "name")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 1, "options", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "devices", 0)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "devices", 1)},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "level")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "raid", 2, "name")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "url"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "url")},
				{path.New("yaml", "boot_device", "luks", "tang", 0, "thumbprint"), path.New("json", "storage", "luks", 0, "clevis", "tang", 0, "thumbprint")},
				{path.New("yaml", "boot_device", "luks", "threshold"), path.New("json", "storage", "luks", 0, "clevis", "threshold")},
				{path.New("yaml", "boot_device", "luks", "tpm2"), path.New("json", "storage", "luks", 0, "clevis", "tpm2")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "device")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "label")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "name")},
				{path.New("yaml", "boot_device", "luks"), path.New("json", "storage", "luks", 0, "wipeVolume")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "device")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "format")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "label")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 0, "wipeFilesystem")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "device")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "format")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "label")},
				{path.New("yaml", "boot_device", "mirror"), path.New("json", "storage", "filesystems", 1, "wipeFilesystem")},
				{path.New("yaml", "storage", "filesystems", 0, "device"), path.New("json", "storage", "filesystems", 2, "device")},
				{path.New("yaml", "storage", "filesystems", 0, "format"), path.New("json", "storage", "filesystems", 2, "format")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 2, "label")},
				{path.New("yaml", "boot_device"), path.New("json", "storage", "filesystems", 2, "wipeFilesystem")},
			},
		},
	}

	// The partition sizes of existing layouts must never change, but
	// we use the constants in tests for clarity.  Ensure no one has
	// changed them.
	assert.Equal(t, biosV1SizeMiB, 1)
	assert.Equal(t, prepV1SizeMiB, 4)
	assert.Equal(t, espV1SizeMiB, 127)
	assert.Equal(t, bootV1SizeMiB, 384)

	for i, test := range tests {
		actual, translations, r := test.in.ToIgn3_2Unvalidated(common.TranslateOptions{})
		assert.Equal(t, test.out, actual, "#%d: translation mismatch", i)
		assert.Equal(t, report.Report{}, r, "#%d: non-empty report", i)
		baseutil.VerifyTranslations(t, translations, test.exceptions, "#%d", i)
	}
}
