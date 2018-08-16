package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"github.com/pivotal-cf-experimental/stembuild/ovftool"
	"github.com/pivotal-cf-experimental/stembuild/stembuildoptions"
	"github.com/pivotal-cf-experimental/stembuild/stemcell"
	"github.com/pivotal-cf-experimental/stembuild/utils"
	"io"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

var (
	applyPatch stembuildoptions.StembuildOptions

	errs        []error
	EnableDebug bool
	DebugColor  bool
)

var Debugf = func(format string, a ...interface{}) {}

const UsageMessage = `
Usage %[1]s [-output DIRNAME] apply-patch <patch manifest yml>

Usage %[1]s [OPTIONS...] [-vmdk FILENAME]
      %[2]s [-output DIRNAME] [-version STEMCELL_VERSION]
      %[2]s [-os OS_VERSION]

Creates a BOSH stemcell from a VHD and PATCH (patch) file.

Usage:
  The VMware 'ovftool' binary must be on your path or Fusion/Workstation
  must be installed (both include the 'ovftool').

  Patch VHD:
		The [vhd], [patch], and [version] flags must be specified either on the
    command line or in the manifest file.  If the [output] flag is not
    specified the stemcell will be created in the current working directory.

  Convert VMDK [-vmdk]:
    The [vmdk] and [version] flags must be specified.  If the [output] flag is
    not specified the stemcell will be created in the current working directory.

Examples:
  %[1]s apply-patch patch-manifest.yml

  where patch-manifest.yml contains:
  ---
  version: '1.2.3'
  vhd_file: disk.vhd
  patch_file: patch.file

    Will create a stemcell using the package patch-manifest.yml (requires a VHD and
    VMDK to exist on paths in your system).

    and

  %[1]s -vmdk disk.vmdk -v 1.2

    Will create a stemcell using [vmdk] 'disk.vmdk' with version 1.2 in the current
		working directory.

Flags:
`

func Init() {
	flag.Usage = func() {
		exe := filepath.Base(os.Args[0])
		pad := strings.Repeat(" ", len(exe))
		fmt.Fprintf(os.Stderr, UsageMessage, exe, pad)
		flag.PrintDefaults()
	}

	flag.StringVar(&applyPatch.VHDFile, "vhd", "", "VHD file to patch")
	flag.StringVar(&applyPatch.VMDKFile, "vmdk", "", "VMDK file to create stemcell from")

	flag.StringVar(&applyPatch.PatchFile, "patch", "",
		"Patch file that will be applied to the VHD")
	flag.StringVar(&applyPatch.PatchFile, "d", "", "Patch file (shorthand)")

	flag.StringVar(&applyPatch.OSVersion, "os", "",
		"OS version must be either 2012R2, 2016 or 1803")

	flag.StringVar(&applyPatch.Version, "version", "",
		"Stemcell version in the form of [DIGITS].[DIGITS] (e.x. 123.01)")
	flag.StringVar(&applyPatch.Version, "v", "", "Stemcell version (shorthand)")

	flag.StringVar(&applyPatch.OutputDir, "output", "",
		"Output directory, default is the current working directory.")
	flag.StringVar(&applyPatch.OutputDir, "o", "", "Output directory (shorthand)")

	flag.BoolVar(&EnableDebug, "debug", false, "Print lots of debugging information")
	flag.BoolVar(&DebugColor, "color", false, "Colorize debug output")
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

func validPatchPath(name string) error {
	var fileError error
	var urlError error
	if fileError = validFile(name); fileError == nil {
		return nil
	}
	if urlError = validUrl(name); urlError == nil {
		return nil
	}
	return fmt.Errorf("The patchfile path is neither a valid file nor a valid url")
}

func validUrl(urlPath string) error {
	if strings.HasPrefix(urlPath, "http://") || strings.HasPrefix(urlPath, "https://") {
		_, err := url.Parse(urlPath)
		return err
	}
	return fmt.Errorf("%s is not a URL", urlPath)
}

func add(err error) {
	errs = append(errs, err)
}

func ValidateFlags() []error {
	var patchManifestFile string

	argCount := flag.NArg()
	switch {
	case argCount == 2:
		command := flag.Arg(0)
		if command != "apply-patch" {
			add(fmt.Errorf("Unrecognized command '%s'", command))
			return errs
		}
		patchManifestFile = flag.Arg(1)
	case argCount != 0:
		add(errors.New("Invalid number of arguments"))
		return errs
	}

	if patchManifestFile != "" {
		Debugf("loading 'apply patch' manifest file: %q", patchManifestFile)
		if err := stembuildoptions.LoadOptionsFromManifest(patchManifestFile, &applyPatch); err != nil {
			add(fmt.Errorf("invalid patch manifest file: %q", err))
			return errs
		}
	}

	Debugf("validating [vmdk] (%s) [vhd] (%s) and [patch] (%s) flags",
		applyPatch.VMDKFile, applyPatch.VHDFile, applyPatch.PatchFile)

	if applyPatch.VMDKFile != "" && applyPatch.VHDFile != "" {
		add(errors.New("both VMDK and VHD flags are specified"))
		return errs
	}
	if applyPatch.VMDKFile == "" && applyPatch.VHDFile == "" {
		add(errors.New("missing VMDK and VHD flags, one must be specified"))
		return errs
	}

	// check for extra flags in vmdk commmand
	Debugf("validating that no extra flags or arguments were provided")
	validateInputs()

	Debugf("validating output directory: %s", applyPatch.OutputDir)
	validateOutputDir()

	Debugf("validating integrity of provided inputs")
	validateChecksums()

	Debugf("validating stemcell version string: %s", applyPatch.Version)
	if err := utils.ValidateVersion(applyPatch.Version); err != nil {
		add(err)
		return errs
	}

	Debugf("validating OS version: %s", applyPatch.OSVersion)
	switch applyPatch.OSVersion {
	case "2012R2", "2016", "1803":
		// Ok
	default:
		add(fmt.Errorf("OS version must be either 2012R2, 2016 or 1803 have: %s", applyPatch.OSVersion))
		return errs
	}

	name := filepath.Join(applyPatch.OutputDir, stemcell.StemcellFilename(applyPatch.Version, applyPatch.OSVersion))
	Debugf("validating that stemcell filename (%s) does not exist", name)
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		add(fmt.Errorf("error with output file (%s): %v (file may already exist)", name, err))
		return errs
	}

	return errs
}

func validateInputs() {
	if applyPatch.VMDKFile != "" {
		Debugf("validating VMDK file [vmdk]: %q", applyPatch.VMDKFile)
		if err := validFile(applyPatch.VMDKFile); err != nil {
			add(fmt.Errorf("invalid [vmdk]: %s", err))
		}
	} else {
		Debugf("validating VHD file [vhd]: %q", applyPatch.VHDFile)
		if err := validFile(applyPatch.VHDFile); err != nil {
			add(fmt.Errorf("invalid [vhd]: %s", err))
		}
		Debugf("validating patch file [patch]: %q", applyPatch.PatchFile)
		if applyPatch.PatchFile == "" {
			add(errors.New("missing required argument 'patch'"))
		}
		if err := validPatchPath(applyPatch.PatchFile); err != nil {
			add(fmt.Errorf("invalid [patch]: %s", err))
		}
	}
}

func validateOutputDir() {
	if applyPatch.OutputDir == "" || applyPatch.OutputDir == "." {
		wd, err := os.Getwd()
		if err != nil {
			add(fmt.Errorf("getting working directory: %s", err))
		}
		Debugf("setting output dir (%s) to working directory: %s", applyPatch.OutputDir, wd)
		applyPatch.OutputDir = wd
	}

	fi, err := os.Stat(applyPatch.OutputDir)
	if err != nil && os.IsNotExist(err) {
		if err = os.Mkdir(applyPatch.OutputDir, 0700); err != nil {
			add(err)
		}
	} else if err != nil || fi == nil {
		add(fmt.Errorf("error opening output directory (%s): %s\n", applyPatch.OutputDir, err))
	} else if !fi.IsDir() {
		add(fmt.Errorf("output argument (%s): is not a directory\n", applyPatch.OutputDir))
	}
}

func validateChecksums() {

	if applyPatch.VHDFileChecksum != "" {
		if validateChecksum(applyPatch.VHDFile, applyPatch.VHDFileChecksum) != nil {
			add(errors.New("the specified base VHD is different from the VHD expected by the diff bundle"))
		}
	}
	if applyPatch.PatchFileChecksum != "" && validFile(applyPatch.PatchFile) == nil {
		if validateChecksum(applyPatch.PatchFile, applyPatch.PatchFileChecksum) != nil {
			add(errors.New("the specified patch file is different from the patch file expected by the diff bundle"))
		}
	}
}

func validateChecksum(filename, expectedChecksum string) error {
	hasher := sha256.New()

	fileContents, err := os.Open(filename)
	if err != nil {
		return errors.New("the specified file cannot be opened")
	}

	defer fileContents.Close()
	if _, err := io.Copy(hasher, fileContents); err != nil {
		return errors.New("the specified file cannot be loaded")
	}

	actualChecksum := hex.EncodeToString(hasher.Sum(nil))

	if expectedChecksum != actualChecksum {
		return errors.New("the actual checksum does not match the expected checksum")
	}

	return nil
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

	applyPatch.OSVersion = strings.ToUpper(applyPatch.OSVersion)

	return nil
}

func realMain(c *stemcell.Config, vmdk, vhd, patch string) error {
	start := time.Now()

	// PATCH HERE
	tmpdir, err := c.TempDir()
	if err != nil {
		return err
	}

	// This is ugly and I'm sorry
	if vmdk == "" {
		Debugf("main: creating vmdk from [vhd] (%s) and [patch] (%s)", vhd, patch)
		vmdk = filepath.Join(tmpdir, "image.vmdk")
		if err := c.ApplyPatch(vhd, patch, vmdk); err != nil {
			return err
		}
	} else {
		Debugf("main: using vmdk (%s)", vmdk)
	}

	stemcellPath, err := c.ConvertVMDK(vmdk, applyPatch.OutputDir)
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

		fmt.Fprintln(os.Stderr, "\nfor usage: stembuild -h")
		os.Exit(1)
	}

	c := stemcell.Config{
		Stop:         make(chan struct{}),
		Debugf:       Debugf,
		BuildOptions: applyPatch,
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

	if err := realMain(&c, applyPatch.VMDKFile, applyPatch.VHDFile, applyPatch.PatchFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		c.Cleanup() // remove temp dir
		os.Exit(1)
	}
	c.Cleanup() // remove temp dir
}
