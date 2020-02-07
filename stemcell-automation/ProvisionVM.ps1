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
    try
    {
        Install-WUCerts
    }
    catch [Exception]
    {
        Write-Log $_.Exception.Message
        Write-Warning "This should not impact the successful execution of stembuild construct. If the root certificates are out of date, Diego cells running on VMs built off of this stemcell may not be able to make outbound network connections."
    }
    Create-VersionFile -Version $Version
    Restart-Computer
}
