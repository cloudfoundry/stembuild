. ./AutomationHelpers.ps1

Describe "CopyPSModules" {
    It "can copy PS Modules to target directory" {
        Mock Write-Log {}
        Mock Expand-Archive {}

        { CopyPSModules } | Should -Not -Throw

        Assert-MockCalled Expand-Archive -Times 1 -Scope It -ParameterFilter { $LiteralPath -eq ".\bosh-psmodules.zip" -and $DestinationPath -eq "C:\Program Files\WindowsPowerShell\Modules\" }
        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Succesfully migrated Bosh Powershell modules to destination dir" }
    }

    It "fails gracefully when expanding archive fails" {
        Mock Expand-Archive { throw "Expand-Archive failed because something went wrong" }
        Mock Write-Log {}

        { CopyPSModules } | Should -Throw "Expand-Archive failed because something went wrong"

        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Expand-Archive failed because something went wrong" }
        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Failed to copy Bosh Powershell Modules into destination dir. See 'c:\provisions\log.log' for mor info." }
    }
}

Describe "InstallCFFeatures" {
    It "executes the Install-CFFeatures powershell cmdlet" {
        Mock Install-CFFeatures {}
        Mock Write-Log {}

        { InstallCFFeatures } | Should -Not -Throw

        Assert-MockCalled Install-CFFeatures -Times 1 -Scope It
        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Successfully installed CF features" }
    }

    It "fails gracefully when installing CF Features" {
        Mock Install-CFFeatures { throw "Something terrible happened while attempting to install a CF feature" }
        Mock Write-Log {}

        { InstallCFFeatures } | Should -Throw "Something terrible happened while attempting to install a CF feature"

        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Something terrible happened while attempting to install a CF feature"}
        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Failed to install the CF features. See 'c:\provisions\log.log' for mor info."}
    }
}

Describe "InstallCFCell" {
    It "executes the Protect-CFCell powershell cmdlet" {
        Mock Protect-CFCell {}
        Mock Write-Log {}

        { InstallCFCell } | Should -Not -Throw

        Assert-MockCalled Protect-CFCell -Times 1
        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Succesfully ran Protect-CFCell" }
    }

    It "fails gracefully when Protect-CFCell powershell cmdlet fails" {
        Mock Protect-CFCell { throw "Something terrible happened while attempting to execute Protect-CFCell" }
        Mock Write-Log {}

        { InstallCFCell } | Should -Throw "Something terrible happened while attempting to execute Protect-CFCell"

        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Something terrible happened while attempting to execute Protect-CFCell"}
        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Failed to execute Protect-CFCell powershell cmdlet. See 'c:\provisions\log.log' for mor info."}
    }

}

Describe "InstallBoshAgent" {
    It "executes the Install-Agent powershell cmdlet" {
        Mock Install-Agent {}
        Mock Write-Log {}

        { InstallBoshAgent } | Should -Not -Throw


        Assert-MockCalled Install-Agent -Times 1 -Scope It -ParameterFilter { $IaaS -eq "vsphere" -and $agentZipPath -eq ".\agent.zip" }
        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Bosh agent successfully installed" }
    }

    It "fails gracefully when Install-Agent powershell cmdlet fails" {
        Mock Install-Agent { throw "Something terrible happened while attempting to execute Install-Agent" }
        Mock Write-Log {}

        { InstallBoshAgent } | Should -Throw "Something terrible happened while attempting to execute Install-Agent"

        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Something terrible happened while attempting to execute Install-Agent" }
        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Failed to execute Install-Agent powershell cmdlet. See 'c:\provisions\log.log' for mor info." }
    }
}

Describe "InstallOpenSSH" {
    It "executes the Install-SSHD powershell cmdlet" {
        Mock Install-SSHD {}
        Mock Write-Log {}

        { InstallOpenSSH } | Should -Not -Throw


        Assert-MockCalled Install-SSHD -Times 1 -Scope It -ParameterFilter { $SSHZipFile -eq ".\OpenSSH-Win64.zip" }
        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "OpenSSH successfully installed" }
    }

    It "fails gracefully when Install-SSHD powershell cmdlet fails" {
        Mock Install-SSHD { throw "Something terrible happened while attempting to execute Install-SSHD" }
        Mock Write-Log {}

        { InstallOpenSSH } | Should -Throw "Something terrible happened while attempting to execute Install-SSHD"

        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Something terrible happened while attempting to execute Install-SSHD" }
        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Failed to execute Install-SSHD powershell cmdlet. See 'c:\provisions\log.log' for mor info." }
    }
}

Describe "CleanUpVM" {
    It "executes the Optimize-Disk and Compress-Disk powershell cmdlet" {
        Mock Optimize-Disk {}
        Mock Compress-Disk {}
        Mock Write-Log {}

        { CleanUpVM } | Should -Not -Throw

        Assert-MockCalled Optimize-Disk -Times 1 -Scope It
        Assert-MockCalled Compress-Disk -Times 1 -Scope It
        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Successfully cleaned up the VM's disk" }
    }

    It "fails gracefully when Optimize-Disk powershell cmdlet fails" {
        Mock Optimize-Disk { throw "Something terrible happened while attempting to execute Optimize-Disk" }
        Mock Compress-Disk {}
        Mock Write-Log {}

        { CleanUpVM } | Should -Throw "Something terrible happened while attempting to execute Optimize-Disk"

        Assert-MockCalled Optimize-Disk -Times 1 -Scope It
        Assert-MockCalled Compress-Disk -Times 0 -Scope It

        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Something terrible happened while attempting to execute Optimize-Disk" }
        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Failed to clean up the VM's disk. See 'c:\provisions\log.log' for mor info." }
    }

    It "fails gracefully when Compress-Disk powershell cmdlet fails" {
        Mock Optimize-Disk {}
        Mock Compress-Disk { throw "Something terrible happened while attempting to execute Compress-Disk" }
        Mock Write-Log {}

        { CleanUpVM } | Should -Throw "Something terrible happened while attempting to execute Compress-Disk"

        Assert-MockCalled Optimize-Disk -Times 1 -Scope It
        Assert-MockCalled Compress-Disk -Times 1 -Scope It

        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Something terrible happened while attempting to execute Compress-Disk" }
        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter { $Message -eq "Failed to clean up the VM's disk. See 'c:\provisions\log.log' for mor info." }
    }
}