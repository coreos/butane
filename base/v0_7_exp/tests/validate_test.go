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

package tests

import (
	"fmt"
	"testing"

	baseutil "github.com/coreos/butane/base/util"
	"github.com/coreos/butane/base/v0_7_exp"
	"github.com/coreos/butane/config/common"

	"github.com/coreos/ignition/v2/config/util"
	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
	"github.com/stretchr/testify/assert"
)

// TestValidateResource tests that multiple sources (i.e. urls and inline) are not allowed but zero or one sources are
func TestValidateResource(t *testing.T) {
	tests := []struct {
		in      v0_7_exp.Resource
		out     error
		errPath path.ContextPath
	}{
		{},
		// source specified
		{
			// contains invalid (by the validator's definition) combinations of fields,
			// but the translator doesn't care and we can check they all get translated at once
			v0_7_exp.Resource{
				Source:      util.StrToPtr("http://example/com"),
				Compression: util.StrToPtr("gzip"),
				Verification: v0_7_exp.Verification{
					Hash: util.StrToPtr("this isn't validated"),
				},
			},
			nil,
			path.New("yaml"),
		},
		// inline specified
		{
			v0_7_exp.Resource{
				Inline:      util.StrToPtr("hello"),
				Compression: util.StrToPtr("gzip"),
				Verification: v0_7_exp.Verification{
					Hash: util.StrToPtr("this isn't validated"),
				},
			},
			nil,
			path.New("yaml"),
		},
		// local specified
		{
			v0_7_exp.Resource{
				Local:       util.StrToPtr("hello"),
				Compression: util.StrToPtr("gzip"),
				Verification: v0_7_exp.Verification{
					Hash: util.StrToPtr("this isn't validated"),
				},
			},
			nil,
			path.New("yaml"),
		},
		// source + inline, invalid
		{
			v0_7_exp.Resource{
				Source:      util.StrToPtr("data:,hello"),
				Inline:      util.StrToPtr("hello"),
				Compression: util.StrToPtr("gzip"),
				Verification: v0_7_exp.Verification{
					Hash: util.StrToPtr("this isn't validated"),
				},
			},
			common.ErrTooManyResourceSources,
			path.New("yaml", "source"),
		},
		// source + local, invalid
		{
			v0_7_exp.Resource{
				Source:      util.StrToPtr("data:,hello"),
				Local:       util.StrToPtr("hello"),
				Compression: util.StrToPtr("gzip"),
				Verification: v0_7_exp.Verification{
					Hash: util.StrToPtr("this isn't validated"),
				},
			},
			common.ErrTooManyResourceSources,
			path.New("yaml", "source"),
		},
		// inline + local, invalid
		{
			v0_7_exp.Resource{
				Inline:      util.StrToPtr("hello"),
				Local:       util.StrToPtr("hello"),
				Compression: util.StrToPtr("gzip"),
				Verification: v0_7_exp.Verification{
					Hash: util.StrToPtr("this isn't validated"),
				},
			},
			common.ErrTooManyResourceSources,
			path.New("yaml", "inline"),
		},
		// source + inline + local, invalid
		{
			v0_7_exp.Resource{
				Source:      util.StrToPtr("data:,hello"),
				Inline:      util.StrToPtr("hello"),
				Local:       util.StrToPtr("hello"),
				Compression: util.StrToPtr("gzip"),
				Verification: v0_7_exp.Verification{
					Hash: util.StrToPtr("this isn't validated"),
				},
			},
			common.ErrTooManyResourceSources,
			path.New("yaml", "source"),
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

func TestValidateTree(t *testing.T) {
	tests := []struct {
		in  v0_7_exp.Tree
		out error
	}{
		{
			in:  v0_7_exp.Tree{},
			out: common.ErrTreeNoLocal,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("validate %d", i), func(t *testing.T) {
			actual := test.in.Validate(path.New("yaml"))
			baseutil.VerifyReport(t, test.in, actual)
			expected := report.Report{}
			expected.AddOnError(path.New("yaml"), test.out)
			assert.Equal(t, expected, actual, "bad report")
		})
	}
}

func TestValidateFileMode(t *testing.T) {
	fileTests := []struct {
		in  v0_7_exp.File
		out error
	}{
		{
			in:  v0_7_exp.File{},
			out: nil,
		},
		{
			in: v0_7_exp.File{
				Mode: util.IntToPtr(0600),
			},
			out: nil,
		},
		{
			in: v0_7_exp.File{
				Mode: util.IntToPtr(600),
			},
			out: common.ErrDecimalMode,
		},
	}

	for i, test := range fileTests {
		t.Run(fmt.Sprintf("validate %d", i), func(t *testing.T) {
			actual := test.in.Validate(path.New("yaml"))
			baseutil.VerifyReport(t, test.in, actual)
			expected := report.Report{}
			expected.AddOnWarn(path.New("yaml", "mode"), test.out)
			assert.Equal(t, expected, actual, "bad report")
		})
	}
}

func TestValidateDirMode(t *testing.T) {
	dirTests := []struct {
		in  v0_7_exp.Directory
		out error
	}{
		{
			in:  v0_7_exp.Directory{},
			out: nil,
		},
		{
			in: v0_7_exp.Directory{
				Mode: util.IntToPtr(01770),
			},
			out: nil,
		},
		{
			in: v0_7_exp.Directory{
				Mode: util.IntToPtr(1770),
			},
			out: common.ErrDecimalMode,
		},
	}

	for i, test := range dirTests {
		t.Run(fmt.Sprintf("validate %d", i), func(t *testing.T) {
			actual := test.in.Validate(path.New("yaml"))
			baseutil.VerifyReport(t, test.in, actual)
			expected := report.Report{}
			expected.AddOnWarn(path.New("yaml", "mode"), test.out)
			assert.Equal(t, expected, actual, "bad report")
		})
	}
}

func TestValidateFilesystem(t *testing.T) {
	tests := []struct {
		in      v0_7_exp.Filesystem
		out     error
		errPath path.ContextPath
	}{
		{
			v0_7_exp.Filesystem{},
			nil,
			path.New("yaml"),
		},
		{
			v0_7_exp.Filesystem{
				Device: "/dev/foo",
			},
			nil,
			path.New("yaml"),
		},
		{
			v0_7_exp.Filesystem{
				Device:        "/dev/foo",
				Format:        util.StrToPtr("zzz"),
				Path:          util.StrToPtr("/z"),
				WithMountUnit: util.BoolToPtr(true),
			},
			nil,
			path.New("yaml"),
		},
		{
			v0_7_exp.Filesystem{
				Device:        "/dev/foo",
				Format:        util.StrToPtr("swap"),
				WithMountUnit: util.BoolToPtr(true),
			},
			nil,
			path.New("yaml"),
		},
		{
			v0_7_exp.Filesystem{
				Device:        "/dev/foo",
				WithMountUnit: util.BoolToPtr(true),
			},
			common.ErrMountUnitNoFormat,
			path.New("yaml", "format"),
		},
		{
			v0_7_exp.Filesystem{
				Device:        "/dev/foo",
				Format:        util.StrToPtr("zzz"),
				WithMountUnit: util.BoolToPtr(true),
			},
			common.ErrMountUnitNoPath,
			path.New("yaml", "path"),
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

// TestValidateUnit tests that multiple sources (i.e. contents and contents_local) are not allowed but zero or one sources are
func TestValidateUnit(t *testing.T) {
	tests := []struct {
		in      v0_7_exp.Unit
		out     error
		errPath path.ContextPath
	}{
		{},
		// contents specified
		{
			v0_7_exp.Unit{
				Contents: util.StrToPtr("hello"),
			},
			nil,
			path.New("yaml"),
		},
		// contents_local specified
		{
			v0_7_exp.Unit{
				ContentsLocal: util.StrToPtr("hello"),
			},
			nil,
			path.New("yaml"),
		},
		// contents + contents_local, invalid
		{
			v0_7_exp.Unit{
				Contents:      util.StrToPtr("hello"),
				ContentsLocal: util.StrToPtr("hello, too"),
			},
			common.ErrTooManySystemdSources,
			path.New("yaml", "contents_local"),
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

// TestValidateDropin tests that multiple sources (i.e. contents and contents_local) are not allowed but zero or one sources are
func TestValidateDropin(t *testing.T) {
	tests := []struct {
		in      v0_7_exp.Dropin
		out     error
		errPath path.ContextPath
	}{
		{},
		// contents specified
		{
			v0_7_exp.Dropin{
				Contents: util.StrToPtr("hello"),
			},
			nil,
			path.New("yaml"),
		},
		// contents_local specified
		{
			v0_7_exp.Dropin{
				ContentsLocal: util.StrToPtr("hello"),
			},
			nil,
			path.New("yaml"),
		},
		// contents + contents_local, invalid
		{
			v0_7_exp.Dropin{
				Contents:      util.StrToPtr("hello"),
				ContentsLocal: util.StrToPtr("hello, too"),
			},
			common.ErrTooManySystemdSources,
			path.New("yaml", "contents_local"),
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

// TestUnkownIgnitionVersion tests that butane will raise a warning but will not fail when an ignition config with an unkown version is specified
func TestUnkownIgnitionVersion(t *testing.T) {
	test := struct {
		in      v0_7_exp.Resource
		out     error
		errPath path.ContextPath
	}{
		v0_7_exp.Resource{
			Inline: util.StrToPtr(`{"ignition": {"version": "10.0.0"}}`),
		},
		common.ErrUnkownIgnitionVersion,
		path.New("yaml", "ignition", "config", "version"),
	}
	path := path.New("yaml", "ignition", "config")
	// Skipping baseutil.VerifyReport because it expects all referenced paths to exist in the struct.
	// In this test, "ignition.config" doesn't exist, so VerifyReport would fail. However, we still need
	// to pass this path to Validate() to trigger the unknown Ignition version warning we're testing for.
	actual := test.in.Validate(path)
	expected := report.Report{}
	expected.AddOnWarn(test.errPath, test.out)
	assert.Equal(t, expected, actual, "bad report")
}
