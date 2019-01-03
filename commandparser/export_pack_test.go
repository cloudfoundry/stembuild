package commandparser

import "github.com/cloudfoundry-incubator/stembuild/filesystem"

func (p *PackageCmd) GetVMDK() string {
	return p.vmdk
}

func (p *PackageCmd) GetOS() string {
	return p.os
}

func (p *PackageCmd) GetStemcellVersion() string {
	return p.stemcellVersion
}

func (p *PackageCmd) GetOutputDir() string {
	return p.outputDir
}

func (p *PackageCmd) ValidateFreeSpaceForPackage(fs filesystem.FileSystem) (bool, uint64, error) {
	return p.validateFreeSpaceForPackage(fs)
}
func (p *PackageCmd) SetVMDK(vmdk string) {
	p.vmdk = vmdk
}
