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

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/pflag"

	"github.com/coreos/fcct/config"
	"github.com/coreos/fcct/config/common"
	"github.com/coreos/fcct/internal/version"
)

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func main() {
	var (
		input       string
		output      string
		helpFlag    bool
		versionFlag bool
	)
	options := common.TranslateBytesOptions{}
	pflag.BoolVarP(&helpFlag, "help", "h", false, "show usage and exit")
	pflag.BoolVarP(&versionFlag, "version", "V", false, "print the version and exit")
	pflag.BoolVarP(&options.Strict, "strict", "s", false, "fail on any warning")
	pflag.BoolVarP(&options.Pretty, "pretty", "p", false, "output formatted json")
	pflag.StringVar(&input, "input", "", "read from input file instead of stdin")
	pflag.Lookup("input").Deprecated = "specify filename directly on command line"
	pflag.Lookup("input").Hidden = true
	pflag.StringVarP(&output, "output", "o", "", "write to output file instead of stdout")
	pflag.StringVarP(&options.FilesDir, "files-dir", "d", "", "allow embedding local files from this directory")

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [input-file]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		pflag.PrintDefaults()
	}
	pflag.Parse()

	args := pflag.Args()
	if len(args) == 1 && input == "" {
		input = args[0]
	} else if len(args) > 0 {
		pflag.Usage()
		os.Exit(2)
	}

	if helpFlag {
		pflag.Usage()
		os.Exit(0)
	}

	if versionFlag {
		fmt.Println(version.String)
		os.Exit(0)
	}

	var infile *os.File = os.Stdin
	var outfile *os.File = os.Stdout
	if input != "" {
		var err error
		infile, err = os.Open(input)
		if err != nil {
			fail("failed to open %s: %v\n", input, err)
		}
		defer infile.Close()
	}

	dataIn, err := ioutil.ReadAll(infile)
	if err != nil {
		fail("failed to read %s: %v\n", infile.Name(), err)
	}

	dataOut, r, err := config.TranslateBytes(dataIn, options)
	fmt.Fprintf(os.Stderr, "%s", r.String())
	if err != nil {
		fail("Error translating config: %v\n", err)
	}

	if output != "" {
		var err error
		outfile, err = os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fail("failed to open %s: %v\n", output, err)
		}
		defer outfile.Close()
	}

	if _, err := outfile.Write(append(dataOut, '\n')); err != nil {
		fail("Failed to write config to %s: %v\n", outfile.Name(), err)
	}
}
