package commandparser

func (p *ConstructCmd) GetWinRMUser() string {
	return p.winrmUsername
}

func (p *ConstructCmd) GetWinRMPwd() string {
	return p.winrmPassword
}

func (p *ConstructCmd) GetStemcellVersion() string {
	return p.stemcellVersion
}

func (p *ConstructCmd) GetWinRMIp() string {
	return p.winrmIP
}
