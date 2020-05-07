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

package v0_2_exp

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/coreos/fcct/base"
	"github.com/coreos/fcct/translate"

	"github.com/coreos/ignition/v2/config/util"
	"github.com/coreos/ignition/v2/config/v3_1_experimental/types"
	"github.com/coreos/vcontext/path"
)

// Most of this is covered by the Ignition translator generic tests, so just test the custom bits

// verifyTranslations ensures all the translations are identity, unless they match a listed one,
// and verifies that all the listed ones exist.
// it returns the offending translation if there is one
func verifyTranslations(set translate.TranslationSet, exceptions ...translate.Translation) *translate.Translation {
	exceptionSet := translate.TranslationSet{
		FromTag: set.FromTag,
		ToTag:   set.ToTag,
		Set:     map[string]translate.Translation{},
	}
	for _, ex := range exceptions {
		exceptionSet.AddTranslation(ex.From, ex.To)
		if tr, ok := set.Set[ex.To.String()]; ok {
			if !reflect.DeepEqual(tr, ex) {
				return &ex
			}
		} else {
			return &ex
		}
	}
	for key, translation := range set.Set {
		if ex, ok := exceptionSet.Set[key]; ok {
			if !reflect.DeepEqual(translation, ex) {
				return &ex
			}
		} else if !reflect.DeepEqual(translation.From.Path, translation.To.Path) {
			return &translation
		}
	}
	return nil
}

// TestTranslateFile tests translating the ct storage.files.[i] entries to ignition storage.files.[i] entries.
func TestTranslateFile(t *testing.T) {
	filesDir, err := ioutil.TempDir("", "translate-test-")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.RemoveAll(filesDir)
	f, err := os.OpenFile(filepath.Join(filesDir, "file-1"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Error(err)
		return
	}
	// no defer; we close immediately after writing
	_, err = f.WriteString("file contents\n")
	f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		in         File
		out        types.File
		exceptions []translate.Translation
		report     string
		options    base.TranslateOptions
	}{
		{
			File{},
			types.File{},
			nil,
			"",
			base.TranslateOptions{},
		},
		{
			// contains invalid (by the validator's definition) combinations of fields,
			// but the translator doesn't care and we can check they all get translated at once
			File{
				Path: "/foo",
				Group: NodeGroup{
					ID:   util.IntToPtr(1),
					Name: util.StrToPtr("foobar"),
				},
				User: NodeUser{
					ID:   util.IntToPtr(1),
					Name: util.StrToPtr("bazquux"),
				},
				Mode: util.IntToPtr(420),
				Append: []Resource{
					{
						Source:      util.StrToPtr("http://example/com"),
						Compression: util.StrToPtr("gzip"),
						HTTPHeaders: HTTPHeaders{
							HTTPHeader{
								Name:  "Header",
								Value: util.StrToPtr("this isn't validated"),
							},
						},
						Verification: Verification{
							Hash: util.StrToPtr("this isn't validated"),
						},
					},
					{
						Inline:      util.StrToPtr("hello"),
						Compression: util.StrToPtr("gzip"),
						HTTPHeaders: HTTPHeaders{
							HTTPHeader{
								Name:  "Header",
								Value: util.StrToPtr("this isn't validated"),
							},
						},
						Verification: Verification{
							Hash: util.StrToPtr("this isn't validated"),
						},
					},
					{
						Local: util.StrToPtr("file-1"),
					},
				},
				Overwrite: util.BoolToPtr(true),
				Contents: Resource{
					Source:      util.StrToPtr("http://example/com"),
					Compression: util.StrToPtr("gzip"),
					HTTPHeaders: HTTPHeaders{
						HTTPHeader{
							Name:  "Header",
							Value: util.StrToPtr("this isn't validated"),
						},
					},
					Verification: Verification{
						Hash: util.StrToPtr("this isn't validated"),
					},
				},
			},
			types.File{
				Node: types.Node{
					Path: "/foo",
					Group: types.NodeGroup{
						ID:   util.IntToPtr(1),
						Name: util.StrToPtr("foobar"),
					},
					User: types.NodeUser{
						ID:   util.IntToPtr(1),
						Name: util.StrToPtr("bazquux"),
					},
					Overwrite: util.BoolToPtr(true),
				},
				FileEmbedded1: types.FileEmbedded1{
					Mode: util.IntToPtr(420),
					Append: []types.Resource{
						{
							Source:      util.StrToPtr("http://example/com"),
							Compression: util.StrToPtr("gzip"),
							HTTPHeaders: types.HTTPHeaders{
								types.HTTPHeader{
									Name:  "Header",
									Value: util.StrToPtr("this isn't validated"),
								},
							},
							Verification: types.Verification{
								Hash: util.StrToPtr("this isn't validated"),
							},
						},
						{
							Source:      util.StrToPtr("data:,hello"),
							Compression: util.StrToPtr("gzip"),
							HTTPHeaders: types.HTTPHeaders{
								types.HTTPHeader{
									Name:  "Header",
									Value: util.StrToPtr("this isn't validated"),
								},
							},
							Verification: types.Verification{
								Hash: util.StrToPtr("this isn't validated"),
							},
						},
						{
							Source: util.StrToPtr("data:,file%20contents%0A"),
						},
					},
					Contents: types.Resource{
						Source:      util.StrToPtr("http://example/com"),
						Compression: util.StrToPtr("gzip"),
						HTTPHeaders: types.HTTPHeaders{
							types.HTTPHeader{
								Name:  "Header",
								Value: util.StrToPtr("this isn't validated"),
							},
						},
						Verification: types.Verification{
							Hash: util.StrToPtr("this isn't validated"),
						},
					},
				},
			},
			[]translate.Translation{
				{
					From: path.New("yaml", "append", 1, "inline"),
					To:   path.New("json", "append", 1, "source"),
				},
				{
					From: path.New("yaml", "append", 2, "local"),
					To:   path.New("json", "append", 2, "source"),
				},
			},
			"",
			base.TranslateOptions{
				FilesDir: filesDir,
			},
		},
		// inline file contents
		{
			File{
				Path: "/foo",
				Contents: Resource{
					Inline: util.StrToPtr("xyzzy"),
				},
			},
			types.File{
				Node: types.Node{
					Path: "/foo",
				},
				FileEmbedded1: types.FileEmbedded1{
					Contents: types.Resource{
						Source: util.StrToPtr("data:,xyzzy"),
					},
				},
			},
			[]translate.Translation{
				{
					From: path.New("yaml", "contents", "inline"),
					To:   path.New("json", "contents", "source"),
				},
			},
			"",
			base.TranslateOptions{},
		},
		// local file contents
		{
			File{
				Path: "/foo",
				Contents: Resource{
					Local: util.StrToPtr("file-1"),
				},
			},
			types.File{
				Node: types.Node{
					Path: "/foo",
				},
				FileEmbedded1: types.FileEmbedded1{
					Contents: types.Resource{
						Source: util.StrToPtr("data:,file%20contents%0A"),
					},
				},
			},
			[]translate.Translation{
				{
					From: path.New("yaml", "contents", "local"),
					To:   path.New("json", "contents", "source"),
				},
			},
			"",
			base.TranslateOptions{
				FilesDir: filesDir,
			},
		},
		// filesDir not specified
		{
			File{
				Path: "/foo",
				Contents: Resource{
					Local: util.StrToPtr("file-1"),
				},
			},
			types.File{
				Node: types.Node{
					Path: "/foo",
				},
			},
			[]translate.Translation{},
			"error at $.contents.local: " + ErrNoFilesDir.Error() + "\n",
			base.TranslateOptions{},
		},
		// attempted directory traversal
		{
			File{
				Path: "/foo",
				Contents: Resource{
					Local: util.StrToPtr("../file-1"),
				},
			},
			types.File{
				Node: types.Node{
					Path: "/foo",
				},
			},
			[]translate.Translation{},
			"error at $.contents.local: " + ErrFilesDirEscape.Error() + "\n",
			base.TranslateOptions{
				FilesDir: filesDir,
			},
		},
		// attempted inclusion of nonexistent file
		{
			File{
				Path: "/foo",
				Contents: Resource{
					Local: util.StrToPtr("file-2"),
				},
			},
			types.File{
				Node: types.Node{
					Path: "/foo",
				},
			},
			[]translate.Translation{},
			"error at $.contents.local: open " + filepath.Join(filesDir, "file-2") + ": no such file or directory\n",
			base.TranslateOptions{
				FilesDir: filesDir,
			},
		},
	}

	for i, test := range tests {
		actual, translations, report := translateFile(test.in, test.options)

		if !reflect.DeepEqual(actual, test.out) {
			t.Errorf("#%d: expected %+v got %+v", i, test.out, actual)
		}

		if report.String() != test.report {
			t.Errorf("#%d: expected report '%+v', got '%+v'", i, test.report, report.String())
		}

		if errT := verifyTranslations(translations, test.exceptions...); errT != nil {
			t.Errorf("#%d: bad translation: %v", i, *errT)
		}
	}
}

// TestTranslateDirectory tests translating the ct storage.directories.[i] entries to ignition storage.directories.[i] entires.
func TestTranslateDirectory(t *testing.T) {
	tests := []struct {
		in  Directory
		out types.Directory
	}{
		{
			Directory{},
			types.Directory{},
		},
		{
			// contains invalid (by the validator's definition) combinations of fields,
			// but the translator doesn't care and we can check they all get translated at once
			Directory{
				Path: "/foo",
				Group: NodeGroup{
					ID:   util.IntToPtr(1),
					Name: util.StrToPtr("foobar"),
				},
				User: NodeUser{
					ID:   util.IntToPtr(1),
					Name: util.StrToPtr("bazquux"),
				},
				Mode:      util.IntToPtr(420),
				Overwrite: util.BoolToPtr(true),
			},
			types.Directory{
				Node: types.Node{
					Path: "/foo",
					Group: types.NodeGroup{
						ID:   util.IntToPtr(1),
						Name: util.StrToPtr("foobar"),
					},
					User: types.NodeUser{
						ID:   util.IntToPtr(1),
						Name: util.StrToPtr("bazquux"),
					},
					Overwrite: util.BoolToPtr(true),
				},
				DirectoryEmbedded1: types.DirectoryEmbedded1{
					Mode: util.IntToPtr(420),
				},
			},
		},
	}

	for i, test := range tests {
		actual, _, report := translateDirectory(test.in, base.TranslateOptions{})
		if !reflect.DeepEqual(actual, test.out) {
			t.Errorf("#%d: expected %+v got %+v", i, test.out, actual)
		}
		if report.String() != "" {
			t.Errorf("#%d: got non-empty report: %v", i, report.String())
		}
	}
}

// TestTranslateLink tests translating the ct storage.links.[i] entries to ignition storage.links.[i] entires.
func TestTranslateLink(t *testing.T) {
	tests := []struct {
		in  Link
		out types.Link
	}{
		{
			Link{},
			types.Link{},
		},
		{
			// contains invalid (by the validator's definition) combinations of fields,
			// but the translator doesn't care and we can check they all get translated at once
			Link{
				Path: "/foo",
				Group: NodeGroup{
					ID:   util.IntToPtr(1),
					Name: util.StrToPtr("foobar"),
				},
				User: NodeUser{
					ID:   util.IntToPtr(1),
					Name: util.StrToPtr("bazquux"),
				},
				Overwrite: util.BoolToPtr(true),
				Target:    "/bar",
				Hard:      util.BoolToPtr(false),
			},
			types.Link{
				Node: types.Node{
					Path: "/foo",
					Group: types.NodeGroup{
						ID:   util.IntToPtr(1),
						Name: util.StrToPtr("foobar"),
					},
					User: types.NodeUser{
						ID:   util.IntToPtr(1),
						Name: util.StrToPtr("bazquux"),
					},
					Overwrite: util.BoolToPtr(true),
				},
				LinkEmbedded1: types.LinkEmbedded1{
					Target: "/bar",
					Hard:   util.BoolToPtr(false),
				},
			},
		},
	}

	for i, test := range tests {
		actual, _, report := translateLink(test.in, base.TranslateOptions{})
		if !reflect.DeepEqual(actual, test.out) {
			t.Errorf("#%d: expected %+v got %+v", i, test.out, actual)
		}
		if report.String() != "" {
			t.Errorf("#%d: got non-empty report: %v", i, report.String())
		}
	}
}

// TestTranslateFilesystem tests translating the ct storage.filesystems.[i] entries to ignition storage.filesystems.[i] entries.
func TestTranslateFilesystem(t *testing.T) {
	tests := []struct {
		in  Filesystem
		out types.Filesystem
	}{
		{
			Filesystem{},
			types.Filesystem{},
		},
		{
			// contains invalid (by the validator's definition) combinations of fields,
			// but the translator doesn't care and we can check they all get translated at once
			Filesystem{
				Device:         "/foo",
				Format:         util.StrToPtr("/bar"),
				Label:          util.StrToPtr("/baz"),
				MountOptions:   []string{"yes", "no", "maybe"},
				Options:        []string{"foo", "foo", "bar"},
				Path:           util.StrToPtr("/quux"),
				UUID:           util.StrToPtr("1234"),
				WipeFilesystem: util.BoolToPtr(true),
				WithMountUnit:  util.BoolToPtr(true),
			},
			types.Filesystem{
				Device:         "/foo",
				Format:         util.StrToPtr("/bar"),
				Label:          util.StrToPtr("/baz"),
				MountOptions:   []types.MountOption{"yes", "no", "maybe"},
				Options:        []types.FilesystemOption{"foo", "foo", "bar"},
				Path:           util.StrToPtr("/quux"),
				UUID:           util.StrToPtr("1234"),
				WipeFilesystem: util.BoolToPtr(true),
			},
		},
	}

	for i, test := range tests {
		actual, _, report := translateFilesystem(test.in, base.TranslateOptions{})
		if !reflect.DeepEqual(actual, test.out) {
			t.Errorf("#%d: expected %+v got %+v", i, test.out, actual)
		}
		if report.String() != "" {
			t.Errorf("#%d: got non-empty report: %v", i, report.String())
		}
	}
}

// TestTranslateIgnition tests translating the ct config.ignition to the ignition config.ignition section.
// It ensure that the version is set as well.
func TestTranslateIgnition(t *testing.T) {
	tests := []struct {
		in  Ignition
		out types.Ignition
	}{
		{
			Ignition{},
			types.Ignition{
				Version: "3.1.0-experimental",
			},
		},
		{
			Ignition{
				Config: IgnitionConfig{
					Merge: []Resource{
						{
							Inline: util.StrToPtr("xyzzy"),
						},
					},
					Replace: Resource{
						Inline: util.StrToPtr("xyzzy"),
					},
				},
			},
			types.Ignition{
				Version: "3.1.0-experimental",
				Config: types.IgnitionConfig{
					Merge: []types.Resource{
						{
							Source: util.StrToPtr("data:,xyzzy"),
						},
					},
					Replace: types.Resource{
						Source: util.StrToPtr("data:,xyzzy"),
					},
				},
			},
		},
		{
			Ignition{
				Proxy: Proxy{
					HTTPProxy: util.StrToPtr("https://example.com:8080"),
					NoProxy:   []string{"example.com"},
				},
			},
			types.Ignition{
				Version: "3.1.0-experimental",
				Proxy: types.Proxy{
					HTTPProxy: util.StrToPtr("https://example.com:8080"),
					NoProxy:   []types.NoProxyItem{types.NoProxyItem("example.com")},
				},
			},
		},
		{
			Ignition{
				Security: Security{
					TLS: TLS{
						CertificateAuthorities: []Resource{
							{
								Inline: util.StrToPtr("xyzzy"),
							},
						},
					},
				},
			},
			types.Ignition{
				Version: "3.1.0-experimental",
				Security: types.Security{
					TLS: types.TLS{
						CertificateAuthorities: []types.Resource{
							{
								Source: util.StrToPtr("data:,xyzzy"),
							},
						},
					},
				},
			},
		},
	}
	for i, test := range tests {
		actual, _, report := translateIgnition(test.in, base.TranslateOptions{})
		if !reflect.DeepEqual(actual, test.out) {
			t.Errorf("#%d: expected %+v got %+v", i, test.out, actual)
		}
		if report.String() != "" {
			t.Errorf("#%d: got non-empty report: %v", i, report.String())
		}
	}
}

// TestToIgn3_1 tests the config.ToIgn3_1 function ensuring it will generate a valid config even when empty. Not much else is
// tested since it uses the Ignition translation code which has it's own set of tests.
func TestToIgn3_1(t *testing.T) {
	tests := []struct {
		in  Config
		out types.Config
	}{
		{
			Config{},
			types.Config{
				Ignition: types.Ignition{
					Version: "3.1.0-experimental",
				},
			},
		},
	}
	for i, test := range tests {
		actual, _, report := test.in.ToIgn3_1(base.TranslateOptions{})
		if report.String() != "" {
			t.Errorf("#%d: got non-empty report: %v", i, report.String())
		}
		if !reflect.DeepEqual(actual, test.out) {
			t.Errorf("#%d: expected %+v got %+v", i, test.out, actual)
		}
	}
}
