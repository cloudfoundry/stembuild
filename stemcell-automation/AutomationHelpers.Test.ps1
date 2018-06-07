. ./AutomationHelpers.ps1

Describe "CopyPSModules" {
    It "can copy PS Modules to target directory" {
        Mock Write-Log {}
        Mock Expand-Archive {}

        CopyPSModules

        Assert-MockCalled Expand-Archive -Times 1 -ParameterFilter { $LiteralPath -eq ".\bosh-psmodules.zip" -and $DestinationPath -eq "C:\Program Files\WindowsPowerShell\Modules\" }
        Assert-MockCalled Write-Log -Times 1 -ParameterFilter { $Message -eq "Succesfully migrated Bosh Powershell modules to destination dir" }
    }

    It "fails gracefully when expanding archive fails" {
        Mock Expand-Archive { throw "Expand-Archive failed because something went wrong" }
        Mock Write-Log { }

        { CopyPSModules } | Should -Throw "Expand-Archive failed because something went wrong"

        Assert-MockCalled Write-Log -Times 1 -ParameterFilter { $Message -eq "Expand-Archive failed because something went wrong" }
        Assert-MockCalled Write-Log -Times 1 -ParameterFilter { $Message -eq "Failed to copy Bosh Powershell Modules into destination dir. See 'c:\provisions\log.log' for mor info." }
    }
}

Describe "InstallCFFeatures" {
    It "invokes the Install-CFFeatures powershell cmdlet" {
        Mock Install-CFFeatures { }
        Mock Write-Log {}

        InstallCFFeatures

        Assert-MockCalled Install-CFFeatures -Times 1
        Assert-MockCalled Write-Log -Times 1 -ParameterFilter { $Message -eq "Successfully installed CF features" }
    }

    It "fails gracefully when installing CF Features" {
        Mock Install-CFFeatures { throw "Something terrible happened while attempting to install a CF feature" }
        Mock Write-Log

        { InstallCFFeatures } | Should -Throw "Something terrible happened while attempting to install a CF feature"

        Assert-MockCalled Write-Log -Times 1 -ParameterFilter { $Message -eq "Something terrible happened while attempting to install a CF feature"}
        Assert-MockCalled Write-Log -Times 1 -ParameterFilter { $Message -eq "Failed to install the CF features. See 'c:\provisions\log.log' for mor info."}
    }
}

Describe "InstallCFCell" {
    It "execute the Protect-CFCell powershell cmdlet" {
        Mock Protect-CFCell { }
        Mock Write-Log

        InstallCFCell

        Assert-MockCalled Protect-CFCell -Times 1
        Assert-MockCalled Write-Log -Times 1 -ParameterFilter { $Message -eq "Succesfully ran Protect-CFCell" }
    }

    It "fails gracefully when Protect-CFCell powershell cmdlet fails" {
        Mock Protect-CFCell { throw "Something terrible happened while attempting to execute Protect-CFCell" }
        Mock Write-Log

        { InstallCFCell } | Should -Throw "Something terrible happened while attempting to execute Protect-CFCell"

        Assert-MockCalled Write-Log -Times 1 -ParameterFilter { $Message -eq "Something terrible happened while attempting to execute Protect-CFCell"}
        Assert-MockCalled Write-Log -Times 1 -ParameterFilter { $Message -eq "Failed to execute Protect-CFCell powershell cmdlet. See 'c:\provisions\log.log' for mor info."}
    }

}