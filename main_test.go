package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pivotal-cf-experimental/stembuild/helpers"
)

func TestCreateImage(t *testing.T) {
	const archive = "testdata/patch-test.tar.gz"

	dirname := helpers.ExtractGzipArchive(t, archive)
	defer os.RemoveAll(dirname)

	vmdkPath := filepath.Join(dirname, "expected.vmdk")

	conf := Config{stop: make(chan struct{})}

	if err := conf.CreateImage(vmdkPath); err != nil {
		t.Errorf("CreateImage: %s", err)
	}

	// the image will be saved to the Config's temp directory
	tmpdir, err := conf.TempDir()
	if err != nil {
		t.Error(err)
	}
	expImagePath := filepath.Join(tmpdir, "image")

	if conf.Image != expImagePath {
		t.Errorf("CreateImage: expected ImagePath to be: %s got: %s",
			expImagePath, conf.Image)
	}

	// Make sure the sha1 sum is correct
	h := sha1.New()
	f, err := os.Open(conf.Image)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(h, f); err != nil {
		t.Fatal(err)
	}
	sum := fmt.Sprintf("%x", h.Sum(nil))

	if conf.Sha1sum != sum {
		t.Errorf("CreateImage: expected sha1: %s got: %s", sum, conf.Sha1sum)
	}

	// extract image
	{
		// expect the image ova to contain only the following file names
		expectedNames := []string{
			"image.ovf",
			"image.mf",
			"image-disk1.vmdk",
		}

		imageDir := helpers.ExtractGzipArchive(t, conf.Image)
		list, err := ioutil.ReadDir(imageDir)
		if err != nil {
			t.Fatal(err)
		}

		var names []string
		infos := make(map[string]os.FileInfo)
		for _, fi := range list {
			names = append(names, fi.Name())
			infos[fi.Name()] = fi
		}

		if len(names) != 3 {
			t.Errorf("CreateImage: expected image (%s) to contain 3 files, found: %d - %s",
				imageDir, len(names), names)
		}
		for _, name := range expectedNames {
			if _, ok := infos[name]; !ok {
				t.Errorf("CreateImage: image (%s) is missing file: %s", names, name)
			}
		}

		// the vmx template should generate an ovf file that
		// does not contain an ethernet section.
		//
		ovf := filepath.Join(imageDir, "image.ovf")
		s, err := readFile(ovf)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(strings.ToLower(s), "ethernet") {
			t.Errorf("CreateImage: ovf contains 'ethernet' block:\n%s\n", s)
		}
	}
}

// this checks that CreateImage can take the relative path of a VMDK
func TestCreateImagePathResolution(t *testing.T) {
	const archive = "testdata/patch-test.tar.gz"

	dirname := helpers.ExtractGzipArchive(t, archive)
	defer os.RemoveAll(dirname)

	// get current working dir
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("could not get working dir")
	}

	if err := os.Chdir(dirname); err != nil {
		t.Errorf("Could not change to test tmp dir: %s", dirname)
	}

	conf := Config{stop: make(chan struct{})}

	if err := conf.CreateImage("expected.vmdk"); err != nil {
		t.Errorf("CreateImage couldn't expand absolute path of VMDK file: %s", err)
	}

	// change back to current working dir
	if err := os.Chdir(cwd); err != nil {
		t.Errorf("Could not change back to working dir: %s", cwd)
	}
}

func TestThatTheManifestIsGeneratedCorrectly(t *testing.T) {
	result := CreateManifest("1", "version", "sha1sum")
	expectedManifest := `---
name: bosh-vsphere-esxi-windows1-go_agent
version: 'version'
sha1: sha1sum
operating_system: windows1
cloud_properties:
  infrastructure: vsphere
  hypervisor: esxi
stemcell_formats:
- vsphere-ovf
- vsphere-ova
`
	if result != expectedManifest {
		t.Errorf("result:\n%s\ndoes not match expected\n%s\n", result, expectedManifest)
	}
}

func TestValidApplyPatchManifestFile(t *testing.T) {
	testCommand := fmt.Sprintf(
		"stembuild apply-patch %s",
		"testdata/valid-apply-patch.yml",
	)
	testArgs := strings.Split(testCommand, " ")
	os.Args = testArgs
	runInit()
	ParseFlags()

	errs := ValidateFlags()

	if len(errs) != 0 {
		t.Errorf("expected no errors, but got errors: %s", errs)
	}
}

func TestInvalidApplyPatchManifestFile(t *testing.T) {
	testCommand := fmt.Sprintf(
		"stembuild apply-patch %s",
		"testdata/invalid-apply-patch.yml",
	)
	testArgs := strings.Split(testCommand, " ")
	os.Args = testArgs
	runInit()
	ParseFlags()

	errs := ValidateFlags()

	if len(errs) != 1 {
		t.Error("expected single error; but got no, or more than one, error(s)")
	}
}
