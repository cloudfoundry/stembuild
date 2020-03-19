param(
    [String]$Version
)

Push-Location $PSScriptRoot

. ./AutomationHelpers.ps1
. ./ProvisionVM.ps1

try
{
    Setup -Version $Version
}
catch [Exception]
{
    Write-Log "Failed to install Bosh dependencies. See 'c:\provision\log.log' for more info."
    Exit 1
}

Pop-Location
