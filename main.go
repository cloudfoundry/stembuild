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
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/pivotal-cf-experimental/stembuild/ovftool"
	"github.com/pivotal-cf-experimental/stembuild/rdiff"
)

const DefaultOSVersion = "2012R2"

var (
	Version     string
	OutputDir   string
	OvaFile     string
	OvfDir      string
	VHDFile     string
	VMDKFile    string
	OSVersion   string
	DeltaFile   string
	EnableDebug bool
	DebugColor  bool
)

var Debugf = func(format string, a ...interface{}) {}

const UsageMessage = `
Usage %[1]s: [OPTIONS...] [-VMDK FILENAME] [[-VHD FILENAME]
      %[2]s  [-DELTA FILENAME]] [-OUTPUT DIRNAME]
      %[2]s  [-VERSION STEMCELL_VERSION] [-OS OS_VERSION]

Creates a BOSH stemcell from a VHD and DELTA (patch) file.

Usage:
  The VMware 'ovftool' binary must be on your path or Fusion/Workstation
  must be installed (both include the 'ovftool').

  Patch VHD [-VHD]:

    The [vhd], [delta] and [version] flags must be specified.  If the [output]
    flag is not specified the stemcell will be created in the current working
    directory.

  Convert VMDK [-VMDK]:
    The [vmdk] and [version] flags must be specified.  If the [output] flag is
    not specified the stemcell will be created in the current working directory.

Examples:
  %[1]s -vmdk disk.vmdk -v 1.2

    Will create a stemcell using [vmdk] 'disk.vmkd' with version 1.2 in the current
		working directory.

  %[1]s -vhd disk.vhd -delta patch.file -v 1.2

    Will create a stemcell using [vhd] 'disk.vhd' and a patch file with
		version 1.2 in the current working directory.

  %[1]s -vhd disk.vhd -delta patch.file -gzip -v 1.2 -output foo

    Will create a stemcell with version 1.2 in the 'foo' directory using gzip
    compressed patch file 'patch.file'.

Flags:
`

func Init() {
	flag.Usage = func() {
		exe := filepath.Base(os.Args[0])
		pad := strings.Repeat(" ", len(exe))
		fmt.Fprintf(os.Stderr, UsageMessage, exe, pad)
		flag.PrintDefaults()
	}

	flag.StringVar(&VHDFile, "vhd", "", "VHD file to patch")
	flag.StringVar(&VMDKFile, "vmdk", "", "VMDK file to create stemcell from")

	flag.StringVar(&DeltaFile, "delta", "",
		"Patch file that will be applied to the VHD")
	flag.StringVar(&DeltaFile, "d", "", "Patch file (shorthand)")

	flag.StringVar(&OSVersion, "os", DefaultOSVersion,
		"OS version must be either 2012R2 or 2016")

	flag.StringVar(&Version, "version", "",
		"Stemcell version in the form of [DIGITS].[DIGITS] (e.x. 123.01)")
	flag.StringVar(&Version, "v", "", "Stemcell version (shorthand)")

	flag.StringVar(&OutputDir, "output", "",
		"Output directory, default is the current working directory.")
	flag.StringVar(&OutputDir, "o", "", "Output directory (shorthand)")

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
	Debugf("validating [vmdk] (%s) [vhd] (%s) and [delta] (%s) flags",
		VMDKFile, VHDFile, DeltaFile)

	var errs []error
	add := func(err error) {
		errs = append(errs, err)
	}

	if VMDKFile != "" && VHDFile != "" {
		add(errors.New("both VMDK and VHD flags are specified"))
		return errs
	}
	if VMDKFile == "" && VHDFile == "" {
		add(errors.New("missing VMDK and VHD flags, one must be specified"))
		return errs
	}

	// check for extra flags
	Debugf("validating that no extra flags or arguments were provided")
	if n := len(flag.Args()); n != 0 {
		add(fmt.Errorf("extra arguments: %s\n", strings.Join(flag.Args(), ", ")))
	}

	if VMDKFile != "" {
		Debugf("validating VMDK file [vmdk]: %q", VMDKFile)
		if err := validFile(VMDKFile); err != nil {
			add(fmt.Errorf("invalid [vmdk]: %s", err))
		}
	} else {
		Debugf("validating VHD file [vhd]: %q", VHDFile)
		if err := validFile(VHDFile); err != nil {
			add(fmt.Errorf("invalid [vhd]: %s", err))
		}
		Debugf("validating patch file [delta]: %q", DeltaFile)
		if DeltaFile == "" {
			add(errors.New("missing required argument 'delta'"))
		}
		if err := validFile(DeltaFile); err != nil {
			add(fmt.Errorf("invalid [delta]: %s", err))
		}
	}

	Debugf("validating output directory: %s", OutputDir)
	if OutputDir == "" {
		add(errors.New("missing required argument 'output'"))
	}
	fi, err := os.Stat(OutputDir)
	if err != nil || fi == nil {
		add(fmt.Errorf("error opening output directory (%s): %s\n", OutputDir, err))
	} else if !fi.IsDir() {
		add(fmt.Errorf("output argument (%s): is not a directory\n", OutputDir))
	}

	Debugf("validating stemcell version string: %s", Version)
	if err := validateVersion(Version); err != nil {
		add(err)
	}

	Debugf("validating OS version: %s", OSVersion)
	switch OSVersion {
	case "2012R2", "2016":
		// Ok
	default:
		add(fmt.Errorf("OS version must be either 2012R2 or 2016 have: %s", OSVersion))
	}

	name := filepath.Join(OutputDir, StemcellFilename(Version, OSVersion))
	Debugf("validating that stemcell filename (%s) does not exist", name)
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		add(fmt.Errorf("file (%s) already exists - refusing to overwrite", name))
	}

	return errs
}

func validateVersion(s string) error {
	Debugf("validating version string: %s", s)
	if s == "" {
		return errors.New("missing required argument 'version'")
	}
	const pattern = `^\d{1,}.\d{1,}$`
	if !regexp.MustCompile(pattern).MatchString(s) {
		Debugf("expected version string to match regex: '%s'", pattern)
		return fmt.Errorf("invalid version (%s) expected format [NUMBER].[NUMBER]", s)
	}
	return nil
}

func StemcellFilename(version, os string) string {
	return fmt.Sprintf("bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz",
		version, os)
}

var ErrInterupt = errors.New("interupt")

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
	name, err := TempDir("stemcell-")
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

	c.Stemcell = filepath.Join(tmpdir, StemcellFilename(Version, OSVersion))
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

func (c *Config) WriteManifest() error {
	const format = `---
name: bosh-vsphere-esxi-windows%[1]s-go_agent
version: %[2]s
sha1: %[3]s
operating_system: windows%[1]s
cloud_properties:
  infrastructure: vsphere
  hypervisor: esxi
`

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

	if _, err := fmt.Fprintf(f, format, OSVersion, Version, c.Sha1sum); err != nil {
		os.Remove(c.Manifest)
		return fmt.Errorf("writing stemcell.MF (%s): %s", c.Manifest, err)
	}
	Debugf("wrote stemcell.MF with sha1: %s and version: %s", c.Sha1sum, Version)

	return nil
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
			Debugf("recieved stop signall killing ovftool process")
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

// ApplyPatch, applies patch file delta to base file vhd, to create file vmdk.
// It is an error if the vmdk file already exists.
func (c *Config) ApplyPatch(vhd, delta, vmdk string) error {
	Debugf("preparing to apply patch: vhd: %s delta: %s vmdk: %s", vhd, delta, vmdk)

	if _, err := os.Stat(vmdk); err == nil {
		return fmt.Errorf("creating [vmdk] file: file exists: %s", vmdk)
	}

	start := time.Now() // this is sometimes interesting

	Debugf("applying patch with rdiff")
	if err := rdiff.Patch(vhd, delta, vmdk); err != nil {
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

	if OutputDir == "" || OutputDir == "." {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting working directory: %s", err)
		}
		Debugf("setting output dir (%s) to working directory: %s", OutputDir, wd)
		OutputDir = wd
	}

	OSVersion = strings.ToUpper(OSVersion)

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

	vmxPath := filepath.Join(tmpdir, "image.vmx")
	vmdkPath, err := filepath.Abs(vmdk)
	if err != nil {
		return err
	}
	if err := WriteVMXTemplate(vmdkPath, vmxPath); err != nil {
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

func (c *Config) ConvertVMDK(vmdk string) (string, error) {
	if err := c.CreateImage(vmdk); err != nil {
		return "", err
	}
	if err := c.WriteManifest(); err != nil {
		return "", err
	}
	if err := c.CreateStemcell(); err != nil {
		return "", err
	}

	stemcellPath := filepath.Join(OutputDir, filepath.Base(c.Stemcell))
	Debugf("moving stemcell (%s) to: %s", c.Stemcell, stemcellPath)

	if err := os.Rename(c.Stemcell, stemcellPath); err != nil {
		return "", err
	}
	return stemcellPath, nil
}

func realMain(c *Config, vmdk, vhd, delta string) error {
	start := time.Now()

	// PATCH HERE
	tmpdir, err := c.TempDir()
	if err != nil {
		return err
	}

	// This is ugly and I'm sorry
	if vmdk == "" {
		Debugf("main: creating vmdk from [vhd] (%s) and [delta] (%s)", vhd, delta)
		vmdk = filepath.Join(tmpdir, "image.vmdk")
		if err := c.ApplyPatch(vhd, delta, vmdk); err != nil {
			return err
		}
	} else {
		Debugf("main: using vmdk (%s)", vmdk)
	}

	stemcellPath, err := c.ConvertVMDK(vmdk)
	if err != nil {
		return err
	}

	Debugf("created stemcell (%s) in: %s", stemcellPath, time.Since(start))
	fmt.Println("created stemell:", stemcellPath)

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
		Usage()
	}

	c := Config{stop: make(chan struct{})}

	// cleanup if interupted
	go func() {
		ch := make(chan os.Signal, 64)
		signal.Notify(ch, os.Interrupt)
		stopping := false
		for sig := range ch {
			Debugf("recieved signal: %s", sig)
			if stopping {
				fmt.Fprintf(os.Stderr, "recieved second (%s) signale - exiting now\n", sig)
				os.Exit(1)
			}
			stopping = true
			fmt.Fprintf(os.Stderr, "recieved (%s) signal cleaning up\n", sig)
			c.Stop()
		}
	}()

	if err := realMain(&c, VMDKFile, VHDFile, DeltaFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
	c.Cleanup() // remove temp dir
}
