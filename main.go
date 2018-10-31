package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/pivotal-cf-experimental/stembuild/ovftool"
	"github.com/pivotal-cf-experimental/stembuild/stembuildoptions"
	"github.com/pivotal-cf-experimental/stembuild/stemcell"
	"github.com/pivotal-cf-experimental/stembuild/utils"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

var (
	stembuildOptions stembuildoptions.StembuildOptions

	errs        []error
	EnableDebug bool
	DebugColor  bool
)

var Debugf = func(format string, a ...interface{}) {}

const UsageMessage = `
Usage %[1]s [OPTIONS...] -vmdk FILENAME
                             -version STEMCELL_VERSION
                             -os OS_VERSION
                             [-output DIRNAME] 

Create a BOSH Stemcell from a VMDK file

Usage:
  The VMware 'ovftool' binary must be on your path or Fusion/Workstation
  must be installed (both include the 'ovftool').

  Convert VMDK [-vmdk]:
    The [vmdk], [version], and [os] flags must be specified.  If the [output] flag is
    not specified the stemcell will be created in the current working directory.

Examples:

  %[1]s -vmdk disk.vmdk -v 1.2

    Will create a stemcell using [vmdk] 'disk.vmdk' with version 1.2 in the current
        working directory.

Flags:
`

func Init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, UsageMessage, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.BoolVar(&DebugColor, "color", false, "Colorize debug output")

	flag.BoolVar(&EnableDebug, "debug", false, "Print lots of debugging information")

	flag.StringVar(&stembuildOptions.OutputDir, "output", "",
		"Output directory, default is the current working directory.")
	flag.StringVar(&stembuildOptions.OutputDir, "o", "", "Output directory (shorthand)")

	flag.StringVar(&stembuildOptions.OSVersion, "os", "",
		"OS version must be either 2012R2, 2016 or 1803")

	flag.StringVar(&stembuildOptions.Version, "version", "",
		"Stemcell version in the form of [DIGITS].[DIGITS] (e.x. 123.01)")
	flag.StringVar(&stembuildOptions.Version, "v", "", "Stemcell version (shorthand)")

	flag.StringVar(&stembuildOptions.VMDKFile, "vmdk", "", "VMDK file to create stemcell from")

}

func Usage() {
	flag.Usage()
	os.Exit(1)
}

func validFile(name string) error {
	fi, err := os.Stat(name)
	if err != nil {
		return err
	}
	if !fi.Mode().IsRegular() {
		return fmt.Errorf("not a regular file: %s", name)
	}
	return nil
}

func add(err error) {
	errs = append(errs, err)
}

func ValidateFlags() []error {
	Debugf("validating [vmdk] (%s) flags", stembuildOptions.VMDKFile)

	if stembuildOptions.VMDKFile == "" {
		add(errors.New("missing VMDK flag"))
		return errs
	}

	// check for extra flags in vmdk commmand
	Debugf("validating that no extra flags or arguments were provided")
	validateInputs()

	Debugf("validating output directory: %s", stembuildOptions.OutputDir)
	validateOutputDir()

	Debugf("validating stemcell version string: %s", stembuildOptions.Version)
	if err := utils.ValidateVersion(stembuildOptions.Version); err != nil {
		add(err)
		return errs
	}

	Debugf("validating OS version: %s", stembuildOptions.OSVersion)
	switch stembuildOptions.OSVersion {
	case "2012R2", "2016", "1803":
		// Ok
	default:
		add(fmt.Errorf("OS version must be either 2012R2, 2016 or 1803 have: %s", stembuildOptions.OSVersion))
		return errs
	}

	name := filepath.Join(stembuildOptions.OutputDir, stemcell.StemcellFilename(stembuildOptions.Version, stembuildOptions.OSVersion))
	Debugf("validating that stemcell filename (%s) does not exist", name)
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		add(fmt.Errorf("error with output file (%s): %v (file may already exist)", name, err))
		return errs
	}

	return errs
}

func validateInputs() {
	Debugf("validating VMDK file [vmdk]: %q", stembuildOptions.VMDKFile)
	if err := validFile(stembuildOptions.VMDKFile); err != nil {
		add(fmt.Errorf("invalid [vmdk]: %s", err))
	}
}

func validateOutputDir() {
	if stembuildOptions.OutputDir == "" || stembuildOptions.OutputDir == "." {
		wd, err := os.Getwd()
		if err != nil {
			add(fmt.Errorf("getting working directory: %s", err))
		}
		Debugf("setting output dir (%s) to working directory: %s", stembuildOptions.OutputDir, wd)
		stembuildOptions.OutputDir = wd
	}

	fi, err := os.Stat(stembuildOptions.OutputDir)
	if err != nil && os.IsNotExist(err) {
		if err = os.Mkdir(stembuildOptions.OutputDir, 0700); err != nil {
			add(err)
		}
	} else if err != nil || fi == nil {
		add(fmt.Errorf("error opening output directory (%s): %s\n", stembuildOptions.OutputDir, err))
	} else if !fi.IsDir() {
		add(fmt.Errorf("output argument (%s): is not a directory\n", stembuildOptions.OutputDir))
	}
}

func ParseFlags() error {
	flag.Parse()

	if EnableDebug {
		if DebugColor {
			Debugf = log.New(os.Stderr, "\033[32m"+"debug: "+"\033[0m", 0).Printf
		} else {
			Debugf = log.New(os.Stderr, "debug: ", 0).Printf
		}
		Debugf("enabled")
	}

	stembuildOptions.OSVersion = strings.ToUpper(stembuildOptions.OSVersion)

	return nil
}

func realMain(c *stemcell.Config, vmdk string) error {
	start := time.Now()

	stemcellPath, err := c.ConvertVMDK(vmdk, stembuildOptions.OutputDir)
	if err != nil {
		return err
	}

	Debugf("created stemcell (%s) in: %s", stemcellPath, time.Since(start))
	fmt.Println("created stemcell:", stemcellPath)

	return nil
}

func main() {
	Init()

	Debugf("parsing flags")
	if err := ParseFlags(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		Usage()
	}

	path, err := ovftool.Ovftool()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not locate 'ovftool' on PATH: %s", err)
		Usage()
	}
	Debugf("using 'ovftool' found at: %s", path)

	if errs := ValidateFlags(); errs != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid arguments")

		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "  %s\n", e)
		}

		fmt.Fprintf(os.Stderr, "\nfor usage: %s -h \n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	c := stemcell.Config{
		Stop:         make(chan struct{}),
		Debugf:       Debugf,
		BuildOptions: stembuildOptions,
	}

	// cleanup if interrupted
	go func() {
		ch := make(chan os.Signal, 64)
		signal.Notify(ch, os.Interrupt)
		stopping := false
		for sig := range ch {
			Debugf("received signal: %s", sig)
			if stopping {
				fmt.Fprintf(os.Stderr, "received second (%s) signal - exiting now\n", sig)
				c.Cleanup() // remove temp dir
				os.Exit(1)
			}
			stopping = true
			fmt.Fprintf(os.Stderr, "received (%s) signal cleaning up\n", sig)
			c.StopConfig()
		}
	}()

	if err := realMain(&c, stembuildOptions.VMDKFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		c.Cleanup() // remove temp dir
		os.Exit(1)
	}
	c.Cleanup() // remove temp dir
}
