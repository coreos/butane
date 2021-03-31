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

package pkga

type Trivial struct {
	A string
	B int
	C bool
}

type Nested struct {
	D string
	Trivial
}

type TrivialReordered struct {
	B int
	A string
	C bool
}

type HasList struct {
	L []Trivial
}

type TrivialSkip struct {
	A string `butane:"auto_skip"`
	B int
	C bool
}
