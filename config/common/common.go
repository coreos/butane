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
	"bytes"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/coreos/fcct/translate"

	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
	"github.com/coreos/vcontext/tree"
	vyaml "github.com/coreos/vcontext/yaml"
	"gopkg.in/yaml.v3"
)

var (
	snakeRe = regexp.MustCompile("([A-Z])")
)

type TranslateOptions struct {
	Pretty bool
	Strict bool
}

type Common struct {
	Version string `yaml:"version"`
	Variant string `yaml:"variant"`
}

// Misc helpers

// Unmarshal unmarshals the data to "to" and also returns a context tree for the source. If strict
// is set it errors out on unused keys.
func Unmarshal(data []byte, to interface{}, strict bool) (tree.Node, error) {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(strict)
	if err := dec.Decode(to); err != nil {
		return nil, err
	}
	return vyaml.UnmarshalToContext(data)
}

// Marshal is a wrapper for marshaling to json with or without pretty-printing the output
func Marshal(from interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(from, "", "  ")
	}
	return json.Marshal(from)
}

// snakePath converts a path.ContextPath with camelCase elements and returns the
// same path but with snake_case elements instead
func snakePath(p path.ContextPath) path.ContextPath {
	ret := path.New(p.Tag)
	for _, part := range p.Path {
		if str, ok := part.(string); ok {
			ret = ret.Append(snake(str))
		} else {
			ret = ret.Append(part)
		}
	}
	return ret
}

// snake converts from camelCase (not CamelCase) to snake_case
func snake(in string) string {
	return strings.ToLower(snakeRe.ReplaceAllString(in, "_$1"))
}

// TranslateReportPaths takes a report from a camelCase json document and a set of translations rules,
// applies those rules and converts all camelCase to snake_case.
func TranslateReportPaths(r *report.Report, ts translate.TranslationSet) {
	for i, ent := range r.Entries {
		context := ent.Context
		if t, ok := ts.Set[context.String()]; ok {
			context = t.From
		}
		context = snakePath(context)
		r.Entries[i].Context = context
	}
}
