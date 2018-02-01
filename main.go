package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/pivotal-cf-experimental/stembuild/ovftool"
	"github.com/pivotal-cf-experimental/stembuild/rdiff"
	"github.com/pivotal-cf-experimental/stembuild/stembuildoptions"
	"github.com/pivotal-cf-experimental/stembuild/utils"
)

const DefaultOSVersion = "2012R2"

var (
	applyPatch stembuildoptions.StembuildOptions

	EnableDebug bool
	DebugColor  bool
)

var Debugf = func(format string, a ...interface{}) {}

const UsageMessage = `
Usage %[1]s [-vhd FILENAME] [-patch FILENAME] [-output DIRNAME]
      %[2]s [-version STEMCELL_VERSION] [-os OS_VERSION]
      %[2]s [OPTIONS...] apply-patch <patch manifest yml>

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
  where patch-manifest.yml contains
  %[2]s ---
  %[2]s version: '1.2.3'
  %[2]s vhd_file: disk.vhd
  %[2]s patch_file: patch.file

  and

  %[1]s -vhd disk.vhd -patch patch.file -v 1.2.3

    Will create a stemcell using [vhd] 'disk.vhd' and a patch file with
    version 1.2.3 in the current working directory.

  %[1]s -vhd disk.vhd -patch patch.file -gzip -v 1.2 -output foo

    Will create a stemcell with version 1.2 in the 'foo' directory using gzip
    compressed patch file 'patch.file'.

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

	flag.StringVar(&applyPatch.OSVersion, "os", DefaultOSVersion,
		"OS version must be either 2012R2 or 2016")

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

func ValidateFlags() []error {
	var (
		errs              []error
		patchManifestFile string
	)

	add := func(err error) {
		errs = append(errs, err)
	}

	argCount := flag.NArg()
	switch {
	case argCount == 2:
		command := flag.Arg(0)
		if command != "apply-patch" {
			add(fmt.Errorf("Unrecognized command '%s'", command))
		}
		patchManifestFile = flag.Arg(1)
	case argCount != 0:
		add(errors.New("Invalid number of arguments"))
	}

	if patchManifestFile != "" {
		Debugf("loading 'apply patch' manifest file: %q", patchManifestFile)
		if err := stembuildoptions.LoadOptionsFromManifest(patchManifestFile, &applyPatch); err != nil {
			add(fmt.Errorf("invalid patch manifest file: %q", err))
			return errs
		}
	}

	if applyPatch.OutputDir == "" || applyPatch.OutputDir == "." {
		wd, err := os.Getwd()
		if err != nil {
			add(fmt.Errorf("getting working directory: %s", err))
			return errs
		}
		Debugf("setting output dir (%s) to working directory: %s", applyPatch.OutputDir, wd)
		applyPatch.OutputDir = wd
	}

	Debugf("validating [vmdk] (%s) [vhd] (%s) and [patch] (%s) flags",
		applyPatch.VMDKFile, applyPatch.VHDFile, applyPatch.PatchFile)

	fmt.Printf("%#v", applyPatch)
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
	if applyPatch.VMDKFile != "" && len(flag.Args()) > 0 {
		add(fmt.Errorf("extra arguments: %s\n", strings.Join(flag.Args(), ", ")))
	}

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
		if err := validFile(applyPatch.PatchFile); err != nil {
			add(fmt.Errorf("invalid [patch]: %s", err))
		}
	}

	Debugf("validating output directory: %s", applyPatch.OutputDir)
	if applyPatch.OutputDir == "" {
		add(errors.New("missing required argument 'output'"))
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

	Debugf("validating stemcell version string: %s", applyPatch.Version)
	if err := utils.ValidateVersion(applyPatch.Version); err != nil {
		add(err)
	}

	Debugf("validating OS version: %s", applyPatch.OSVersion)
	switch applyPatch.OSVersion {
	case "2012R2", "2016":
		// Ok
	default:
		add(fmt.Errorf("OS version must be either 2012R2 or 2016 have: %s", applyPatch.OSVersion))
	}

	name := filepath.Join(applyPatch.OutputDir, StemcellFilename(applyPatch.Version, applyPatch.OSVersion))
	Debugf("validating that stemcell filename (%s) does not exist", name)
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		add(fmt.Errorf("error with output file (%s): %v (file may already exist)", name, err))
	}

	return errs
}

func StemcellFilename(version, os string) string {
	return fmt.Sprintf("bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz",
		version, os)
}

var ErrInterupt = errors.New("interrupt")

type CancelReadSeeker struct {
	rs   io.ReadSeeker
	stop chan struct{}
}

func (r *CancelReadSeeker) Seek(offset int64, whence int) (int64, error) {
	select {
	case <-r.stop:
		return 0, ErrInterupt
	default:
		return r.rs.Seek(offset, whence)
	}
}

func (r *CancelReadSeeker) Read(p []byte) (int, error) {
	select {
	case <-r.stop:
		return 0, ErrInterupt
	default:
		return r.rs.Read(p)
	}
}

type CancelWriter struct {
	w    io.Writer
	stop chan struct{}
}

func (w *CancelWriter) Write(p []byte) (int, error) {
	select {
	case <-w.stop:
		return 0, ErrInterupt
	default:
		return w.w.Write(p)
	}
}

type CancelReader struct {
	r    io.Reader
	stop chan struct{}
}

func (r *CancelReader) Read(p []byte) (int, error) {
	select {
	case <-r.stop:
		return 0, ErrInterupt
	default:
		return r.r.Read(p)
	}
}

type Config struct {
	Image    string
	Stemcell string
	Manifest string
	Sha1sum  string
	tmpdir   string
	stop     chan struct{}
}

// returns a io.Writer that returns an error when Config c is stopped
func (c *Config) Writer(w io.Writer) *CancelWriter {
	return &CancelWriter{w: w, stop: c.stop}
}

// returns a io.Reader that returns an error when Config c is stopped
func (c *Config) Reader(r io.Reader) *CancelReader {
	return &CancelReader{r: r, stop: c.stop}
}

func (c *Config) Stop() {
	Debugf("stopping config")
	defer c.Cleanup() // make sure this runs!
	close(c.stop)
}

func (c *Config) Cleanup() {
	if c.tmpdir == "" {
		return
	}
	// check if directory exists to make Cleanup idempotent
	if _, err := os.Stat(c.tmpdir); err == nil {
		Debugf("deleting temp directory: %s", c.tmpdir)
		os.RemoveAll(c.tmpdir)
	}
}

func (c *Config) AddTarFile(tr *tar.Writer, name string) error {
	Debugf("adding file (%s) to tar archive", name)
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	hdr, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return err
	}
	if err := tr.WriteHeader(hdr); err != nil {
		return err
	}
	if _, err := io.Copy(tr, c.Reader(f)); err != nil {
		return err
	}
	return nil
}

func (c *Config) TempDir() (string, error) {
	if c.tmpdir != "" {
		if _, err := os.Stat(c.tmpdir); err != nil {
			Debugf("unable to stat temp dir (%s) was it deleted?", c.tmpdir)
			return "", fmt.Errorf("opening temp directory: %s", c.tmpdir)
		}
		return c.tmpdir, nil
	}
	name, err := ioutil.TempDir("", "stemcell-")
	if err != nil {
		return "", fmt.Errorf("creating temp directory: %s", err)
	}
	c.tmpdir = name
	Debugf("created temp directory: %s", name)
	return c.tmpdir, nil
}

func (c *Config) CreateStemcell() error {
	Debugf("creating stemcell")

	// programming errors - panic!
	if c.Manifest == "" {
		panic("CreateStemcell: empty manifest")
	}
	if c.Image == "" {
		panic("CreateStemcell: empty image")
	}

	tmpdir, err := c.TempDir()
	if err != nil {
		return err
	}

	c.Stemcell = filepath.Join(tmpdir, StemcellFilename(applyPatch.Version, applyPatch.OSVersion))
	stemcell, err := os.OpenFile(c.Stemcell, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer stemcell.Close()
	Debugf("created temp stemcell: %s", c.Stemcell)

	errorf := func(format string, a ...interface{}) error {
		stemcell.Close()
		os.Remove(c.Stemcell)
		return fmt.Errorf(format, a...)
	}

	t := time.Now()
	w := gzip.NewWriter(c.Writer(stemcell))
	tr := tar.NewWriter(w)

	Debugf("adding image file to stemcell tarball: %s", c.Image)
	if err := c.AddTarFile(tr, c.Image); err != nil {
		return errorf("creating stemcell: %s", err)
	}

	Debugf("adding manifest file to stemcell tarball: %s", c.Manifest)
	if err := c.AddTarFile(tr, c.Manifest); err != nil {
		return errorf("creating stemcell: %s", err)
	}

	if err := tr.Close(); err != nil {
		return errorf("creating stemcell: %s", err)
	}

	if err := w.Close(); err != nil {
		return errorf("creating stemcell: %s", err)
	}

	Debugf("created stemcell in: %s", time.Since(t))

	return nil
}

func (c *Config) WriteManifest(manifest string) error {
	// programming error - this should never happen...
	if c.Manifest != "" {
		panic("already created manifest: " + c.Manifest)
	}

	tmpdir, err := c.TempDir()
	if err != nil {
		return err
	}

	c.Manifest = filepath.Join(tmpdir, "stemcell.MF")
	f, err := os.OpenFile(c.Manifest, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("creating stemcell.MF (%s): %s", c.Manifest, err)
	}
	defer f.Close()
	Debugf("created temp stemcell.MF file: %s", c.Manifest)

	if _, err := fmt.Fprintf(f, manifest); err != nil {
		os.Remove(c.Manifest)
		return fmt.Errorf("writing stemcell.MF (%s): %s", c.Manifest, err)
	}
	Debugf("wrote stemcell.MF with sha1: %s and version: %s", c.Sha1sum, applyPatch.Version)

	return nil
}

func CreateManifest(osVersion, version, sha1sum string) string {
	const format = `---
name: bosh-vsphere-esxi-windows%[1]s-go_agent
version: '%[2]s'
sha1: %[3]s
operating_system: windows%[1]s
cloud_properties:
  infrastructure: vsphere
  hypervisor: esxi
stemcell_formats:
- vsphere-ovf
- vsphere-ova
`
	return fmt.Sprintf(format, osVersion, version, sha1sum)

}

func ExtractOVA(ova, dirname string) error {
	Debugf("extracting ova file (%s) to directory: %s", ova, dirname)
	tf, err := os.Open(ova)
	if err != nil {
		return err
	}
	defer tf.Close()
	return ExtractArchive(tf, dirname)
}

func ExtractArchive(archive io.Reader, dirname string) error {
	Debugf("extracting archive to directory: %s", dirname)

	tr := tar.NewReader(archive)

	limit := 100
	for ; limit >= 0; limit-- {
		h, err := tr.Next()
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("tar: reading from archive: %s", err)
			}
			break
		}

		// expect a flat archive
		name := h.Name
		if filepath.Base(name) != name {
			return fmt.Errorf("tar: archive contains subdirectory: %s", name)
		}

		// only allow regular files
		mode := h.FileInfo().Mode()
		if !mode.IsRegular() {
			return fmt.Errorf("tar: unexpected file mode (%s): %s", name, mode)
		}

		path := filepath.Join(dirname, name)
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
		if err != nil {
			return fmt.Errorf("tar: opening file (%s): %s", path, err)
		}
		defer f.Close()

		if _, err := io.Copy(f, tr); err != nil {
			return fmt.Errorf("tar: writing file (%s): %s", path, err)
		}
	}
	if limit <= 0 {
		return errors.New("tar: too many files in archive")
	}
	return nil
}

func (c *Config) ConvertVMX2OVA(vmx, ova string) error {
	const errFmt = "converting vmx to ova: %s\n" +
		"-- BEGIN STDERR OUTPUT -- :\n%s\n-- END STDERR OUTPUT --\n"

	ovfpath, err := ovftool.Ovftool()
	if err != nil {
		return err
	}

	// ignore stdout
	var stderr bytes.Buffer

	cmd := exec.Command(ovfpath, vmx, ova)
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("ovftool: %s", err)
	}
	Debugf("converting vmx to ova with cmd: %s %s", cmd.Path, cmd.Args)

	// Wait for process exit or interupt
	errCh := make(chan error, 1)
	go func() { errCh <- cmd.Wait() }()

	select {
	case <-c.stop:
		if cmd.Process != nil {
			Debugf("received stop signall killing ovftool process")
			cmd.Process.Kill()
		}
		return ErrInterupt
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf(errFmt, err, stderr.String())
		}
	}

	return nil
}

// ApplyPatch, applies patch file patch to base file vhd, to create file vmdk.
// It is an error if the vmdk file already exists.
func (c *Config) ApplyPatch(vhd, patch, vmdk string) error {
	Debugf("preparing to apply patch: vhd: %s patch: %s vmdk: %s", vhd, patch, vmdk)

	if _, err := os.Stat(vmdk); err == nil {
		return fmt.Errorf("creating [vmdk] file: file exists: %s", vmdk)
	}

	start := time.Now() // this is sometimes interesting

	Debugf("applying patch with rdiff")
	if err := rdiff.Patch(vhd, patch, vmdk); err != nil {
		return err
	}

	Debugf("applied patch in: %s", time.Since(start))
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

// CreateImage, converts a vmdk to a gzip compressed image file and records the
// sha1 sum of the resulting image.
func (c *Config) CreateImage(vmdk string) error {
	Debugf("Creating [image] from [vmdk]: %s", vmdk)

	tmpdir, err := c.TempDir()
	if err != nil {
		return err
	}

	var hwVersion int
	switch applyPatch.OSVersion {
	case "2012R2":
		hwVersion = 9
	case "2016":
		hwVersion = 10
	}

	vmxPath := filepath.Join(tmpdir, "image.vmx")
	vmdkPath, err := filepath.Abs(vmdk)
	if err != nil {
		return err
	}
	if err := WriteVMXTemplate(vmdkPath, hwVersion, vmxPath); err != nil {
		return err
	}

	ovaPath := filepath.Join(tmpdir, "image.ova")
	if err := c.ConvertVMX2OVA(vmxPath, ovaPath); err != nil {
		return err
	}

	// reader
	r, err := os.Open(ovaPath)
	if err != nil {
		return err
	}
	defer r.Close()

	// image file (writer)
	c.Image = filepath.Join(tmpdir, "image")
	f, err := os.OpenFile(c.Image, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// calculate sha1 while writing image file
	h := sha1.New()
	w := gzip.NewWriter(io.MultiWriter(f, h))

	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	c.Sha1sum = fmt.Sprintf("%x", h.Sum(nil))
	Debugf("Sha1 of image (%s): %s", c.Image, c.Sha1sum)
	return nil
}

func (c *Config) ConvertVMDK(vmdk string, outputDir string) (string, error) {
	if err := c.CreateImage(vmdk); err != nil {
		return "", err
	}
	if err := c.WriteManifest(CreateManifest(applyPatch.OSVersion, applyPatch.Version, c.Sha1sum)); err != nil {
		return "", err
	}
	if err := c.CreateStemcell(); err != nil {
		return "", err
	}

	stemcellPath := filepath.Join(outputDir, filepath.Base(c.Stemcell))
	Debugf("moving stemcell (%s) to: %s", c.Stemcell, stemcellPath)

	if err := os.Rename(c.Stemcell, stemcellPath); err != nil {
		return "", err
	}
	return stemcellPath, nil
}

func realMain(c *Config, vmdk, vhd, patch string) error {
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

	c := Config{stop: make(chan struct{})}

	// cleanup if interupted
	go func() {
		ch := make(chan os.Signal, 64)
		signal.Notify(ch, os.Interrupt)
		stopping := false
		for sig := range ch {
			Debugf("received signal: %s", sig)
			if stopping {
				fmt.Fprintf(os.Stderr, "received second (%s) signale - exiting now\n", sig)
				os.Exit(1)
			}
			stopping = true
			fmt.Fprintf(os.Stderr, "received (%s) signal cleaning up\n", sig)
			c.Stop()
		}
	}()

	if err := realMain(&c, applyPatch.VMDKFile, applyPatch.VHDFile, applyPatch.PatchFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
	c.Cleanup() // remove temp dir
}
