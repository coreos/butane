package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ajeddeloh/fcct/config"
	"github.com/ajeddeloh/fcct/config/common"
)

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func main() {
	var (
		input    string
		output   string
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
