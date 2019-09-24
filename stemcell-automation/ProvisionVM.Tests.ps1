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
        Mock Enable-HyperV { $provisionerCalls.Add("Extract-LGPO") }
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
    It "does not install SecurityPoliciesAndRegistries" {
        ProvisionVM

        Assert-MockCalled -Times 0 -CommandName Install-SecurityPoliciesAndRegistries
    }

    It "extracts LGPO before enabling SSH" {
        ProvisionVM

        Assert-MockCalled -CommandName Extract-LGPO

        $provisionerCalls.IndexOf("Extract-LGPO") | Should -BeGreaterOrEqual 0
        $provisionerCalls.IndexOf("Extract-LGPO") | Should -BeLessThan $provisionerCalls.IndexOf("Enable-SSHD")
    }

}
