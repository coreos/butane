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
	"errors"
	"testing"

	"github.com/coreos/fcct/translate/tests/pkga"
	"github.com/coreos/fcct/translate/tests/pkgb"

	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
	"github.com/stretchr/testify/assert"
)

type testOptions struct{}

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

	trans := NewTranslator("", "", testOptions{})

	ts, r := trans.Translate(&in, &got)
	assert.Equal(t, got, expected, "bad translation")
	assert.Equal(t, ts, exTrans, "bad translation")
	assert.Equal(t, r.String(), "", "non-empty report")
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

	trans := NewTranslator("", "", testOptions{})

	ts, r := trans.Translate(&in, &got)
	assert.Equal(t, got, expected, "bad translation")
	assert.Equal(t, ts, exTrans, "bad translation")
	assert.Equal(t, r.String(), "", "non-empty report")
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

	trans := NewTranslator("", "", testOptions{})

	ts, r := trans.Translate(&in, &got)
	assert.Equal(t, got, expected, "bad translation")
	assert.Equal(t, ts, exTrans, "bad translation")
	assert.Equal(t, r.String(), "", "non-empty report")
}

func TestTranslateTrivialSkip(t *testing.T) {
	in := pkga.TrivialSkip{
		A: "asdf",
		B: 5,
		C: true,
	}

	expected := pkgb.TrivialSkip{
		B: 5,
		C: true,
	}
	exTrans := mkTrans(
		fp("B"), fp("B"),
		fp("C"), fp("C"),
	)

	got := pkgb.TrivialSkip{}

	trans := NewTranslator("", "", testOptions{})

	ts, r := trans.Translate(&in, &got)
	assert.Equal(t, got, expected, "bad translation")
	assert.Equal(t, ts, exTrans, "bad translation")
	assert.Equal(t, r.String(), "", "non-empty report")
}

func TestCustomTranslatorTrivial(t *testing.T) {
	tr := func(a pkga.Trivial, options testOptions) (pkgb.Nested, TranslationSet, report.Report) {
		ts := mkTrans(fp("A"), fp("A"),
			fp("B"), fp("B"),
			fp("C"), fp("C"),
			fp("C"), fp("D"),
		)
		var r report.Report
		r.AddOnInfo(fp("A"), errors.New("info"))
		return pkgb.Nested{
			Trivial: pkgb.Trivial{
				A: a.A,
				B: a.B,
				C: a.C,
			},
			D: "abc",
		}, ts, r
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

	trans := NewTranslator("", "", testOptions{})
	trans.AddCustomTranslator(tr)

	ts, r := trans.Translate(&in, &got)
	assert.Equal(t, got, expected, "bad translation")
	assert.Equal(t, ts, exTrans, "bad translation")
	assert.Equal(t, r.String(), "info at $.A: info\n", "bad report")
}

func TestCustomTranslatorTrivialWithAutomaticResume(t *testing.T) {
	trans := NewTranslator("", "", testOptions{})
	tr := func(a pkga.Trivial, options testOptions) (pkgb.Nested, TranslationSet, report.Report) {
		ret := pkgb.Nested{
			D: "abc",
		}
		ts, r := trans.Translate(&a, &ret.Trivial)
		ts.AddTranslation(fp("C"), fp("D"))
		return ret, ts, r
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

	ts, r := trans.Translate(&in, &got)
	assert.Equal(t, got, expected, "bad translation")
	assert.Equal(t, ts, exTrans, "bad translation")
	assert.Equal(t, r.String(), "", "non-empty report")
}

func TestCustomTranslatorList(t *testing.T) {
	tr := func(a pkga.Trivial, options testOptions) (pkgb.Nested, TranslationSet, report.Report) {
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
		}, ts, report.Report{}
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

	trans := NewTranslator("", "", testOptions{})
	trans.AddCustomTranslator(tr)

	ts, r := trans.Translate(&in, &got)
	assert.Equal(t, got, expected, "bad translation")
	assert.Equal(t, ts, exTrans, "bad translation")
	assert.Equal(t, r.String(), "", "non-empty report")
}

func TestAddIdentity(t *testing.T) {
	ts := NewTranslationSet("1", "2")
	ts.AddIdentity("foo", "bar")
	expectedFoo := Translation{
		From: path.New("1", "foo"),
		To:   path.New("2", "foo"),
	}
	expectedBar := Translation{
		From: path.New("1", "bar"),
		To:   path.New("2", "bar"),
	}
	expectedFoo2 := Translation{
		From: path.New("1", "pre", "foo"),
		To:   path.New("2", "pre", "foo"),
	}
	expectedBar2 := Translation{
		From: path.New("1", "pre", "bar"),
		To:   path.New("2", "pre", "bar"),
	}
	ts2 := NewTranslationSet("1", "2")
	ts2.MergeP("pre", ts)
	ts3 := NewTranslationSet("1", "2")
	ts3.Merge(ts.Prefix("pre"))

	assert.Equal(t, ts.Set["$.foo"], expectedFoo, "foo not added correctly")
	assert.Equal(t, ts.Set["$.bar"], expectedBar, "bar not added correctly")
	assert.Equal(t, ts2.Set["$.pre.foo"], expectedFoo2, "foo not added correctly")
	assert.Equal(t, ts3.Set["$.pre.bar"], expectedBar2, "bar not added correctly")
	assert.Equal(t, ts3.Set["$.pre.foo"], expectedFoo2, "foo not added correctly")
	assert.Equal(t, ts2.Set["$.pre.bar"], expectedBar2, "bar not added correctly")
}
