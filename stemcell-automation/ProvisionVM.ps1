. ./AutomationHelpers.ps1

function ProvisionVM() {
    param (
        [string]$Version,
        [switch]$FailOnInstallWUCerts
    )

    RunQuickerDism -IgnoreErrors $True
    CopyPSModules
    Set-RegKeys
    InstallBoshAgent
    InstallOpenSSH
    Extract-LGPO
    Install-SecurityPoliciesAndRegistries
    Enable-SSHD
    InstallCFFeatures

    try
    {
        Install-WUCerts
    }
    catch [Exception]
    {
        Write-Log $_.Exception.Message

        if ($FailOnInstallWUCerts) {
            throw $_.Exception
        } else {
            Write-Warning "Failed to retrieve updated root certificates from the public Windows Update Server. This should not impact the successful execution of stembuild construct. If your root certificates are out of date, Diego cells running on VMs built from this stemcell may not be able to make outbound network connections."
        }
    }
    Create-VersionFile -Version $Version
    RunQuickerDism -IgnoreErrors $True
    Restart-Computer
}
