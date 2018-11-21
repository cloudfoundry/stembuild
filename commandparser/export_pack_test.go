package commandparser

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
