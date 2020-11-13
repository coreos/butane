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

package util

import (
	"bytes"
	"reflect"
	"regexp"
	"strings"

	"github.com/coreos/fcct/config/common"
	"github.com/coreos/fcct/translate"

	"github.com/clarketm/json"
	ignvalidate "github.com/coreos/ignition/v2/config/validate"
	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
	"github.com/coreos/vcontext/tree"
	"github.com/coreos/vcontext/validate"
	vyaml "github.com/coreos/vcontext/yaml"
	"gopkg.in/yaml.v3"
)

var (
	snakeRe = regexp.MustCompile("([A-Z])")
)

// Misc helpers

// Translate translates cfg to the corresponding Ignition config version
// using the named translation method on cfg, and returns the marshaled
// Ignition config.  It returns a report of any errors or warnings in the
// source and resultant config.  If the report has fatal errors or it
// encounters other problems translating, an error is returned.
func Translate(cfg interface{}, translateMethod string, options common.TranslateOptions) (interface{}, report.Report, error) {
	// Get method, and zero return value for error returns.
	method := reflect.ValueOf(cfg).MethodByName(translateMethod)
	zeroValue := reflect.Zero(method.Type().Out(0)).Interface()

	// Validate the input.
	r := validate.Validate(cfg, "yaml")
	if r.IsFatal() {
		return zeroValue, r, common.ErrInvalidSourceConfig
	}

	// Perform the translation.
	translateRet := method.Call([]reflect.Value{reflect.ValueOf(options)})
	final := translateRet[0].Interface()
	translations := translateRet[1].Interface().(translate.TranslationSet)
	translateReport := translateRet[2].Interface().(report.Report)
	r.Merge(translateReport)
	if r.IsFatal() {
		return zeroValue, r, common.ErrInvalidSourceConfig
	}

	// Check for invalid duplicated keys.
	dupsReport := validate.ValidateCustom(final, "json", ignvalidate.ValidateDups)
	translateReportPaths(&dupsReport, translations)
	r.Merge(dupsReport)

	// Validate JSON semantics.
	jsonReport := validate.Validate(final, "json")
	translateReportPaths(&jsonReport, translations)
	r.Merge(jsonReport)

	if r.IsFatal() {
		return zeroValue, r, common.ErrInvalidGeneratedConfig
	}
	return final, r, nil
}

// TranslateBytes unmarshals the FCC specified in input into the struct
// pointed to by container, translates it to the corresponding Ignition
// config version using the named translation method, and returns the
// marshaled Ignition config.  It returns a report of any errors or warnings
// in the source and resultant config.  If the report has fatal errors or it
// encounters other problems translating, an error is returned.
func TranslateBytes(input []byte, container interface{}, translateMethod string, options common.TranslateBytesOptions) ([]byte, report.Report, error) {
	cfg := container

	// Unmarshal the YAML.
	contextTree, err := unmarshal(input, cfg, options.Strict)
	if err != nil {
		return nil, report.Report{}, err
	}

	// Check for unused keys.
	unusedKeyCheck := func(v reflect.Value, c path.ContextPath) report.Report {
		return ignvalidate.ValidateUnusedKeys(v, c, contextTree)
	}
	r := validate.ValidateCustom(cfg, "yaml", unusedKeyCheck)
	r.Correlate(contextTree)
	if r.IsFatal() {
		return nil, r, common.ErrInvalidSourceConfig
	}

	// Perform the translation.
	translateRet := reflect.ValueOf(cfg).MethodByName(translateMethod).Call([]reflect.Value{reflect.ValueOf(options.TranslateOptions)})
	final := translateRet[0].Interface()
	translateReport := translateRet[1].Interface().(report.Report)
	errVal := translateRet[2]
	translateReport.Correlate(contextTree)
	r.Merge(translateReport)
	if !errVal.IsNil() {
		return nil, r, errVal.Interface().(error)
	}
	if r.IsFatal() {
		return nil, r, common.ErrInvalidSourceConfig
	}

	// Marshal the JSON.
	outbytes, err := marshal(final, options.Pretty)
	return outbytes, r, err
}

// unmarshal unmarshals the data to "to" and also returns a context tree for the source. If strict
// is set it errors out on unused keys.
func unmarshal(data []byte, to interface{}, strict bool) (tree.Node, error) {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(strict)
	if err := dec.Decode(to); err != nil {
		return nil, err
	}
	return vyaml.UnmarshalToContext(data)
}

// marshal is a wrapper for marshaling to json with or without pretty-printing the output
func marshal(from interface{}, pretty bool) ([]byte, error) {
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

// translateReportPaths takes a report from a camelCase json document and a set of translations rules,
// applies those rules and converts all camelCase to snake_case.
func translateReportPaths(r *report.Report, ts translate.TranslationSet) {
	for i, ent := range r.Entries {
		context := ent.Context
		if t, ok := ts.Set[context.String()]; ok {
			context = t.From
		}
		context = snakePath(context)
		r.Entries[i].Context = context
	}
}
