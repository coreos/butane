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

package common

import (
	"errors"
)

var (
	// common field parsing
	ErrNoVariant      = errors.New("error parsing variant; must be specified")
	ErrInvalidVersion = errors.New("error parsing version; must be a valid semver")

	// high-level errors for fatal reports
	ErrInvalidSourceConfig    = errors.New("source config is invalid")
	ErrInvalidGeneratedConfig = errors.New("config generated was invalid")

	// resources and trees
	ErrTooManyResourceSources = errors.New("only one of the following can be set: inline, local, source")
	ErrFilesDirEscape         = errors.New("local file path traverses outside the files directory")
	ErrFileType               = errors.New("trees may only contain files, directories, and symlinks")
	ErrNodeExists             = errors.New("matching filesystem node has existing contents or different type")
	ErrNoFilesDir             = errors.New("local file paths are relative to a files directory that must be specified with -d/--files-dir")
	ErrTreeNotDirectory       = errors.New("root of tree must be a directory")
	ErrTreeNoLocal            = errors.New("local is required")

	// filesystem nodes
	ErrDecimalMode = errors.New("unreasonable mode would be reasonable if specified in octal; remember to add a leading zero")

	// mount units
	ErrMountUnitNoPath   = errors.New("path is required if with_mount_unit is true and format is not swap")
	ErrMountUnitNoFormat = errors.New("format is required if with_mount_unit is true")

	// boot device
	ErrUnknownBootDeviceLayout = errors.New("layout must be one of: aarch64, ppc64le, x86_64")
	ErrTooFewMirrorDevices     = errors.New("mirroring requires at least two devices")
)
