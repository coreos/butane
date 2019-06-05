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

package v0_1

import (
	"github.com/coreos/ignition/v2/config/translate"
	"github.com/coreos/ignition/v2/config/v3_0/types"
)

func (c Config) ToIgn3_0() (types.Config, error) {
	ret := types.Config{}
	tr := translate.NewTranslator()
	tr.AddCustomTranslator(translateIgnition)
	tr.AddCustomTranslator(translateFile)
	tr.AddCustomTranslator(translateDirectory)
	tr.AddCustomTranslator(translateLink)
	tr.Translate(&c, &ret)
	return ret, nil
}

func translateIgnition(from Ignition) (to types.Ignition) {
	tr := translate.NewTranslator()
	to.Version = "3.0.0"
	tr.Translate(&from.Config, &to.Config)
	tr.Translate(&from.Security, &to.Security)
	tr.Translate(&from.Timeouts, &to.Timeouts)
	return
}

func translateFile(from File) (to types.File) {
	tr := translate.NewTranslator()
	tr.Translate(&from.Group, &to.Group)
	tr.Translate(&from.User, &to.User)
	tr.Translate(&from.Append, &to.Append)
	tr.Translate(&from.Contents, &to.Contents)
	to.Overwrite = from.Overwrite
	to.Path = from.Path
	to.Mode = from.Mode
	return
}

func translateDirectory(from Directory) (to types.Directory) {
	tr := translate.NewTranslator()
	tr.Translate(&from.Group, &to.Group)
	tr.Translate(&from.User, &to.User)
	to.Overwrite = from.Overwrite
	to.Path = from.Path
	to.Mode = from.Mode
	return
}

func translateLink(from Link) (to types.Link) {
	tr := translate.NewTranslator()
	tr.Translate(&from.Group, &to.Group)
	tr.Translate(&from.User, &to.User)
	to.Target = from.Target
	to.Hard = from.Hard
	to.Overwrite = from.Overwrite
	to.Path = from.Path
	return
}
