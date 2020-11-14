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

package v0_2

import (
	"testing"

	"github.com/coreos/fcct/config/common"

	"github.com/coreos/ignition/v2/config/util"
	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
	"github.com/stretchr/testify/assert"
)

// TestValidateResource tests that multiple sources (i.e. urls and inline) are not allowed but zero or one sources are
func TestValidateResource(t *testing.T) {
	tests := []struct {
		in      Resource
		out     error
		errPath path.ContextPath
	}{
		{},
		// source specified
		{
			// contains invalid (by the validator's definition) combinations of fields,
			// but the translator doesn't care and we can check they all get translated at once
			Resource{
				Source:      util.StrToPtr("http://example/com"),
				Compression: util.StrToPtr("gzip"),
				Verification: Verification{
					Hash: util.StrToPtr("this isn't validated"),
				},
			},
			nil,
			path.New("yaml"),
		},
		// inline specified
		{
			Resource{
				Inline:      util.StrToPtr("hello"),
				Compression: util.StrToPtr("gzip"),
				Verification: Verification{
					Hash: util.StrToPtr("this isn't validated"),
				},
			},
			nil,
			path.New("yaml"),
		},
		// local specified
		{
			Resource{
				Local:       util.StrToPtr("hello"),
				Compression: util.StrToPtr("gzip"),
				Verification: Verification{
					Hash: util.StrToPtr("this isn't validated"),
				},
			},
			nil,
			path.New("yaml"),
		},
		// source + inline, invalid
		{
			Resource{
				Source:      util.StrToPtr("data:,hello"),
				Inline:      util.StrToPtr("hello"),
				Compression: util.StrToPtr("gzip"),
				Verification: Verification{
					Hash: util.StrToPtr("this isn't validated"),
				},
			},
			common.ErrTooManyResourceSources,
			path.New("yaml", "source"),
		},
		// source + local, invalid
		{
			Resource{
				Source:      util.StrToPtr("data:,hello"),
				Local:       util.StrToPtr("hello"),
				Compression: util.StrToPtr("gzip"),
				Verification: Verification{
					Hash: util.StrToPtr("this isn't validated"),
				},
			},
			common.ErrTooManyResourceSources,
			path.New("yaml", "source"),
		},
		// inline + local, invalid
		{
			Resource{
				Inline:      util.StrToPtr("hello"),
				Local:       util.StrToPtr("hello"),
				Compression: util.StrToPtr("gzip"),
				Verification: Verification{
					Hash: util.StrToPtr("this isn't validated"),
				},
			},
			common.ErrTooManyResourceSources,
			path.New("yaml", "inline"),
		},
		// source + inline + local, invalid
		{
			Resource{
				Source:      util.StrToPtr("data:,hello"),
				Inline:      util.StrToPtr("hello"),
				Local:       util.StrToPtr("hello"),
				Compression: util.StrToPtr("gzip"),
				Verification: Verification{
					Hash: util.StrToPtr("this isn't validated"),
				},
			},
			common.ErrTooManyResourceSources,
			path.New("yaml", "source"),
		},
	}

	for i, test := range tests {
		actual := test.in.Validate(path.New("yaml"))
		expected := report.Report{}
		expected.AddOnError(test.errPath, test.out)
		assert.Equal(t, expected, actual, "#%d: bad report", i)
	}
}

func TestValidateTree(t *testing.T) {
	tests := []struct {
		in  Tree
		out error
	}{
		{
			in:  Tree{},
			out: common.ErrTreeNoLocal,
		},
	}

	for i, test := range tests {
		actual := test.in.Validate(path.New("yaml"))
		expected := report.Report{}
		expected.AddOnError(path.New("yaml"), test.out)
		assert.Equal(t, expected, actual, "#%d: bad report", i)
	}
}

func TestValidateMode(t *testing.T) {
	fileTests := []struct {
		in  File
		out error
	}{
		{
			in:  File{},
			out: nil,
		},
		{
			in: File{
				Mode: util.IntToPtr(0600),
			},
			out: nil,
		},
		{
			in: File{
				Mode: util.IntToPtr(600),
			},
			out: common.ErrDecimalMode,
		},
	}

	for i, test := range fileTests {
		actual := test.in.Validate(path.New("yaml"))
		expected := report.Report{}
		expected.AddOnWarn(path.New("yaml", "mode"), test.out)
		assert.Equal(t, expected, actual, "#%d: bad report", i)
	}

	dirTests := []struct {
		in  Directory
		out error
	}{
		{
			in:  Directory{},
			out: nil,
		},
		{
			in: Directory{
				Mode: util.IntToPtr(01770),
			},
			out: nil,
		},
		{
			in: Directory{
				Mode: util.IntToPtr(1770),
			},
			out: common.ErrDecimalMode,
		},
	}

	for i, test := range dirTests {
		actual := test.in.Validate(path.New("yaml"))
		expected := report.Report{}
		expected.AddOnWarn(path.New("yaml", "mode"), test.out)
		assert.Equal(t, expected, actual, "#%d: bad report", i)
	}
}
