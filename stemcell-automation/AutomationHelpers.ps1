function CopyPSModules {
    try {
        Expand-Archive -LiteralPath ".\bosh-psmodules.zip" -DestinationPath "C:\Program Files\WindowsPowerShell\Modules\" -Force
        Write-Log "Succesfully migrated Bosh Powershell modules to destination dir"
    } catch [ Exception ] {
        Write-Log $_.Exception.Message
        Write-Log "Failed to copy Bosh Powershell Modules into destination dir. See 'c:\provisions\log.log' for mor info."
        throw $_.Exception
    }
}

function InstallCFFeatures {
    try {
        Install-CFFeatures
        Write-Log "Successfully installed CF features"
    } catch [ Exception ] {
        Write-Log $_.Exception.Message
        Write-Log "Failed to install the CF features. See 'c:\provisions\log.log' for mor info."
        throw $_.Exception
    }
}

function InstallCFCell {
    try {
        Protect-CFCell
        Write-Log "Succesfully ran Protect-CFCell"
    } catch [ Exception ] {
        Write-Log $_.Exception.Message
        Write-Log "Failed to execute Protect-CFCell powershell cmdlet. See 'c:\provisions\log.log' for mor info."
        throw $_.Exception
    }
}

function InstallBoshAgent {
    try {
        Install-Agent -Iaas "vsphere" -agentZipPath ".\agent.zip"
        Write-Log "Bosh agent successfully installed"
    } catch [ Exception ] {
        Write-Log $_.Exception.Message
        Write-Log "Failed to execute Install-Agent powershell cmdlet. See 'c:\provisions\log.log' for mor info."
        throw $_.Exception
    }
}

function InstallOpenSSH {
    try {
        Install-SSHD -SSHZipFile ".\OpenSSH-Win64.zip"
        Write-Log "OpenSSH successfully installed"
    } catch [ Exception ] {
        Write-Log $_.Exception.Message
        Write-Log "Failed to execute Install-SSHD powershell cmdlet. See 'c:\provisions\log.log' for mor info."
        throw $_.Exception
    }
}

function CleanUpVM {
    try {
        Optimize-Disk
        Compress-Disk
        Write-Log "Successfully cleaned up the VM's disk"
    } catch [ Exception ] {
        Write-Log $_.Exception.Message
        Write-Log "Failed to clean up the VM's disk. See 'c:\provisions\log.log' for mor info."
        throw $_.Exception
    }
}

function Is-Special() {
    param([parameter(Mandatory=$true)] [string]$c)

    return $c -cmatch '[!-/:-@[-`{-~]'
}

function Valid-Password() {
    param([parameter(Mandatory=$true)] [string]$Password)

    $digits = 0
    $special = 0
    $alphaLow = 0
    $alphaHigh = 0

    if ($Password.Length -lt 8) {
        return $false
    }

    $tmp = $Password.ToCharArray()

    foreach ($c in $Password.ToCharArray()) {
        if ($c -cmatch '\d') {
            $digits = 1
        } elseif ($c -cmatch '[a-z]') {
            $alphaLow = 1
        } elseif ($c -cmatch '[A-Z]') {
            $alphaHigh = 1
        } elseif (Is-Special $c) {
            $special = 1
        } else {
            #Invalid char
            return $false
        }
    }
    return ($digits + $special + $alphaLow + $alphaHigh) -ge 3
}

function GenerateRandomPassword {
    $CharList = "!`"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_``abcdefghijklmnopqrstuvwxyz{|}~".ToCharArray()
    $limit = 200
    $count = 0

    while ($limit-- -gt 0) {
        $passwd = (Get-Random -InputObject $CharList -Count 24) -join ''
        if (Valid-Password -Password $passwd) {
            Write-Log "Successfully generated password"
            return $passwd
        }
    }
    Write-Log "Failed to generate password after 200 attempts"
    throw "Unable to generate a valid password after 200 attempts"
}

function SysprepVM {
    try {
        Expand-Archive -LiteralPath ".\LGPO.zip" -DestinationPath "C:\Windows\"
        Write-Log "Successfully migrated LGPO to destination dir"

        $randomPassword = GenerateRandomPassword

        Invoke-Sysprep -IaaS "vsphere" -NewPassword $randomPassword
    } catch [ Exception ] {
        Write-Log $_.Exception.Message
        Write-Log "Failed to Sysprep the VM's. See 'c:\provisions\log.log' for mor info."
        throw $_.Exception
    }
}

function Check-Dependencies{
    try
    {
        $depsObj = ( Get-Content -Path "$PSScriptRoot/deps.json") -join '`n' | ConvertFrom-Json
        if ($depsObj.psobject.properties.Count -eq 0 -or $depsObj.psobject.properties.Count -eq $null)  {
            throw "Dependency file is empty"
        }

        $hashtable = @{ }
        $depsObj.psobject.properties | Foreach { $hashtable[$_.Name] = $_.Value }

        $corruptedOrMissingFile = $false
        foreach ($item in $hashtable.GetEnumerator())
        {
            $fileName = $item.Key
            $expectedFileHash = $item.Value
            if (Test-Path -Path "$PSScriptRoot/$fileName")
            {
                $fileHash = Get-FileHash -Path "$PSScriptRoot/$fileName"
                if ($fileHash.Hash -notmatch $expectedFileHash)
                {
                    Write-Log "$PSScriptRoot/$fileName does not have the correct hash"
                    $corruptedOrMissingFile = $true
                }
            }
            else {
                Write-Log "$PSScriptRoot/$fileName is required but was not found"
                $corruptedOrMissingFile = $true
            }
        }

        if ($corruptedOrMissingFile) {
            throw "One or more files are corrupted or missing."
        }

    } catch [Exception] {
        Write-Log $_.Exception.Message
        Write-Log "Failed to validate required dependencies. See 'c:\provisions\log.log' for more info."
        throw $_.Exception
    }

    Write-Log "Found all dependencies"
}