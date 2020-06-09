. ./ProvisionVM.ps1

Describe "ProvisionVM" {
    BeforeEach {
        [System.Collections.ArrayList]$provisionerCalls = @()

        Mock Set-RegKeys { $provisionerCalls.Add("Set-RegKeys") }
        Mock InstallCFFeatures { $provisionerCalls.Add("InstallCFFeatures") }
        Mock CopyPSModules { $provisionerCalls.Add("CopyPSModules") }
        Mock InstallOpenSSH { $provisionerCalls.Add("InstallOpenSSH") }
        Mock InstallBoshAgent { $provisionerCalls.Add("InstallBoshAgent") }
        Mock Enable-SSHD { $provisionerCalls.Add("Enable-SSHD") }
        Mock Install-SecurityPoliciesAndRegistries { $provisionerCalls.Add("Install-SecurityPoliciesAndRegistries") }
        Mock Extract-LGPO { $provisionerCalls.Add("Extract-LGPO") }
        Mock Install-WUCerts { $provisionerCalls.Add("Install-WUCerts") }
        Mock Create-VersionFile { $provisionerCalls.Add("Create-VersionFile") }

        if (!(Get-Command "Restart-Computer" -errorAction SilentlyContinue))
        {
            function Restart-Computer() {
                throw "what is happening I should never be invoked"
            }
        }

        Mock Restart-Computer { $provisionerCalls.Add("Restart-Computer") }
    }

    It "sets registry keys to stop zombie load and meltdown exploits" {
        ProvisionVM

        Assert-MockCalled -CommandName Set-RegKeys
    }
    It "installs CFFeatures" {
        ProvisionVM

        Assert-MockCalled -CommandName InstallCFFeatures
    }

    It "copy PSModules is the first provisioner called" {
        ProvisionVM

        Assert-MockCalled -CommandName CopyPSModules
        $provisionerCalls.IndexOf("CopyPSModules") | Should -Be 0
    }
    It "installs BoshAgent" {
        ProvisionVM

        Assert-MockCalled -CommandName InstallBoshAgent
    }
    It "installs OpenSSH before enabling SSH" {
        ProvisionVM

        Assert-MockCalled -CommandName InstallOpenSSH
        $provisionerCalls.IndexOf("InstallOpenSSH") | Should -BeGreaterOrEqual 0
        $provisionerCalls.IndexOf("InstallOpenSSH") | Should -BeLessThan $provisionerCalls.IndexOf("Enable-SSHD")
    }
    It "enables SSHD" {
        ProvisionVM

        Assert-MockCalled -CommandName Enable-SSHD
    }
    It "installs SecurityPoliciesAndRegistries after extracting LGPO" {
        ProvisionVM

        Assert-MockCalled -CommandName Install-SecurityPoliciesAndRegistries

        $provisionerCalls.IndexOf("Extract-LGPO") | Should -BeGreaterOrEqual 0
        $provisionerCalls.IndexOf("Extract-LGPO") | Should -BeLessThan $provisionerCalls.IndexOf("Install-SecurityPoliciesAndRegistries")
    }

    It "extracts LGPO before enabling SSH" {
        ProvisionVM

        Assert-MockCalled -CommandName Extract-LGPO

        $provisionerCalls.IndexOf("Extract-LGPO") | Should -BeGreaterOrEqual 0
        $provisionerCalls.IndexOf("Extract-LGPO") | Should -BeLessThan $provisionerCalls.IndexOf("Enable-SSHD")
    }

    It "fails gracefully when Install-WUCerts helper fails" {
        Mock Install-WUCerts { throw "Something went wrong trying to Install-WUCerts" }
        Mock Write-Log { }
        Mock Write-Warning { }

        { ProvisionVM } | Should -Not -Throw

        Assert-MockCalled Install-WUCerts -Times 1 -Scope It
        Assert-MockCalled Write-Log -Times 1 -Scope It -ParameterFilter {$Message -eq "Something went wrong trying to Install-WUCerts" }
        Assert-MockCalled Write-Warning -Times 1 -Scope It -ParameterFilter {$Message -eq "Failed to retrieve updated root certificates from the public Windows Update Server. This should not impact the successful execution of stembuild construct. If your root certificates are out of date, Diego cells running on VMs built from this stemcell may not be able to make outbound network connections." }

    }

    It "installs WU certs" {
        ProvisionVM

        Assert-MockCalled -CommandName Install-WUCerts
    }

    It "creates a version file" {
        ProvisionVM

        Assert-MockCalled -CommandName Create-VersionFile
    }

    It "restarts as the last command" {
        ProvisionVM

        Assert-MockCalled -CommandName Restart-Computer
        $lastIndex = $provisionerCalls.Count - 1
        $provisionerCalls.IndexOf("Restart-Computer") | Should -Be $lastIndex
    }

}
