// Copyright 2019 Red Hat, Inc.
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
// limitations under the License.

package translate

import (
	"reflect"
	"strings"
)

// fieldName returns the name uses when (un)marshalling a field. t should be a reflect.Value of a struct,
// index is the field index, and tag is the struct tag used when (un)marshalling (e.g. "json" or "yaml")
func fieldName(t reflect.Value, index int, tag string) string {
	f := t.Type().Field(index)
	if tag == "" {
		return f.Name
	}
	return strings.Split(f.Tag.Get(tag), ",")[0]
}
