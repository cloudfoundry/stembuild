package config

import "errors"

type SourceConfig struct {
	Vmdk     string
	VmName   string
	Username string
	Password string
	URL      string
}

type Source int

const (
	VMDK Source = iota
	VCENTER
	NIL
)

//Three package_parameters:
//1. If VMDK provided and VMName not provided, ignore other VM credentials (simplifies a lot)
//2. Use VCenterCredentials struct
//3. use below implementation:

func (c SourceConfig) GetSource() (Source, error) {
	if c.vmdkProvided() && c.partialvCenterProvided() {
		return NIL, errors.New("configuration provided for VMDK & vCenter sources")
	}

	if c.vmdkProvided() {
		return VMDK, nil
	}

	if c.vcenterProvided() {
		return VCENTER, nil
	}

	if c.partialvCenterProvided() {
		return NIL, errors.New("missing vCenter configurations")
	}

	return NIL, errors.New("no configuration was provided")
}

func (c SourceConfig) vmdkProvided() bool {
	return c.Vmdk != ""
}

func (c SourceConfig) vcenterProvided() bool {
	if c.VmName != "" && c.Username != "" && c.Password != "" && c.URL != "" {
		return true
	}
	return false
}

//At least one vCenter configuration given
func (c SourceConfig) partialvCenterProvided() bool {
	if c.VmName != "" || c.Username != "" || c.Password != "" || c.URL != "" {
		return true
	}
	return false
}
