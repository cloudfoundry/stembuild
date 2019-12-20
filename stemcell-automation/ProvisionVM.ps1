. ./AutomationHelpers.ps1

function ProvisionVM() {
    param (
        [string]$Version
    )

    CopyPSModules
    InstallBoshAgent
    InstallOpenSSH
    Extract-LGPO
    Install-SecurityPoliciesAndRegistries
    Enable-SSHD
    InstallCFFeatures
    Enable-HyperV
    Install-WUCerts
    Create-VersionFile -Version $Version
    Restart-Computer
}
