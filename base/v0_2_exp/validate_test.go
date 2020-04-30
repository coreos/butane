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

	"github.com/coreos/ignition/v2/config/util"
	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
)

// TestValidateResource tests that multiple sources (i.e. urls and inline) are not allowed but zero or one sources are
func TestValidateResource(t *testing.T) {
	tests := []struct {
		in  Resource
		out error
	}{
		{},
		{
			// contains invalid (by the validator's definition) combinations of fields,
			// but the translator doesn't care and we can check they all get translated at once
			Resource{
				Source:      util.StrToPtr("http://example/com"),
				Compression: util.StrToPtr("gzip"),
				Verification: Verification{
					Hash: util.StrToPtr("this isn't validated"),
				},
			},
			nil,
		},
		{
			Resource{
				Inline:      util.StrToPtr("hello"),
				Compression: util.StrToPtr("gzip"),
				Verification: Verification{
					Hash: util.StrToPtr("this isn't validated"),
				},
			},
			nil,
		},
		{
			Resource{
				Source:      util.StrToPtr("data:,hello"),
				Inline:      util.StrToPtr("hello"),
				Compression: util.StrToPtr("gzip"),
				Verification: Verification{
					Hash: util.StrToPtr("this isn't validated"),
				},
			},
			ErrInlineAndSource,
		},
	}

	for i, test := range tests {
		actual := test.in.Validate(path.New("yaml"))
		expected := report.Report{}
		// hardcode inline for now since that's the only place errors occur. Move into the
		// test struct once there's more than one place
		expected.AddOnError(path.New("yaml", "inline"), test.out)

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("#%d: expected %+v got %+v", i, expected, actual)
		}
	}
}
