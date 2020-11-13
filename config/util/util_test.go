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
	"testing"
)

func TestSnake(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{},
		{
			"foo",
			"foo",
		},
		{
			"snakeCase",
			"snake_case",
		},
		{
			"longSnakeCase",
			"long_snake_case",
		},
		{
			"snake_already",
			"snake_already",
		},
	}

	for i, test := range tests {
		if snake(test.in) != test.out {
			t.Errorf("#%d: expected %q got %q", i, test.out, snake(test.in))
		}
	}
}
