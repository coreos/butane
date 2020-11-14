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
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/coreos/fcct/config/common"
	"github.com/coreos/fcct/translate"

	"github.com/clarketm/json"
	"github.com/coreos/ignition/v2/config/util"
	"github.com/coreos/ignition/v2/config/v3_1/types"
	"github.com/coreos/vcontext/path"
)

// Most of this is covered by the Ignition translator generic tests, so just test the custom bits

func format(from interface{}) string {
	buf, err := json.MarshalIndent(from, "", "  ")
	if err != nil {
		return fmt.Sprintf("<Error marshaling to JSON: %v>", err)
	}
	return string(buf)
}

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
	zzz := "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"
	zzz_gz := "data:;base64,H4sIAAAAAAAC/6oajAAQAAD//5tA8d+VAAAA"
	random := "\xc0\x9cl\x01\x89i\xa5\xbfW\xe4\x1b\xf4J_\xb79P\xa3#\xa7"
	random_b64 := "data:;base64,wJxsAYlppb9X5Bv0Sl+3OVCjI6c="

	filesDir, err := ioutil.TempDir("", "translate-test-")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.RemoveAll(filesDir)
	fileContents := map[string]string{
		"file-1": "file contents\n",
		"file-2": zzz,
		"file-3": random,
	}
	for name, contents := range fileContents {
		err := ioutil.WriteFile(filepath.Join(filesDir, name), []byte(contents), 0644)
		if err != nil {
			t.Error(err)
			return
		}
	}

	tests := []struct {
		in         File
		out        types.File
		exceptions []translate.Translation
		report     string
		options    common.TranslateOptions
	}{
		{
			File{},
			types.File{},
			nil,
			"",
			common.TranslateOptions{},
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
			common.TranslateOptions{
				FilesDir: filesDir,
			},
		},
		// inline file contents
		{
			File{
				Path: "/foo",
				Contents: Resource{
					// String is too short for auto gzip compression
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
			common.TranslateOptions{},
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
			common.TranslateOptions{
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
			"error at $.contents.local: " + common.ErrNoFilesDir.Error() + "\n",
			common.TranslateOptions{},
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
			"error at $.contents.local: " + common.ErrFilesDirEscape.Error() + "\n",
			common.TranslateOptions{
				FilesDir: filesDir,
			},
		},
		// attempted inclusion of nonexistent file
		{
			File{
				Path: "/foo",
				Contents: Resource{
					Local: util.StrToPtr("file-missing"),
				},
			},
			types.File{
				Node: types.Node{
					Path: "/foo",
				},
			},
			[]translate.Translation{},
			"error at $.contents.local: open " + filepath.Join(filesDir, "file-missing") + ": no such file or directory\n",
			common.TranslateOptions{
				FilesDir: filesDir,
			},
		},
		// inline and local automatic file encoding
		{
			File{
				Path: "/foo",
				Contents: Resource{
					// gzip
					Inline: util.StrToPtr(zzz),
				},
				Append: []Resource{
					{
						// gzip
						Local: util.StrToPtr("file-2"),
					},
					{
						// base64
						Inline: util.StrToPtr(random),
					},
					{
						// base64
						Local: util.StrToPtr("file-3"),
					},
					{
						// URL-escaped
						Inline:      util.StrToPtr(zzz),
						Compression: util.StrToPtr("invalid"),
					},
					{
						// URL-escaped
						Local:       util.StrToPtr("file-2"),
						Compression: util.StrToPtr("invalid"),
					},
				},
			},
			types.File{
				Node: types.Node{
					Path: "/foo",
				},
				FileEmbedded1: types.FileEmbedded1{
					Contents: types.Resource{
						Source:      util.StrToPtr(zzz_gz),
						Compression: util.StrToPtr("gzip"),
					},
					Append: []types.Resource{
						{
							Source:      util.StrToPtr(zzz_gz),
							Compression: util.StrToPtr("gzip"),
						},
						{
							Source: util.StrToPtr(random_b64),
						},
						{
							Source: util.StrToPtr(random_b64),
						},
						{
							Source:      util.StrToPtr("data:," + zzz),
							Compression: util.StrToPtr("invalid"),
						},
						{
							Source:      util.StrToPtr("data:," + zzz),
							Compression: util.StrToPtr("invalid"),
						},
					},
				},
			},
			[]translate.Translation{
				{
					From: path.New("yaml", "contents", "inline"),
					To:   path.New("json", "contents", "source"),
				},
				{
					From: path.New("yaml", "contents", "inline"),
					To:   path.New("json", "contents", "compression"),
				},
				{
					From: path.New("yaml", "append", 0, "local"),
					To:   path.New("json", "append", 0, "source"),
				},
				{
					From: path.New("yaml", "append", 0, "local"),
					To:   path.New("json", "append", 0, "compression"),
				},
				{
					From: path.New("yaml", "append", 1, "inline"),
					To:   path.New("json", "append", 1, "source"),
				},
				{
					From: path.New("yaml", "append", 2, "local"),
					To:   path.New("json", "append", 2, "source"),
				},
				{
					From: path.New("yaml", "append", 3, "inline"),
					To:   path.New("json", "append", 3, "source"),
				},
				{
					From: path.New("yaml", "append", 4, "local"),
					To:   path.New("json", "append", 4, "source"),
				},
			},
			"",
			common.TranslateOptions{
				FilesDir: filesDir,
			},
		},
		// Test disable automatic gzip compression
		{
			File{
				Path: "/foo",
				Contents: Resource{
					Inline: util.StrToPtr(zzz),
				},
			},
			types.File{
				Node: types.Node{
					Path: "/foo",
				},
				FileEmbedded1: types.FileEmbedded1{
					Contents: types.Resource{
						Source: util.StrToPtr("data:," + zzz),
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
			common.TranslateOptions{
				NoResourceAutoCompression: true,
			},
		},
	}

	for i, test := range tests {
		actual, translations, report := translateFile(test.in, test.options)

		if !reflect.DeepEqual(actual, test.out) {
			t.Errorf("#%d: expected %v, got %v", i, format(test.out), format(actual))
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
		actual, _, report := translateDirectory(test.in, common.TranslateOptions{})
		if !reflect.DeepEqual(actual, test.out) {
			t.Errorf("#%d: expected %v, got %v", i, format(test.out), format(actual))
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
		actual, _, report := translateLink(test.in, common.TranslateOptions{})
		if !reflect.DeepEqual(actual, test.out) {
			t.Errorf("#%d: expected %v, got %v", i, format(test.out), format(actual))
		}
		if report.String() != "" {
			t.Errorf("#%d: got non-empty report: %v", i, report.String())
		}
	}
}

// TestTranslateFilesystem tests translating the fcct storage.filesystems.[i] entries to ignition storage.filesystems.[i] entries.
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
		// Filesystem doesn't have a custom translator, so embed in a
		// complete config
		in := Config{
			Storage: Storage{
				Filesystems: []Filesystem{test.in},
			},
		}
		expected := []types.Filesystem{test.out}
		actual, _, report := in.ToIgn3_1Unvalidated(common.TranslateOptions{})
		if !reflect.DeepEqual(actual.Storage.Filesystems, expected) {
			t.Errorf("#%d: expected %v, got %v", i, format(expected), format(actual.Storage.Filesystems))
		}
		if report.String() != "" {
			t.Errorf("#%d: got non-empty report: %v", i, report.String())
		}
	}
}

// TestTranslateTree tests translating the FCC storage.trees.[i] entries to ignition storage.files.[i] entries.
func TestTranslateTree(t *testing.T) {
	tests := []struct {
		options    *common.TranslateOptions // defaulted if not specified
		dirDirs    map[string]os.FileMode   // relative path -> mode
		dirFiles   map[string]os.FileMode   // relative path -> mode
		dirLinks   map[string]string        // relative path -> target
		dirSockets []string                 // relative path
		inTrees    []Tree
		inFiles    []File
		inDirs     []Directory
		inLinks    []Link
		outFiles   []types.File
		outLinks   []types.Link
		report     string
	}{
		// smoke test
		{},
		// basic functionality
		{
			dirFiles: map[string]os.FileMode{
				"tree/executable":            0700,
				"tree/file":                  0600,
				"tree/overridden":            0644,
				"tree/overridden-executable": 0700,
				"tree/subdir/file":           0644,
				// compressed contents
				"tree/subdir/subdir/subdir/subdir/subdir/subdir/subdir/subdir/subdir/file": 0644,
				"tree2/file": 0600,
			},
			dirLinks: map[string]string{
				"tree/subdir/bad-link":        "../nonexistent",
				"tree/subdir/link":            "../file",
				"tree/subdir/overridden-link": "../file",
			},
			inTrees: []Tree{
				{
					Local: "tree",
				},
				{
					Local: "tree2",
					Path:  util.StrToPtr("/etc"),
				},
			},
			inFiles: []File{
				{
					Path: "/overridden",
					Mode: util.IntToPtr(0600),
					User: NodeUser{
						Name: util.StrToPtr("bovik"),
					},
				},
				{
					Path: "/overridden-executable",
					Mode: util.IntToPtr(0600),
					User: NodeUser{
						Name: util.StrToPtr("bovik"),
					},
				},
			},
			inLinks: []Link{
				{
					Path: "/subdir/overridden-link",
					User: NodeUser{
						Name: util.StrToPtr("bovik"),
					},
				},
			},
			outFiles: []types.File{
				{
					Node: types.Node{
						Path: "/overridden",
						User: types.NodeUser{
							Name: util.StrToPtr("bovik"),
						},
					},
					FileEmbedded1: types.FileEmbedded1{
						Contents: types.Resource{
							Source: util.StrToPtr("data:,tree%2Foverridden"),
						},
						Mode: util.IntToPtr(0600),
					},
				},
				{
					Node: types.Node{
						Path: "/overridden-executable",
						User: types.NodeUser{
							Name: util.StrToPtr("bovik"),
						},
					},
					FileEmbedded1: types.FileEmbedded1{
						Contents: types.Resource{
							Source: util.StrToPtr("data:,tree%2Foverridden-executable"),
						},
						Mode: util.IntToPtr(0600),
					},
				},
				{
					Node: types.Node{
						Path: "/executable",
					},
					FileEmbedded1: types.FileEmbedded1{
						Contents: types.Resource{
							Source: util.StrToPtr("data:,tree%2Fexecutable"),
						},
						Mode: util.IntToPtr(0755),
					},
				},
				{
					Node: types.Node{
						Path: "/file",
					},
					FileEmbedded1: types.FileEmbedded1{
						Contents: types.Resource{
							Source: util.StrToPtr("data:,tree%2Ffile"),
						},
						Mode: util.IntToPtr(0644),
					},
				},
				{
					Node: types.Node{
						Path: "/subdir/file",
					},
					FileEmbedded1: types.FileEmbedded1{
						Contents: types.Resource{
							Source: util.StrToPtr("data:,tree%2Fsubdir%2Ffile"),
						},
						Mode: util.IntToPtr(0644),
					},
				},
				{
					Node: types.Node{
						Path: "/subdir/subdir/subdir/subdir/subdir/subdir/subdir/subdir/subdir/file",
					},
					FileEmbedded1: types.FileEmbedded1{
						Contents: types.Resource{
							Source:      util.StrToPtr("data:;base64,H4sIAAAAAAAC/yopSk3VLy5NSsksIptKy8xJBQQAAP//gkRzjkgAAAA="),
							Compression: util.StrToPtr("gzip"),
						},
						Mode: util.IntToPtr(0644),
					},
				},
				{
					Node: types.Node{
						Path: "/etc/file",
					},
					FileEmbedded1: types.FileEmbedded1{
						Contents: types.Resource{
							Source: util.StrToPtr("data:,tree2%2Ffile"),
						},
						Mode: util.IntToPtr(0644),
					},
				},
			},
			outLinks: []types.Link{
				{
					Node: types.Node{
						Path: "/subdir/overridden-link",
						User: types.NodeUser{
							Name: util.StrToPtr("bovik"),
						},
					},
					LinkEmbedded1: types.LinkEmbedded1{
						Target: "../file",
					},
				},
				{
					Node: types.Node{
						Path: "/subdir/bad-link",
					},
					LinkEmbedded1: types.LinkEmbedded1{
						Target: "../nonexistent",
					},
				},
				{
					Node: types.Node{
						Path: "/subdir/link",
					},
					LinkEmbedded1: types.LinkEmbedded1{
						Target: "../file",
					},
				},
			},
		},
		// collisions
		{
			dirFiles: map[string]os.FileMode{
				"tree0/file":         0600,
				"tree1/directory":    0600,
				"tree2/link":         0600,
				"tree3/file-partial": 0600, // should be okay
				"tree4/link-partial": 0600,
				"tree5/tree-file":    0600, // set up for tree/tree collision
				"tree6/tree-file":    0600,
				"tree15/tree-link":   0600,
			},
			dirLinks: map[string]string{
				"tree7/file":          "file",
				"tree8/directory":     "file",
				"tree9/link":          "file",
				"tree10/file-partial": "file",
				"tree11/link-partial": "file", // should be okay
				"tree12/tree-file":    "file",
				"tree13/tree-link":    "file", // set up for tree/tree collision
				"tree14/tree-link":    "file",
			},
			inTrees: []Tree{
				{
					Local: "tree0",
				},
				{
					Local: "tree1",
				},
				{
					Local: "tree2",
				},
				{
					Local: "tree3",
				},
				{
					Local: "tree4",
				},
				{
					Local: "tree5",
				},
				{
					Local: "tree6",
				},
				{
					Local: "tree7",
				},
				{
					Local: "tree8",
				},
				{
					Local: "tree9",
				},
				{
					Local: "tree10",
				},
				{
					Local: "tree11",
				},
				{
					Local: "tree12",
				},
				{
					Local: "tree13",
				},
				{
					Local: "tree14",
				},
				{
					Local: "tree15",
				},
			},
			inFiles: []File{
				{
					Path: "/file",
					Contents: Resource{
						Source: util.StrToPtr("data:,foo"),
					},
				},
				{
					Path: "/file-partial",
				},
			},
			inDirs: []Directory{
				{
					Path: "/directory",
				},
			},
			inLinks: []Link{
				{
					Path:   "/link",
					Target: "file",
				},
				{
					Path: "/link-partial",
				},
			},
			report: "error at $.storage.trees.0: " + common.ErrNodeExists.Error() + "\n" +
				"error at $.storage.trees.1: " + common.ErrNodeExists.Error() + "\n" +
				"error at $.storage.trees.2: " + common.ErrNodeExists.Error() + "\n" +
				"error at $.storage.trees.4: " + common.ErrNodeExists.Error() + "\n" +
				"error at $.storage.trees.6: " + common.ErrNodeExists.Error() + "\n" +
				"error at $.storage.trees.7: " + common.ErrNodeExists.Error() + "\n" +
				"error at $.storage.trees.8: " + common.ErrNodeExists.Error() + "\n" +
				"error at $.storage.trees.9: " + common.ErrNodeExists.Error() + "\n" +
				"error at $.storage.trees.10: " + common.ErrNodeExists.Error() + "\n" +
				"error at $.storage.trees.12: " + common.ErrNodeExists.Error() + "\n" +
				"error at $.storage.trees.14: " + common.ErrNodeExists.Error() + "\n" +
				"error at $.storage.trees.15: " + common.ErrNodeExists.Error() + "\n",
		},
		// files-dir escape
		{
			inTrees: []Tree{
				{
					Local: "../escape",
				},
			},
			report: "error at $.storage.trees.0: " + common.ErrFilesDirEscape.Error() + "\n",
		},
		// no files-dir
		{
			options: &common.TranslateOptions{},
			inTrees: []Tree{
				{
					Local: "tree",
				},
			},
			report: "error at $.storage.trees.0: " + common.ErrNoFilesDir.Error() + "\n",
		},
		// non-file/dir/symlink in directory tree
		{
			dirSockets: []string{
				"tree/socket",
			},
			inTrees: []Tree{
				{
					Local: "tree",
				},
			},
			report: "error at $.storage.trees.0: " + common.ErrFileType.Error() + "\n",
		},
		// unreadable file
		{
			dirDirs: map[string]os.FileMode{
				"tree/subdir": 0000,
				"tree2":       0000,
			},
			dirFiles: map[string]os.FileMode{
				"tree/file": 0000,
			},
			inTrees: []Tree{
				{
					Local: "tree",
				},
				{
					Local: "tree2",
				},
			},
			report: "error at $.storage.trees.0: open %FilesDir%/tree/file: permission denied\n" +
				"error at $.storage.trees.0: open %FilesDir%/tree/subdir: permission denied\n" +
				"error at $.storage.trees.1: open %FilesDir%/tree2: permission denied\n",
		},
		// local is not a directory
		{
			dirFiles: map[string]os.FileMode{
				"tree": 0600,
			},
			inTrees: []Tree{
				{
					Local: "tree",
				},
				{
					Local: "nonexistent",
				},
			},
			report: "error at $.storage.trees.0: " + common.ErrTreeNotDirectory.Error() + "\n" +
				"error at $.storage.trees.1: stat %FilesDir%/nonexistent: no such file or directory\n",
		},
	}

	for i, test := range tests {
		filesDir, err := ioutil.TempDir("", "translate-test-")
		if err != nil {
			t.Error(err)
			return
		}
		defer os.RemoveAll(filesDir)
		for path, mode := range test.dirDirs {
			absPath := filepath.Join(filesDir, path)
			if err := os.MkdirAll(absPath, 0755); err != nil {
				t.Error(err)
				return
			}
			if err := os.Chmod(absPath, mode); err != nil {
				t.Error(err)
				return
			}
		}
		for path, mode := range test.dirFiles {
			absPath := filepath.Join(filesDir, path)
			if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
				t.Error(err)
				return
			}
			if err := ioutil.WriteFile(absPath, []byte(path), mode); err != nil {
				t.Error(err)
				return
			}
		}
		for path, target := range test.dirLinks {
			absPath := filepath.Join(filesDir, path)
			if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
				t.Error(err)
				return
			}
			if err := os.Symlink(target, absPath); err != nil {
				t.Error(err)
				return
			}
		}
		for _, path := range test.dirSockets {
			absPath := filepath.Join(filesDir, path)
			if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
				t.Error(err)
				return
			}
			listener, err := net.ListenUnix("unix", &net.UnixAddr{
				Name: absPath,
				Net:  "unix",
			})
			if err != nil {
				t.Error(err)
				return
			}
			defer listener.Close()
		}

		config := Config{
			Storage: Storage{
				Files:       test.inFiles,
				Directories: test.inDirs,
				Links:       test.inLinks,
				Trees:       test.inTrees,
			},
		}
		options := common.TranslateOptions{
			FilesDir: filesDir,
		}
		if test.options != nil {
			options = *test.options
		}
		actual, _, report := config.ToIgn3_1Unvalidated(options)

		expectedReport := strings.ReplaceAll(test.report, "%FilesDir%", filesDir)
		if report.String() != expectedReport {
			t.Errorf("#%d: expected report '%+v', got '%+v'", i, expectedReport, report.String())
		}
		if expectedReport != "" {
			continue
		}

		if !reflect.DeepEqual(actual.Storage.Files, test.outFiles) {
			t.Errorf("#%d: expected files %v, got %v", i, format(test.outFiles), format(actual.Storage.Files))
		}

		if len(actual.Storage.Directories) != 0 {
			t.Errorf("#%d: expected empty directories, got %v", i, format(actual.Storage.Directories))
		}

		if !reflect.DeepEqual(actual.Storage.Links, test.outLinks) {
			t.Errorf("#%d: expected links %v, got %v", i, format(test.outLinks), format(actual.Storage.Links))
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
				Version: "3.1.0",
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
				Version: "3.1.0",
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
				Version: "3.1.0",
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
				Version: "3.1.0",
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
		actual, _, report := translateIgnition(test.in, common.TranslateOptions{})
		if !reflect.DeepEqual(actual, test.out) {
			t.Errorf("#%d: expected %v, got %v", i, format(test.out), format(actual))
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
					Version: "3.1.0",
				},
			},
		},
	}
	for i, test := range tests {
		actual, _, report := test.in.ToIgn3_1Unvalidated(common.TranslateOptions{})
		if report.String() != "" {
			t.Errorf("#%d: got non-empty report: %v", i, report.String())
		}
		if !reflect.DeepEqual(actual, test.out) {
			t.Errorf("#%d: expected %v, got %v", i, format(test.out), format(actual))
		}
	}
}
