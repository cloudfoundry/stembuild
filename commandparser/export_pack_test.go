package commandparser

import "github.com/cloudfoundry-incubator/stembuild/filesystem"

func (p *PackageCmd) GetVMDK() string {
	return p.sourceConfig.Vmdk
}

func (p *PackageCmd) GetOS() string {
	return p.outputConfig.Os
}

func (p *PackageCmd) GetStemcellVersion() string {
	return p.outputConfig.StemcellVersion
}

func (p *PackageCmd) GetOutputDir() string {
	return p.outputConfig.OutputDir
}

func (p *PackageCmd) ValidateFreeSpaceForPackage(fs filesystem.FileSystem) error {
	return p.ValidateFreeSpaceForPackage(fs)
}
func (p *PackageCmd) SetVMDK(vmdk string) {
	p.sourceConfig.Vmdk = vmdk
}
