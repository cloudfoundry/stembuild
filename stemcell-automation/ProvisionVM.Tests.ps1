. ./ProvisionVM.ps1

Describe "ProvisionVM" {
    BeforeEach {
        Mock InstallCFFeatures
        Mock CopyPSModules
        Mock InstallOpenSSH
        Mock InstallBoshAgent
        Mock Enable-SSHD
        Mock Install-SecurityPoliciesAndRegistries
        Mock Enable-HyperV
    }

    It "installs CFFeatures" {
        ProvisionVM

        Assert-MockCalled -CommandName InstallCFFeatures
    }
    It "enables Hyper-V when possible" {
        ProvisionVM

        Assert-MockCalled -CommandName Enable-HyperV
    }

    It "copy PSModules" {
        ProvisionVM

        Assert-MockCalled -CommandName CopyPSModules
    }
    It "installs BoshAgent" {
        ProvisionVM

        Assert-MockCalled -CommandName InstallBoshAgent
    }
    It "installs OpenSSH" {
        ProvisionVM

        Assert-MockCalled -CommandName InstallOpenSSH
    }
    It "enables SSHD" {
        ProvisionVM

        Assert-MockCalled -CommandName Enable-SSHD
    }
    It "does not install SecurityPoliciesAndRegistries" {
        ProvisionVM

        Assert-MockCalled -Times 0 -CommandName Install-SecurityPoliciesAndRegistries
    }

}
