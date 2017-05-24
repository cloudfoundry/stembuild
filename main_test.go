package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

var versionTests = []struct {
	s  string
	ok bool
}{
	{"1.2", true},
	{"-1.2", false},
	{"1.-2", false},
	{"001.002", true},
	{"0a1.002", false},
	{"1.a", false},
	{"a1.2", false},
	{"a.2", false},
	{"1.2 a", false},
}

func TestValidateVersion(t *testing.T) {
	for _, x := range versionTests {
		err := validateVersion(x.s)
		if (err == nil) != x.ok {
			if x.ok {
				t.Errorf("failed to validate version: %s\n", x.s)
			} else {
				t.Errorf("expected error for version (%s) but got: %v\n", x.s, err)
			}
		}
	}
}

func readdirnames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

func TestExtractOVA_Valid(t *testing.T) {
	const Count = 9
	const NameFmt = "file-%d"

	tmpdir, err := ioutil.TempDir("", "test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	if err := ExtractOVA("testdata/tar/valid.tar", tmpdir); err != nil {
		t.Fatal(err)
	}

	var expFileNames []string
	for i := 0; i <= Count; i++ {
		expFileNames = append(expFileNames, fmt.Sprintf("file-%d", i))
	}
	sort.Strings(expFileNames)

	names, err := readdirnames(tmpdir)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expFileNames, names) {
		t.Errorf("ExtractOVA: filenames want: %v got: %v", expFileNames, names)
	}

	// the content of each file is it's index
	// and a newline so 'file-2' contains "2\n"
	validFile := func(name string) error {
		path := filepath.Join(tmpdir, name)
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		var i int
		if _, err := fmt.Sscanf(name, NameFmt, &i); err != nil {
			return err
		}
		exp := fmt.Sprintf("%d\n", i)
		if s := string(b); s != exp {
			t.Errorf("ExtractOVA: file (%s) want: %s got: %s", name, exp, s)
		}
		return nil
	}

	for _, name := range names {
		if err := validFile(name); err != nil {
			t.Error(err)
		}
	}
}

func TestExtractOVA_Invalid(t *testing.T) {
	var tests = []struct {
		archive string
		reason  string
	}{
		{
			"has-sub-dir.tar",
			"subdirectories are not supported",
		},
		{
			"too-many-files.tar",
			"too many files read from archive (this is capped at 100)",
		},
		{
			"symlinks.tar",
			"symlinks are not supported",
		},
	}

	for _, x := range tests {
		tmpdir, err := ioutil.TempDir("", "test-")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpdir)

		filename := filepath.Join("testdata", "tar", x.archive)
		if err := ExtractOVA(filename, tmpdir); err == nil {
			t.Errorf("ExtractOVA (%s): expected error because:", x.archive, x.reason)
		}
	}
}

func TestApplyPatch(t *testing.T) {
	const archive = "testdata/patch-test.tar.gz"

	dirname := extractGzipArchive(t, archive)
	defer os.RemoveAll(dirname)

	vhd := filepath.Join(dirname, "original.vhd")
	delta := filepath.Join(dirname, "delta.patch")
	expVMDK := filepath.Join(dirname, "expected.vmdk")

	// output file
	newVMDK := filepath.Join(dirname, "image.vmdk")

	conf := Config{stop: make(chan struct{})}

	// Normal operation
	{
		if err := conf.ApplyPatch(vhd, delta, newVMDK); err != nil {
			t.Fatal(err)
		}
		expSrc, err := ioutil.ReadFile(expVMDK)
		if err != nil {
			t.Fatalf("ApplyPatch: reading expected vmdk file (%s): %s", expVMDK, err)
		}
		vmdkSrc, err := ioutil.ReadFile(newVMDK)
		if err != nil {
			t.Fatalf("ApplyPatch: reading vmdk file (%s): %s", newVMDK, err)
		}
		if !bytes.Equal(expSrc, vmdkSrc) {
			t.Fatalf("ApplyPatch: patched vmdk (%s) does not match expected vmdk (%s)",
				newVMDK, expVMDK)
		}
	}

	// Error when VMDK file already exists
	{
		if err := conf.ApplyPatch(vhd, delta, newVMDK); err == nil {
			t.Error("ApplyPatch: expected an error when the VMDK already exists got")
		}
	}
}

func TestCreateImage(t *testing.T) {
	const archive = "testdata/patch-test.tar.gz"

	dirname := extractGzipArchive(t, archive)
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
		t.Error("CreateImage: expected sha1: %s got: %s", sum, conf.Sha1sum)
	}

	// extract image
	{
		// expect the image ova to contain only the following file names
		expectedNames := []string{
			"image.ovf",
			"image.mf",
			"image-disk1.vmdk",
		}

		imageDir := extractGzipArchive(t, conf.Image)
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

		// the vmx template should generate an ovf file that does not
		// contain an ethernet section.
		//
		ovf := filepath.Join(imageDir, "image.ovf")
		s, err := readFile(ovf)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(strings.ToLower(s), "ethernet") {
			t.Errorf("CreateImage: ovf contains 'ethernet':\n%s\n", s)
		}
	}

}

func readFile(name string) (string, error) {
	b, err := ioutil.ReadFile(name)
	return string(b), err
}

// extractGzipArchive, extracts the tgz archive name to a temp directory
// returning the filepath of the temp directory.
func extractGzipArchive(t *testing.T, name string) string {
	t.Logf("extractGzipArchive: extracting tgz: %s", name)

	tmpdir, err := ioutil.TempDir("", "test-")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("extractGzipArchive: using temp directory: %s", tmpdir)

	f, err := os.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	w, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	if err := ExtractArchive(w, tmpdir); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return tmpdir
}
