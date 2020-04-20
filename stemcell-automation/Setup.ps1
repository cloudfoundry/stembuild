param(
    [string]$Organization = "",
    [string]$Owner = "",
    [switch]$SkipRandomPassword,
    [String]$Version
)

Push-Location $PSScriptRoot

. ./AutomationHelpers.ps1
. ./ProvisionVM.ps1

try
{
    Setup -Organization $Organization -Owner $Owner -SkipRandomPassword:$SkipRandomPassword -Version $Version
}
catch [Exception]
{
    Write-Log "Failed to install Bosh dependencies. See 'c:\provision\log.log' for more info."
    DeleteScheduledTask
    Exit 1
}

Pop-Location
