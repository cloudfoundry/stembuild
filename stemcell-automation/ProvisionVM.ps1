. ./AutomationHelpers.ps1

function ProvisionVM() {
    CopyPSModules
    InstallBoshAgent
    InstallOpenSSH
    Enable-SSHD
    InstallCFFeatures
    Enable-HyperV
}
