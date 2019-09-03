function CopyPSModules
{
    try
    {
        Expand-Archive -LiteralPath ".\bosh-psmodules.zip" -DestinationPath "C:\Program Files\WindowsPowerShell\Modules\" -Force
        Write-Log "Succesfully migrated Bosh Powershell modules to destination dir"
    }
    catch [Exception]
    {
        Write-Log $_.Exception.Message
        Write-Log "Failed to copy Bosh Powershell Modules into destination dir. See 'c:\provisions\log.log' for mor info."
        throw $_.Exception
    }
}

function InstallCFFeatures
{
    try
    {
        Install-CFFeatures2016
        Write-Log "Successfully installed CF features"
        Restart-Computer
    }
    catch [Exception]
    {
        Write-Log $_.Exception.Message
        #TODO: Fix spelling!
        Write-Log "Failed to install the CF features. See 'c:\provisions\log.log' for mor info."
        throw $_.Exception
    }
}

function InstallCFCell
{
    try
    {
        Protect-CFCell
        Write-Log "Succesfully ran Protect-CFCell"
    }
    catch [Exception]
    {
        Write-Log $_.Exception.Message
        Write-Log "Failed to execute Protect-CFCell powershell cmdlet. See 'c:\provisions\log.log' for mor info."
        throw $_.Exception
    }
}

function InstallBoshAgent
{
    try
    {
        Install-Agent -Iaas "vsphere" -agentZipPath ".\agent.zip"
        Write-Log "Bosh agent successfully installed"
    }
    catch [Exception]
    {
        Write-Log $_.Exception.Message
        Write-Log "Failed to execute Install-Agent powershell cmdlet. See 'c:\provisions\log.log' for mor info."
        throw $_.Exception
    }
}

function InstallOpenSSH
{
    try
    {
        Install-SSHD -SSHZipFile ".\OpenSSH-Win64.zip"
        Write-Log "OpenSSH successfully installed"
    }
    catch [Exception]
    {
        Write-Log $_.Exception.Message
        Write-Log "Failed to execute Install-SSHD powershell cmdlet. See 'c:\provisions\log.log' for mor info."
        throw $_.Exception
    }
}

#This function contains all our registry changes related to fixing the zombieload and meltdown bugs
function Set-RegKeys
{
    $PathExists = Test-Path 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\QualityCompat'
    try
    {
        if ($PathExists -eq $False)
        {
            New-Item -Path 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\QualityCompat'
        }

        New-ItemProperty -Path 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\QualityCompat' -Value 0 -Name 'cadca5fe-87d3-4b96-b7fb-a231484277cc' -force
        New-ItemProperty -Path 'HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Virtualization' -Value 1.0 -Name 'MinVmVersionForCpuBasedMitigations' -force
        New-ItemProperty -Path 'HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management' -Value 3 -Name 'FeatureSettingsOverrideMask' -force
        New-ItemProperty -Path 'HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management' -Value 72 -Name 'FeatureSettingsOverride' -force
        Write-Log "Meltdown/zombieload registry keys successfully added"
    }
    catch [Exception]
    {
        Write-Log $_.Exception.Message
        Write-Log "Failed to set meltdown/zombieload registry keys. See 'c:\provisions\log.log' for mor info."
        throw $_.Exception
    }
}

function CleanUpVM
{
    try
    {
        Optimize-Disk
        Compress-Disk
        Write-Log "Successfully cleaned up the VM's disk"
    }
    catch [Exception]
    {
        Write-Log $_.Exception.Message
        Write-Log "Failed to clean up the VM's disk. See 'c:\provisions\log.log' for mor info."
        throw $_.Exception
    }
}

function Is-Special()
{
    param ([parameter(Mandatory = $true)] [string]$c)

    return $c -cmatch '[!-/:-@[-`{-~]'
}

function Valid-Password()
{
    param ([parameter(Mandatory = $true)] [string]$Password)

    $digits = 0
    $special = 0
    $alphaLow = 0
    $alphaHigh = 0

    if ($Password.Length -lt 8)
    {
        return $false
    }

    $tmp = $Password.ToCharArray()

    foreach ($c in $Password.ToCharArray())
    {
        if ($c -cmatch '\d')
        {
            $digits = 1
        }
        elseif ($c -cmatch '[a-z]')
        {
            $alphaLow = 1
        }
        elseif ($c -cmatch '[A-Z]')
        {
            $alphaHigh = 1
        }
        elseif (Is-Special $c)
        {
            $special = 1
        }
        else
        {
            #Invalid char
            return $false
        }
    }
    return ($digits + $special + $alphaLow + $alphaHigh) -ge 3
}

function GenerateRandomPassword
{
    $CharList = "!`"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_``abcdefghijklmnopqrstuvwxyz{|}~".ToCharArray()
    $limit = 200
    $count = 0

    while ($limit-- -gt 0)
    {
        $passwd = (Get-Random -InputObject $CharList -Count 24) -join ''
        if (Valid-Password -Password $passwd)
        {
            Write-Log "Successfully generated password"
            return $passwd
        }
    }
    Write-Log "Failed to generate password after 200 attempts"
    throw "Unable to generate a valid password after 200 attempts"
}

function SysprepVM
{
    Param (
        [string]$Organization = "",
        [string]$Owner = "",
        [bool]$SkipRandomPassword = $false
    )

    try
    {
        Expand-Archive -LiteralPath ".\LGPO.zip" -DestinationPath "C:\Windows\"
        Write-Log "Successfully migrated LGPO to destination dir"

        if ($SkipRandomPassword) {
            Invoke-Sysprep -IaaS "vsphere" -Organization $Organization -Owner $Owner
        }else {
            $randomPassword = GenerateRandomPassword
            Invoke-Sysprep -IaaS "vsphere" -NewPassword $randomPassword -Organization $Organization -Owner $Owner
        }
    }
    catch [Exception]
    {
        Write-Log $_.Exception.Message
        Write-Log "Failed to Sysprep the VM's. See 'c:\provisions\log.log' for mor info."
        throw $_.Exception
    }
}

function Check-Dependencies
{
    try
    {
        $depsObj = (Get-Content -Path "$PSScriptRoot/deps.json") -join '' | ConvertFrom-Json
        if ($depsObj.psobject.properties.Count -eq 0 -or $depsObj.psobject.properties.Count -eq $null)
        {
            throw "Dependency file is empty"
        }

        $corruptedOrMissingFile = $false
        $depsObj.psobject.properties | ForEach {
            $fileName = $_.Name
            $expectedFileHash = $_.Value.sha
            if (Test-Path -Path "$PSScriptRoot/$fileName")
            {
                $fileHash = Get-FileHash -Path "$PSScriptRoot/$fileName"
                if ($fileHash.Hash -notmatch $expectedFileHash)
                {
                    Write-Log "$PSScriptRoot/$fileName does not have the correct hash"
                    $corruptedOrMissingFile = $true
                }
            }
            else
            {
                Write-Log "$PSScriptRoot/$fileName is required but was not found"
                $corruptedOrMissingFile = $true
            }
        }

        if ($corruptedOrMissingFile)
        {
            throw "One or more files are corrupted or missing."
        }

    }
    catch [Exception]
    {
        Write-Log $_.Exception.Message
        Write-Log "Failed to validate required dependencies. See 'c:\provisions\log.log' for more info."
        throw $_.Exception
    }

    Write-Log "Found all dependencies"
}

function Get-OSVersionString
{
    return [System.Environment]::OSVersion.Version.ToString()
}

function Validate-OSVersion
{
    try
    {
        $osVersion = Get-OSVersionString
        if ($osVersion -match "10\.0\.16299\..+")
        {
            Write-Log "Found correct OS version: Windows Server 2016, Version 1709"
        }
        elseif ($osVersion -match "10\.0\.17134\..+")
        {
            Write-Log "Found correct OS version: Windows Server 2016, Version 1803"
        }
        elseif ($osVersion -match "10\.0\.17763\..+")
        {
            Write-Log "Found correct OS version: Windows Server 2019"
        }
        else {
            throw "OS Version Mismatch: Please use Windows Server 2019 or Windows Server 2016, Version 1709 or 1803"
        }
    }
    catch [Exception]
    {
        Write-Log $_.Exception.Message
        Write-Log "Failed to validate the OS version. See 'c:\provisions\log.log' for more info."
        throw $_.Exception
    }
}

function DeleteScheduledTask {
    try {
        if ((Get-ScheduledTask | ForEach { $_.TaskName }) -ccontains "BoshCompleteVMPrep") {
            Unregister-ScheduledTask -TaskName BoshCompleteVMPrep -Confirm:$false
            Write-Log "Successfully deleted the 'BoshCompleteVMPrep' scheduled task"
        }
        else {
            Write-Log "BoshCompleteVMPrep schedule task was not registered"
        }
    }
    catch [Exception] {
        Write-Log $_.Exception.Message
        Write-Log "Failed to unregister the BoshCompleteVMPrep scheduled task. See 'c:\provisions\log.log' for more info."
        throw $_.Exception
    }
}

function Create-VMPrepTaskAction {
    param(
        [string]$Organization="",
        [string]$Owner="",
        [switch]$SkipRandomPassword
    )

    $arguments = "-NoExit -File ""$PSScriptRoot\Complete-VMPrep.ps1"""
    if ($Organization -ne "") {
        $arguments += " -Organization ""$Organization"""
    }

    if ($Owner -ne "") {
        $arguments += " -Owner ""$Owner"""
    }

    if ($SkipRandomPassword) {
        $arguments += " -SkipRandomPassword"
    }

    New-ScheduledTaskAction -Execute "powershell.exe" -Argument $arguments
}

function Remove-SSHKeys
{
    $SSHDir = "C:\Program Files\OpenSSH"

    Push-Location $SSHDir
    New-Item -ItemType Directory -Path "$env:ProgramData\ssh" -ErrorAction Ignore

    "Removing any existing host keys"
    Remove-Item -Path "$env:ProgramData\ssh\ssh_host_*" -ErrorAction Ignore
    Pop-Location
}

function Run-LGPO
{
    param (
        [string]$LGPOPath = $( Throw "Provide LGPO path" ),
        [string]$InfFilePath = $( Throw "Provide Inf file path" )
    )
    & $LGPOPath /s $InfFilePath
}

function Enable-SSHD
{
    if ((Get-NetFirewallRule | where { $_.DisplayName -ieq 'SSH' }) -eq $null)
    {
        "Creating firewall rule for SSH"
        New-NetFirewallRule -Protocol TCP -LocalPort 22 -Direction Inbound -Action Allow -DisplayName SSH
    }
    else
    {
        "Firewall rule for SSH already exists"
    }

    $InfFilePath = "$env:WINDIR\Temp\enable-ssh.inf"

    $InfFileContents = @'
[Unicode]
Unicode=yes
[Version]
signature=$CHICAGO$
Revision=1:w
[Registry Values]
[System Access]
[Privilege Rights]
SeDenyNetworkLogonRight=*S-1-5-32-546
SeAssignPrimaryTokenPrivilege=*S-1-5-19,*S-1-5-20,*S-1-5-80-3847866527-469524349-687026318-516638107-1125189541
'@
    $LGPOPath = "$env:WINDIR\LGPO.exe"
    if (Test-Path $LGPOPath)
    {
        Out-File -FilePath $InfFilePath -Encoding unicode -InputObject $InfFileContents -Force
        Try
        {
            Run-LGPO -LGPOPath $LGPOPath -InfFilePath $InfFilePath
        }
        Catch
        {
            throw "LGPO.exe failed with: $_.Exception.Message"
        }
    }
    else
    {
        "Did not find $LGPOPath. Assuming existing security policies are sufficient to support ssh."
    }

    Set-Service -Name sshd -StartupType Automatic
    # ssh-agent is not the same as ssh-agent in *nix openssh
    Set-Service -Name ssh-agent -StartupType Automatic

    Remove-SSHKeys
}

function Install-SecurityPoliciesAndRegistries
{
    try
    {
        $osVersion2019Regex = "10\.0\.17763\..+"
        $osVersion = Get-OSVersionString
        Write-Log "osVersion: $osVersion"
        if ($osVersion -match $osVersion2019Regex) {
            Set-InternetExplorerRegistries
            Write-Log "Succesfully ran Set-InternetExplorerRegistries"
        } else {
            Write-Log "Did not run Set-InternetExplorerRegistries because OS version was not 2019"
        }
    }
    catch [Exception]
    {
        Write-Log $_.Exception.Message
        Write-Log "Failed to execute Set-InternetExplorerRegistries powershell cmdlet. See 'c:\provisions\log.log' for mor info."
        throw $_.Exception
    }
}


