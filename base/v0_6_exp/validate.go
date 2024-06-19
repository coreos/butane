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

package v0_6_exp

import (
	baseutil "github.com/coreos/butane/base/util"
	"github.com/coreos/butane/config/common"
	"github.com/coreos/ignition/v2/config/shared/errors"
	"github.com/coreos/ignition/v2/config/util"
	exp "github.com/coreos/ignition/v2/config/v3_5_experimental"
	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
)

func (rs Resource) Validate(c path.ContextPath) (r report.Report) {
	var field string
	sources := 0
	var config string
	var butaneReport report.Report
	if rs.Local != nil {
		sources++
		field = "local"
		config = *rs.Local
	}
	if rs.Inline != nil {
		sources++
		field = "inline"
		config = *rs.Inline
	}
	if rs.Source != nil {
		sources++
		field = "source"
		config = *rs.Source
	}
	if sources > 1 {
		r.AddOnError(c.Append(field), common.ErrTooManyResourceSources)
	}
	if field == "local" || field == "inline" {
		_, report, err := exp.Parse([]byte(config))
		if len(report.Entries) > 0 {
			butaneReport = ConvertToButaneReport(report, field)
			r.Merge(butaneReport)
		}
		if err != nil {
			r.AddOnError(c.Append(field), errors.ErrUnknownVersion)
		}
	}
	return
}

func ConvertToButaneReport(ignitionReport report.Report, field string) report.Report {
	var butaneRep report.Report
	for _, entry := range ignitionReport.Entries {
		
		adjustedPath := []interface{}{field}
		adjustedPath = append(adjustedPath, entry.Context.Path...)


		butaneEntry := report.Entry{
			Kind:    entry.Kind,
			Message: entry.Message,
			Context: path.ContextPath{
				Path: adjustedPath,      // convert ignition path to butane path
				Tag:  entry.Context.Tag,
			},
			Marker: entry.Marker,
		}
		butaneRep.Entries = append(butaneRep.Entries, butaneEntry)
	}
	return butaneRep
}

// func TranslatePath() {
// path translating logic // TODO
// }

func (fs Filesystem) Validate(c path.ContextPath) (r report.Report) {
	if !util.IsTrue(fs.WithMountUnit) {
		return
	}
	if util.NilOrEmpty(fs.Format) {
		r.AddOnError(c.Append("format"), common.ErrMountUnitNoFormat)
	} else if *fs.Format != "swap" && util.NilOrEmpty(fs.Path) {
		r.AddOnError(c.Append("path"), common.ErrMountUnitNoPath)
	}
	return
}

func (d Directory) Validate(c path.ContextPath) (r report.Report) {
	if d.Mode != nil {
		r.AddOnWarn(c.Append("mode"), baseutil.CheckForDecimalMode(*d.Mode, true))
	}
	return
}

func (f File) Validate(c path.ContextPath) (r report.Report) {
	if f.Mode != nil {
		r.AddOnWarn(c.Append("mode"), baseutil.CheckForDecimalMode(*f.Mode, false))
	}
	return
}

func (t Tree) Validate(c path.ContextPath) (r report.Report) {
	if t.Local == "" {
		r.AddOnError(c, common.ErrTreeNoLocal)
	}
	return
}

func (rs Unit) Validate(c path.ContextPath) (r report.Report) {
	if rs.ContentsLocal != nil && rs.Contents != nil {
		r.AddOnError(c.Append("contents_local"), common.ErrTooManySystemdSources)
	}
	return
}

func (rs Dropin) Validate(c path.ContextPath) (r report.Report) {
	if rs.ContentsLocal != nil && rs.Contents != nil {
		r.AddOnError(c.Append("contents_local"), common.ErrTooManySystemdSources)
	}
	return
}
