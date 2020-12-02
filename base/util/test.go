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

package util

import (
	"reflect"

	"github.com/coreos/fcct/translate"
)

// helper functions for writing tests

// VerifyTranslations ensures all the translations are identity, unless they
// match a listed one, and verifies that all the listed ones exist.
// it returns the offending translation if there is one
func VerifyTranslations(set translate.TranslationSet, exceptions ...translate.Translation) *translate.Translation {
	exceptionSet := translate.NewTranslationSet(set.FromTag, set.ToTag)
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
