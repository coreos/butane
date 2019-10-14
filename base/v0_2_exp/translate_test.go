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
	"reflect"
	"testing"

	"github.com/coreos/fcct/translate"

	"github.com/coreos/ignition/v2/config/util"
	"github.com/coreos/ignition/v2/config/v3_1_experimental/types"
	"github.com/coreos/vcontext/path"
)

// Most of this is covered by the Ignition translator generic tests, so just test the custom bits

// verifyTranslations ensures all the translations are identity, unless they match a listed one
// it returns the offending translation if there is one
func verifyTranslations(set translate.TranslationSet, exceptions ...translate.Translation) *translate.Translation {
	exceptionSet := translate.TranslationSet{
		FromTag: set.FromTag,
		ToTag:   set.ToTag,
		Set:     map[string]translate.Translation{},
	}
	for _, ex := range exceptions {
		exceptionSet.AddTranslation(ex.From, ex.To)
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

// TestTranslateFile tests translating the ct storage.files.[i] entries to ignition storage.files.[i] entires.
func TestTranslateFile(t *testing.T) {
	tests := []struct {
		in         File
		out        types.File
		exceptions []translate.Translation
	}{
		{
			File{},
			types.File{},
			nil,
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
				Append: []FileContents{
					{
						Source:      util.StrToPtr("http://example/com"),
						Compression: util.StrToPtr("gzip"),
						Verification: Verification{
							Hash: util.StrToPtr("this isn't validated"),
						},
					},
					{
						Inline:      util.StrToPtr("hello"),
						Compression: util.StrToPtr("gzip"),
						Verification: Verification{
							Hash: util.StrToPtr("this isn't validated"),
						},
					},
				},
				Overwrite: util.BoolToPtr(true),
				Contents: FileContents{
					Source:      util.StrToPtr("http://example/com"),
					Compression: util.StrToPtr("gzip"),
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
					Append: []types.FileContents{
						{
							Source:      util.StrToPtr("http://example/com"),
							Compression: util.StrToPtr("gzip"),
							Verification: types.Verification{
								Hash: util.StrToPtr("this isn't validated"),
							},
						},
						{
							Source:      util.StrToPtr("data:,hello"),
							Compression: util.StrToPtr("gzip"),
							Verification: types.Verification{
								Hash: util.StrToPtr("this isn't validated"),
							},
						},
					},
					Contents: types.FileContents{
						Source:      util.StrToPtr("http://example/com"),
						Compression: util.StrToPtr("gzip"),
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
			},
		},
	}

	for i, test := range tests {
		actual, translations := translateFile(test.in)

		if !reflect.DeepEqual(actual, test.out) {
			t.Errorf("#%d: expected %+v got %+v", i, test.out, actual)
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
		actual, _ := translateDirectory(test.in)
		if !reflect.DeepEqual(actual, test.out) {
			t.Errorf("#%d: expected %+v got %+v", i, test.out, actual)
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
		actual, _ := translateLink(test.in)
		if !reflect.DeepEqual(actual, test.out) {
			t.Errorf("#%d: expected %+v got %+v", i, test.out, actual)
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
	}
	for i, test := range tests {
		actual, _ := translateIgnition(test.in)
		if !reflect.DeepEqual(actual, test.out) {
			t.Errorf("#%d: expected %+v got %+v", i, test.out, actual)
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
		actual, _, err := test.in.ToIgn3_1()
		if err != nil {
			t.Errorf("#%d: got error: %v", i, err)
		}

		if !reflect.DeepEqual(actual, test.out) {
			t.Errorf("#%d: expected %+v got %+v", i, test.out, actual)
		}
	}
}
