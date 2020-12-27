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

import "io/fs"

type TranslateOptions struct {
	// FS allows embedding local files.
	FS fs.FS
	// NoResourceAutoCompression disables automatic compression
	// of inline/local resources.
	NoResourceAutoCompression bool
	// DebugPrintTranslations reports translations to stderr.
	DebugPrintTranslations    bool
}

type TranslateBytesOptions struct {
	TranslateOptions
	Pretty bool
	Strict bool
}
