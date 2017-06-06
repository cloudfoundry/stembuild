package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pivotal-cf-experimental/pcf-make-stemcell/rdiff"
)

const UsageMessage = `Usage: %s [OPTIONS] signature [BASIS [SIGNATURE]]
          [OPTIONS] delta SIGNATURE [NEWFILE [DELTA]]
          [OPTIONS] patch BASIS [DELTA [NEWFILE]]
`

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, UsageMessage, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func Usage() {
	flag.Usage()
	os.Exit(1)
}

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		Usage()
	}
	switch args[0] {
	case "signature":
		if len(args) != 3 {
			Usage()
		}
		basis := args[1]
		signature := args[2]
		if err := rdiff.Signature(basis, signature, false); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "delta":
		if len(args) != 4 {
			Usage()
		}
		signature := args[1]
		newfile := args[2]
		delta := args[3]
		if err := rdiff.Delta(signature, newfile, delta); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "patch":
		if len(args) != 4 {
			Usage()
		}
		basis := args[1]
		delta := args[2]
		newfile := args[3]
		if err := rdiff.Patch(basis, delta, newfile); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid args: %s\n", args)
		Usage()
	}
}
