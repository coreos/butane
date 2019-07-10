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
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/coreos/fcct/config"
	"github.com/coreos/fcct/config/common"
)

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func main() {
	var (
		input  string
		output string
	)
	options := common.TranslateOptions{}
	flag.BoolVar(&options.Strict, "strict", false, "fail on any warning")
	flag.BoolVar(&options.Pretty, "pretty", false, "output formatted json")
	flag.StringVar(&input, "input", "", "read from input file instead of stdin")
	flag.StringVar(&output, "output", "", "write to output file instead of stdout")

	flag.Parse()

	var infile *os.File = os.Stdin
	var outfile *os.File = os.Stdout
	if input != "" {
		var err error
		infile, err = os.Open(input)
		if err != nil {
			fail("failed to open %s: %v", input, err)
		}
		defer infile.Close()
	}

	dataIn, err := ioutil.ReadAll(infile)
	if err != nil {
		fail("failed to read %s: %v", infile.Name(), err)
	}

	dataOut, err := config.Translate(dataIn, options)
	if err != nil {
		fail("Error translating config: %v", err)
	}

	if output != "" {
		var err error
		outfile, err = os.Open(output)
		if err != nil {
			fail("failed to open %s: %v", output, err)
		}
		defer outfile.Close()
	}

	if _, err := outfile.Write(dataOut); err != nil {
		fail("Failed to write config to %s: %v", outfile.Name(), err)
	}
}
