param(
    [string]$Organization = "",
    [string]$Owner = "",
    [switch]$SkipRandomPassword
)

Push-Location $PSScriptRoot

. ./AutomationHelpers.ps1

function Write-Log
{
    Param (
        [Parameter(Mandatory = $True, Position = 1)][string]$Message,
        [string]$LogFile = "C:\provision\log.log"
    )

    $LogDir = (split-path $LogFile -parent)
    If ((Test-Path $LogDir) -ne $True)
    {
        New-Item -Path $LogDir -ItemType Directory -Force
    }

    $msg = "{0} {1}" -f (Get-Date -Format o), $Message
    Add-Content -Path $LogFile -Value $msg -Encoding 'UTF8'
    Write-Host $msg
}

try {
    Validate-OSVersion
    Check-Dependencies
    Set-RegKeys

    # create the scheduled task to run second script here!
    $Sta = Create-VMPrepTaskAction -Organization $Organization -Owner $Owner
    if($SkipRandomPassword) {
        $Sta = Create-VMPrepTaskAction -Organization $Organization -Owner $Owner -SkipRandomPassword
    }
    $STPrin = New-ScheduledTaskPrincipal -UserID "NT AUTHORITY\SYSTEM" -LogonType ServiceAccount -RunLevel Highest
    $Stt = New-ScheduledTaskTrigger -AtStartup
    Register-ScheduledTask BoshCompleteVMPrep -Action $Sta -Trigger $Stt -Principal $STPrin -Description "Bosh Stemcell Automation task to complete the vm preparation"
    Write-Log "Successfully registered the Bosh Stemcell Automation scheduled task"

    CopyPSModules
    InstallBoshAgent
    InstallOpenSSH
    Enable-SSHD
    InstallCFFeatures
} catch [Exception] {
    Write-Log "Failed to install Bosh dependencies. See 'c:\provisions\log.log' for more info."
    DeleteScheduledTask
    Exit 1
}

Pop-Location
