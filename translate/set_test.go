// Copyright 2022 Red Hat, Inc.
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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslationSetMap(t *testing.T) {
	create := func() TranslationSet {
		return mkTrans(
			fp(), fp(),
			fp("a"), fp("A"),
			fp("a", 0), fp("A", 0),
			fp("a", 0, "b"), fp("A", 0, "B"),
			fp("a", 0, "b", "c"), fp("A", 0, "B"),
			fp("a", 0, "b", "d"), fp("A", 0, "B", 0),
			fp("a", 0, "b", "e"), fp("A", 0, "B", 0, "C"),
			fp("a", 0, "b", "f"), fp("A", 0, "B", 0, "D"),
			fp("clobbered"), fp("A", 0, "B", 0, "G"),
			fp("a", 0, "b", "g"), fp("A", 0, "B", 1),
			fp("a", 0, "b", "h"), fp("A", 0, "B", 1, "E"),
			fp("a", 0, "b", "i"), fp("A", 0, "B", 1, "F"),
		)
	}
	ts := create()
	result := ts.Map(mkTrans(
		fp("A", 0, "B", 0, "C"), fp("A", 0, "B", 0, "G"),
		fp("A", 0, "B", 0, "D"), fp("A", 0, "H"),
		fp("missing"), fp("B"),
	))
	assert.Equal(t, create(), ts, "original was changed")
	assert.Equal(t, mkTrans(
		fp(), fp(),
		fp("a"), fp("A"),
		fp("a", 0), fp("A", 0),
		fp("a", 0, "b"), fp("A", 0, "B"),
		fp("a", 0, "b", "c"), fp("A", 0, "B"),
		fp("a", 0, "b", "d"), fp("A", 0, "B", 0),
		fp("a", 0, "b", "e"), fp("A", 0, "B", 0, "G"),
		fp("a", 0, "b", "f"), fp("A", 0, "H"),
		fp("a", 0, "b", "g"), fp("A", 0, "B", 1),
		fp("a", 0, "b", "h"), fp("A", 0, "B", 1, "E"),
		fp("a", 0, "b", "i"), fp("A", 0, "B", 1, "F"),
	), result, "bad mapping")
}
