. ./AutomationHelpers.ps1

function ProvisionVM() {
    CopyPSModules
    InstallBoshAgent
    InstallOpenSSH
    Extract-LGPO
    Enable-SSHD
    InstallCFFeatures
    Enable-HyperV
}
