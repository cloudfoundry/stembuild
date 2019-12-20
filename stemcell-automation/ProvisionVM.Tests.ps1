. ./ProvisionVM.ps1

Describe "ProvisionVM" {
    BeforeEach {
        [System.Collections.ArrayList]$provisionerCalls = @()

        Mock InstallCFFeatures { $provisionerCalls.Add("InstallCFFeatures") }
        Mock CopyPSModules { $provisionerCalls.Add("CopyPSModules") }
        Mock InstallOpenSSH { $provisionerCalls.Add("InstallOpenSSH") }
        Mock InstallBoshAgent { $provisionerCalls.Add("InstallBoshAgent") }
        Mock Enable-SSHD { $provisionerCalls.Add("Enable-SSHD") }
        Mock Install-SecurityPoliciesAndRegistries { $provisionerCalls.Add("Install-SecurityPoliciesAndRegistries") }
        Mock Extract-LGPO { $provisionerCalls.Add("Extract-LGPO") }
        Mock Enable-HyperV { $provisionerCalls.Add("Enable-HyperV") }
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

    It "installs CFFeatures" {
        ProvisionVM

        Assert-MockCalled -CommandName InstallCFFeatures
    }
    It "enables Hyper-V when possible" {
        ProvisionVM

        Assert-MockCalled -CommandName Enable-HyperV
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
