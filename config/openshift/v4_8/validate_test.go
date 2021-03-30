// Copyright 2021 Red Hat, Inc
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

package v4_8

import (
	"testing"

	"github.com/coreos/fcct/config/common"

	"github.com/coreos/ignition/v2/config/shared/errors"
	"github.com/coreos/ignition/v2/config/util"
	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
	"github.com/stretchr/testify/assert"
)

func TestValidateMetadata(t *testing.T) {
	tests := []struct {
		in      Metadata
		out     error
		errPath path.ContextPath
	}{
		// missing name
		{
			Metadata{
				Labels: map[string]string{
					ROLE_LABEL_KEY: "r",
				},
			},
			common.ErrNameRequired,
			path.New("yaml", "name"),
		},
		// missing role
		{
			Metadata{
				Name: "n",
			},
			common.ErrRoleRequired,
			path.New("yaml", "labels", ROLE_LABEL_KEY),
		},
		// empty role
		{
			Metadata{
				Name: "n",
				Labels: map[string]string{
					ROLE_LABEL_KEY: "",
				},
			},
			common.ErrRoleRequired,
			path.New("yaml", "labels", ROLE_LABEL_KEY),
		},
	}

	for i, test := range tests {
		actual := test.in.Validate(path.New("yaml"))
		expected := report.Report{}
		expected.AddOnError(test.errPath, test.out)
		assert.Equal(t, expected, actual, "#%d: bad report", i)
	}
}

func TestValidateOpenShift(t *testing.T) {
	tests := []struct {
		in      OpenShift
		out     error
		errPath path.ContextPath
	}{
		// empty struct
		{
			OpenShift{},
			nil,
			path.New("yaml"),
		},
		// bad kernel type
		{
			OpenShift{
				KernelType: util.StrToPtr("hurd"),
			},
			common.ErrInvalidKernelType,
			path.New("yaml", "kernel_type"),
		},
	}

	for i, test := range tests {
		actual := test.in.Validate(path.New("yaml"))
		expected := report.Report{}
		expected.AddOnError(test.errPath, test.out)
		assert.Equal(t, expected, actual, "#%d: bad report", i)
	}
}

// TestReportCorrelation tests that errors are correctly correlated to their source lines
func TestReportCorrelation(t *testing.T) {
	tests := []struct {
		in      string
		message string
		line    int64
	}{
		// FCCT unused key check
		{
			`
                         metadata:
                           name: something
                           labels:
                             machineconfiguration.openshift.io/role: r
                         storage:
                           files:
                           - path: /z
                             q: z`,
			"Unused key q",
			9,
		},
		// FCCT YAML validation error
		{
			`
                         metadata:
                           name: something
                           labels:
                             machineconfiguration.openshift.io/role: r
                         storage:
                           files:
                           - path: /z
                             contents:
                               source: https://example.com
                               inline: z`,
			common.ErrTooManyResourceSources.Error(),
			10,
		},
		// FCCT YAML validation warning
		{
			`
                         metadata:
                           name: something
                           labels:
                             machineconfiguration.openshift.io/role: r
                         storage:
                           files:
                           - path: /z
                             mode: 644`,
			common.ErrDecimalMode.Error(),
			9,
		},
		// FCCT translation error
		{
			`
                         metadata:
                           name: something
                           labels:
                             machineconfiguration.openshift.io/role: r
                         storage:
                           files:
                           - path: /z
                             contents:
                               local: z`,
			common.ErrNoFilesDir.Error(),
			10,
		},
		// Ignition validation error, leaf node
		{
			`
                         metadata:
                           name: something
                           labels:
                             machineconfiguration.openshift.io/role: r
                         storage:
                           files:
                           - path: z`,
			errors.ErrPathRelative.Error(),
			8,
		},
		// Ignition validation error, partition
		{
			`
                         metadata:
                           name: something
                           labels:
                             machineconfiguration.openshift.io/role: r
                         storage:
                           disks:
                           - device: /dev/z
                             partitions:
                               - start_mib: 5`,
			errors.ErrNeedLabelOrNumber.Error(),
			10,
		},
		// Ignition validation error, partition list
		{
			`
                         metadata:
                           name: something
                           labels:
                             machineconfiguration.openshift.io/role: r
                         storage:
                           disks:
                           - device: /dev/z
                             partitions:
                               - number: 1
                                 should_exist: false
                               - label: z`,
			errors.ErrZeroesWithShouldNotExist.Error(),
			10,
		},
		// Ignition duplicate key check, paths
		{
			`
                         metadata:
                           name: something
                           labels:
                             machineconfiguration.openshift.io/role: r
                         storage:
                           files:
                           - path: /z
                           - path: /z`,
			errors.ErrDuplicate.Error(),
			9,
		},
	}

	for i, test := range tests {
		for _, raw := range []bool{false, true} {
			_, r, _ := ToConfigBytes([]byte(test.in), common.TranslateBytesOptions{
				Raw: raw,
			})
			assert.Len(t, r.Entries, 1, "#%d: unexpected report length, raw %v", i, raw)
			assert.Equal(t, test.message, r.Entries[0].Message, "#%d: bad error, raw %v", i, raw)
			assert.NotNil(t, r.Entries[0].Marker.StartP, "#%d: marker start is nil, raw %v", i, raw)
			assert.Equal(t, test.line, r.Entries[0].Marker.StartP.Line, "#%d: incorrect error line, raw %v", i, raw)
		}
	}
}
