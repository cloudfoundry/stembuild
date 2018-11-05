package vmxtemplate_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/pivotal-cf-experimental/stembuild/vmxtemplate"
)

func parseVMX(vmx string) (map[string]string, error) {
	m := make(map[string]string)
	for _, s := range strings.Split(vmx, "\n") {
		if s == "" {
			continue
		}
		n := strings.IndexByte(s, '=')
		if n == -1 {
			return nil, fmt.Errorf("parse vmx: invalid line: %s", s)
		}
		k := strings.TrimSpace(s[:n])
		v, err := strconv.Unquote(strings.TrimSpace(s[n+1:]))
		if err != nil {
			return nil, err
		}
		if _, ok := m[k]; ok {
			return nil, fmt.Errorf("parse vmx: duplicate key: %s", k)
		}
		m[k] = v
	}
	if len(m) == 0 {
		return nil, errors.New("parse vmx: empty vmx")
	}
	return m, nil
}

func checkVMXTemplate(t *testing.T, hwVersion int, vmdkPath, vmxContent string) {
	const vmdkPathKeyName = "scsi0:0.fileName"
	const hwVersionKeyName = "virtualHW.version"

	m, err := parseVMX(vmxContent)
	if err != nil {
		t.Fatal(err)
	}
	if s := m[vmdkPathKeyName]; s != vmdkPath {
		t.Errorf("VMXTemplate: key: %q want: %q got: %q", vmdkPathKeyName, vmdkPath, s)
	}

	expectedHWVersion := strconv.Itoa(hwVersion)
	if s := m[hwVersionKeyName]; s != expectedHWVersion {
		t.Errorf("VMXTemplate: key: %q want: %q got: %q", hwVersionKeyName, expectedHWVersion, s)
	}
}

const vmdkPath = "FooBarBaz.vmdk"
const virtualHWVersion = 60

func TestVMXTemplate(t *testing.T) {
	var buf bytes.Buffer
	if err := vmxtemplate.VMXTemplate(vmdkPath, virtualHWVersion, &buf); err != nil {
		t.Fatal(err)
	}
	checkVMXTemplate(t, virtualHWVersion, vmdkPath, buf.String())

	if err := vmxtemplate.VMXTemplate("", 0, &buf); err == nil {
		t.Error("VMXTemplate: expected error for empty vmx filename")
	}
}

func TestWriteVMXTemplate(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test-")
	if err != nil {
		t.Fatal(err)
	}
	vmxPath := filepath.Join(tmpdir, "FooBarBaz.vmx")

	if err := vmxtemplate.WriteVMXTemplate(vmdkPath, virtualHWVersion, vmxPath); err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadFile(vmxPath)
	if err != nil {
		t.Fatal(err)
	}
	checkVMXTemplate(t, virtualHWVersion, vmdkPath, string(b))

	if err := os.Remove(vmxPath); err != nil {
		t.Fatal(err)
	}

	// vmx file is deleted if there is an error
	if err := vmxtemplate.WriteVMXTemplate("", 0, vmxPath); err == nil {
		t.Error("WriteVMXTemplate: expected error for empty vmx filename")
	}
	if _, err := os.Stat(vmxPath); err == nil {
		t.Errorf("WriteVMXTemplate: failed to delete vmx file on error: %s", vmxPath)
	}
}
