package patch

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type ApplyPatch struct {
	PatchFile string `yaml:"patch_file"`
	OSVersion string `yaml:"os_version"`
	OutputDir string `yaml:"output_dir"`
	Version   string `yaml:"version"`
	VHDFile   string `yaml:"vhd_file"`
	VMDKFile  string `yaml:"vmdk_file"`
}

// Copy into `d` the values in `s` which are empty in `d`.
func (d *ApplyPatch) CopyInto(s ApplyPatch) {
	if d.PatchFile == "" {
		d.PatchFile = s.PatchFile
	}

	if d.OSVersion == "" {
		d.OSVersion = s.OSVersion
	}

	if d.OutputDir == "" {
		d.OutputDir = s.OutputDir
	}

	if d.Version == "" {
		d.Version = s.Version
	}

	if d.VHDFile == "" {
		d.VHDFile = s.VHDFile
	}

	if d.VMDKFile == "" {
		d.VMDKFile = s.VMDKFile
	}
}

func LoadPatchManifest(fileName string, patchArgs *ApplyPatch) error {
	_, err := os.Stat(fileName)
	if err != nil {
		return err
	}

	var patchManifest ApplyPatch
	manifestFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(manifestFile, &patchManifest)
	if err != nil {
		return err
	}

	patchArgs.CopyInto(patchManifest)

	return nil
}
