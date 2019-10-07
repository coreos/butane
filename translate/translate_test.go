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
	"testing"

	"github.com/coreos/fcct/translate/tests/pkga"
	"github.com/coreos/fcct/translate/tests/pkgb"

	"github.com/coreos/vcontext/path"
	"github.com/stretchr/testify/assert"
)

// Note: we need different input and output types which unfortunately means a lot of tests

// mkTrans makes a TranslationSet with no tag in the paths consuming pairs of args. i.e:
// mkTrans(from1, to1, from2, to2) -> a set wiht from1->to1, from2->to2
// This is just a shorthand for making writing tests easier
func mkTrans(paths ...path.ContextPath) TranslationSet {
	ret := TranslationSet{Set: map[string]Translation{}}
	if len(paths)%2 == 1 {
		panic("Odd number of args to mkTrans")
	}
	for i := 0; i < len(paths); i += 2 {
		ret.AddTranslation(paths[i], paths[i+1])
	}
	return ret
}

// fp means "fastpath"; super shorthand, we'll use it a lot
func fp(parts ...interface{}) path.ContextPath {
	return path.New("", parts...)
}

func TestTranslateTrivial(t *testing.T) {
	in := pkga.Trivial{
		A: "asdf",
		B: 5,
		C: true,
	}

	expected := pkgb.Trivial{
		A: "asdf",
		B: 5,
		C: true,
	}
	exTrans := mkTrans(
		fp("A"), fp("A"),
		fp("B"), fp("B"),
		fp("C"), fp("C"),
	)

	got := pkgb.Trivial{}

	trans := NewTranslator("", "")

	ts := trans.Translate(&in, &got)
	assert.Equal(t, got, expected, "bad translation")
	assert.Equal(t, ts, exTrans, "bad translation")
}

func TestTranslateNested(t *testing.T) {
	in := pkga.Nested{
		D: "foobar",
		Trivial: pkga.Trivial{
			A: "asdf",
			B: 5,
			C: true,
		},
	}

	expected := pkgb.Nested{
		D: "foobar",
		Trivial: pkgb.Trivial{
			A: "asdf",
			B: 5,
			C: true,
		},
	}
	exTrans := mkTrans(
		fp("A"), fp("A"),
		fp("B"), fp("B"),
		fp("C"), fp("C"),
		fp("D"), fp("D"),
	)

	got := pkgb.Nested{}

	trans := NewTranslator("", "")

	ts := trans.Translate(&in, &got)
	assert.Equal(t, got, expected, "bad translation")
	assert.Equal(t, ts, exTrans, "bad translation")
}

func TestTranslateTrivialReordered(t *testing.T) {
	in := pkga.TrivialReordered{
		A: "asdf",
		B: 5,
		C: true,
	}

	expected := pkgb.TrivialReordered{
		A: "asdf",
		B: 5,
		C: true,
	}
	exTrans := mkTrans(
		fp("A"), fp("A"),
		fp("B"), fp("B"),
		fp("C"), fp("C"),
	)

	got := pkgb.TrivialReordered{}

	trans := NewTranslator("", "")

	ts := trans.Translate(&in, &got)
	assert.Equal(t, got, expected, "bad translation")
	assert.Equal(t, ts, exTrans, "bad translation")
}

func TestCustomTranslatorTrivial(t *testing.T) {
	tr := func(a pkga.Trivial) (pkgb.Nested, TranslationSet) {
		ts := mkTrans(fp("A"), fp("A"),
			fp("B"), fp("B"),
			fp("C"), fp("C"),
			fp("C"), fp("D"),
		)
		return pkgb.Nested{
			Trivial: pkgb.Trivial{
				A: a.A,
				B: a.B,
				C: a.C,
			},
			D: "abc",
		}, ts
	}
	in := pkga.Trivial{
		A: "asdf",
		B: 5,
		C: true,
	}

	expected := pkgb.Nested{
		D: "abc",
		Trivial: pkgb.Trivial{
			A: "asdf",
			B: 5,
			C: true,
		},
	}
	exTrans := mkTrans(
		fp("A"), fp("A"),
		fp("B"), fp("B"),
		fp("C"), fp("C"),
		fp("C"), fp("D"),
	)

	got := pkgb.Nested{}

	trans := NewTranslator("", "")
	trans.AddCustomTranslator(tr)

	ts := trans.Translate(&in, &got)
	assert.Equal(t, got, expected, "bad translation")
	assert.Equal(t, ts, exTrans, "bad translation")
}

func TestCustomTranslatorTrivialWithAutomaticResume(t *testing.T) {
	trans := NewTranslator("", "")
	tr := func(a pkga.Trivial) (pkgb.Nested, TranslationSet) {
		ret := pkgb.Nested{
			D: "abc",
		}
		ts := trans.Translate(&a, &ret.Trivial)
		ts.AddTranslation(fp("C"), fp("D"))
		return ret, ts
	}
	in := pkga.Trivial{
		A: "asdf",
		B: 5,
		C: true,
	}
	exTrans := mkTrans(
		fp("A"), fp("A"),
		fp("B"), fp("B"),
		fp("C"), fp("C"),
		fp("C"), fp("D"),
	)

	expected := pkgb.Nested{
		D: "abc",
		Trivial: pkgb.Trivial{
			A: "asdf",
			B: 5,
			C: true,
		},
	}

	got := pkgb.Nested{}

	trans.AddCustomTranslator(tr)

	ts := trans.Translate(&in, &got)
	assert.Equal(t, got, expected, "bad translation")
	assert.Equal(t, ts, exTrans, "bad translation")
}

func TestCustomTranslatorList(t *testing.T) {
	tr := func(a pkga.Trivial) (pkgb.Nested, TranslationSet) {
		ts := mkTrans(fp("A"), fp("A"),
			fp("B"), fp("B"),
			fp("C"), fp("C"),
			fp("C"), fp("D"),
		)
		return pkgb.Nested{
			Trivial: pkgb.Trivial{
				A: a.A,
				B: a.B,
				C: a.C,
			},
			D: "abc",
		}, ts
	}
	in := pkga.HasList{
		L: []pkga.Trivial{
			{
				A: "asdf",
				B: 5,
				C: true,
			},
		},
	}

	expected := pkgb.HasList{
		L: []pkgb.Nested{
			{
				D: "abc",
				Trivial: pkgb.Trivial{
					A: "asdf",
					B: 5,
					C: true,
				},
			},
		},
	}
	exTrans := mkTrans(
		fp("L", 0, "A"), fp("L", 0, "A"),
		fp("L", 0, "B"), fp("L", 0, "B"),
		fp("L", 0, "C"), fp("L", 0, "C"),
		fp("L", 0, "C"), fp("L", 0, "D"),
	)

	got := pkgb.HasList{}

	trans := NewTranslator("", "")
	trans.AddCustomTranslator(tr)

	ts := trans.Translate(&in, &got)
	assert.Equal(t, got, expected, "bad translation")
	assert.Equal(t, ts, exTrans, "bad translation")
}
