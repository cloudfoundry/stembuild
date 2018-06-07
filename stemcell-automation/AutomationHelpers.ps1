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
    try
    {
        Protect-CFCell
        Write-Log "Succesfully ran Protect-CFCell"
    } catch [ Exception ] {
        Write-Log $_.Exception.Message
        Write-Log "Failed to execute Protect-CFCell powershell cmdlet. See 'c:\provisions\log.log' for mor info."
        throw $_.Exception
    }
}
