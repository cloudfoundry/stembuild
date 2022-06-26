﻿<#
.Synopsis
  Sysprep Utilities
.Description
  This cmdlet enables enabling a local security policy for a stemcell
#>
function Enable-LocalSecurityPolicy {
  Param (
    [string]$PolicySource =$(throw "Policy backup filepath is required")
  )
  Write-Log "Starting LocalSecurityPolicy"

  # Convert registry.txt files into registry.pol files
  $MachineDir="$PolicySource/DomainSysvol/GPO/Machine"
  LGPO.exe /r "$MachineDir/registry.txt" /w "$MachineDir/registry.pol"
  if ($LASTEXITCODE -ne 0) {
    Write-Error "Generating policy: Machine"
  }

  $UserDir="$PolicySource/DomainSysvol/GPO/User"
  LGPO.exe /r "$UserDir/registry.txt" /w "$UserDir/registry.pol"
  if ($LASTEXITCODE -ne 0) {
    Write-Error "Generating policy: User"
  }

  # Apply policies
  LGPO.exe /g "$PolicySource/DomainSysvol" /v
  if ($LASTEXITCODE -ne 0) {
    Write-Error "Applying policy: $PolicySource/DomainSysvol"
  }

  Write-Log "Ending LocalSecurityPolicy"
}

<#
.Synopsis
  Sysprep Utilities
.Description
  This cmdlet creates the Unattend file for sysprep
#>
function Create-Unattend {
  Param (
    [string]$UnattendDestination = "C:\Windows\Panther\Unattend",
    [string]$NewPassword,
    [string]$ProductKey,
    [string]$Organization,
    [string]$Owner
  )
  Write-Log "Starting Create-Unattend"

  New-Item -ItemType directory $UnattendDestination -Force
  $UnattendPath = Join-Path $UnattendDestination "unattend.xml"

  Write-Log "Writing unattend.xml to $UnattendPath"

  $ProductKeyXML=""
  if ($ProductKey -ne "") {
    $ProductKeyXML="<ProductKey>$ProductKey</ProductKey>"
  }

  $OrganizationXML="<RegisteredOrganization />"
  if ($Organization -ne "" -and $Organization -ne $null) {
    $OrganizationXML="<RegisteredOrganization>$Organization</RegisteredOrganization>"
  }

  $OwnerXML="<RegisteredOwner />"
  if ($Owner -ne "" -and $Owner -ne $null) {
    $OwnerXML="<RegisteredOwner>$Owner</RegisteredOwner>"
  }

  $AdministratorPasswordXML = ""
  if ($NewPassword -ne "" -and $NewPassword -ne $null) {
    $NewPassword = [system.convert]::ToBase64String([system.text.encoding]::Unicode.GetBytes($NewPassword + "AdministratorPassword"))
    $AdministratorPasswordXML = @"
      <UserAccounts>
        <AdministratorPassword>
          <Value>$NewPassword</Value>
          <PlainText>false</PlainText>
        </AdministratorPassword>
      </UserAccounts>
"@
  }

  $PostUnattend = @"
<?xml version="1.0" encoding="utf-8"?>
<unattend xmlns="urn:schemas-microsoft-com:unattend">
  <settings pass="specialize">
    <component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
      <OEMInformation>
        <HelpCustomized>false</HelpCustomized>
      </OEMInformation>
      <ComputerName>*</ComputerName>
      <TimeZone>UTC</TimeZone>
      $ProductKeyXML
      $OrganizationXML
      $OwnerXML
    </component>
    <component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-ServerManager-SvrMgrNc" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
      <DoNotOpenServerManagerAtLogon>true</DoNotOpenServerManagerAtLogon>
    </component>
    <component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-OutOfBoxExperience" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
      <DoNotOpenInitialConfigurationTasksAtLogon>true</DoNotOpenInitialConfigurationTasksAtLogon>
    </component>
    <component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-Security-SPP-UX" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
      <SkipAutoActivation>true</SkipAutoActivation>
    </component>
    <component name="Microsoft-Windows-NetBT" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
        <Interfaces>
            <Interface wcm:action="add">
                <NetbiosOptions>2</NetbiosOptions>
                <Identifier>Ethernet0</Identifier>
            </Interface>
        </Interfaces>
    </component>
    <component name="Microsoft-Windows-Deployment" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
      <RunAsynchronous>
        <RunAsynchronousCommand wcm:action="add">
          <Path>powershell Enable-AgentService</Path>
          <Order>1</Order>
          <Description>Enable Bosh Agent Service</Description>
        </RunAsynchronousCommand>
      </RunAsynchronous>
    </component>
  </settings>
  <settings pass="generalize">
    <component name="Microsoft-Windows-PnpSysprep" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
      <PersistAllDeviceInstalls>false</PersistAllDeviceInstalls>
      <DoNotCleanUpNonPresentDevices>false</DoNotCleanUpNonPresentDevices>
    </component>
  </settings>
  <settings pass="oobeSystem">
    <component name="Microsoft-Windows-International-Core" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
      <InputLocale>en-US</InputLocale>
      <SystemLocale>en-US</SystemLocale>
      <UILanguage>en-US</UILanguage>
      <UserLocale>en-US</UserLocale>
    </component>
    <component name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
      <OOBE>
        <HideEULAPage>true</HideEULAPage>
        <ProtectYourPC>3</ProtectYourPC>
        <NetworkLocation>Home</NetworkLocation>
        <HideWirelessSetupInOOBE>true</HideWirelessSetupInOOBE>
      </OOBE>
      <TimeZone>UTC</TimeZone>
      $AdministratorPasswordXML
    </component>
  </settings>
</unattend>
"@

  Out-File -FilePath $UnattendPath -InputObject $PostUnattend -Encoding utf8
}

<#
.Synopsis
  Sanity check that the unattend.xml shipped with GCP has not changed.
.Description
  Sanity check that the unattend.xml shipped with GCP has not changed.
#>
function Check-Default-GCP-Unattend() {

  [xml]$Expected = @'
<?xml version="1.0" encoding="utf-8"?>
<unattend xmlns="urn:schemas-microsoft-com:unattend">
  <!--
  For more information about unattended.xml please refer too
  http://technet.microsoft.com/en-us/library/cc722132(v=ws.10).aspx
  -->
  <settings pass="generalize">
    <component name="Microsoft-Windows-PnpSysprep" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
      <PersistAllDeviceInstalls>true</PersistAllDeviceInstalls>
    </component>
  </settings>
  <settings pass="specialize">
    <component name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
      <!-- Random ComputerName, will be replaced by specialize script -->
      <ComputerName></ComputerName>
      <TimeZone>Greenwich Standard Time</TimeZone>
    </component>
  </settings>
  <settings pass="oobeSystem">
    <!-- Setting Location Information -->
    <component name="Microsoft-Windows-International-Core" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
      <InputLocale>en-us</InputLocale>
      <SystemLocale>en-us</SystemLocale>
      <UILanguage>en-us</UILanguage>
      <UserLocale>en-us</UserLocale>
    </component>
    <component name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
      <OOBE>
        <!-- Setting EULA -->
        <HideEULAPage>true</HideEULAPage>
        <!-- Setting network location to public -->
        <NetworkLocation>Other</NetworkLocation>
        <!-- Hide Wirelss setup -->
        <HideWirelessSetupInOOBE>true</HideWirelessSetupInOOBE>
        <ProtectYourPC>1</ProtectYourPC>
        <SkipMachineOOBE>true</SkipMachineOOBE>
        <SkipUserOOBE>true</SkipUserOOBE>
      </OOBE>
      <!-- Setting timezone to GMT -->
      <ShowWindowsLive>false</ShowWindowsLive>
      <TimeZone>Greenwich Standard Time</TimeZone>
      <!--Setting OEM information -->
      <OEMInformation>
        <Manufacturer>Google Cloud Platform</Manufacturer>
        <Model>Google Compute Engine Virtual Machine</Model>
        <SupportURL>https://support.google.com/enterprisehelp/answer/142244?hl=en#cloud</SupportURL>
        <Logo>C:\Program Files\Google Compute Engine\sysprep\gcp.bmp</Logo>
      </OEMInformation>
    </component>
  </settings>
</unattend>
'@

  $UnattendPath = "C:\Program Files\Google\Compute Engine\sysprep\unattended.xml"
  [xml]$Unattend = (Get-Content -Path $UnattendPath)

  if (-Not ($Unattend.xml.Equals($Expected.xml))) {
  Write-Error "The unattend.xml shipped with GCP has changed."
  }
}

function Create-Unattend-GCP() {
  Param (
    [string]$UnattendDestination = "C:\Program Files\Google\Compute Engine\sysprep"
  )
  $UnattendXML = @'
<?xml version="1.0" encoding="utf-8"?>
<unattend xmlns="urn:schemas-microsoft-com:unattend">
  <!--
  For more information about unattended.xml please refer too
  http://technet.microsoft.com/en-us/library/cc722132(v=ws.10).aspx
  -->
  <settings pass="generalize">
    <component name="Microsoft-Windows-PnpSysprep" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
      <PersistAllDeviceInstalls>true</PersistAllDeviceInstalls>
    </component>
  </settings>
  <settings pass="specialize">
    <component name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
      <!-- Random ComputerName, will be replaced by specialize script -->
      <ComputerName></ComputerName>
      <TimeZone>UTC</TimeZone>
    </component>
    <component name="Microsoft-Windows-Deployment" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
      <RunAsynchronous>
        <RunAsynchronousCommand wcm:action="add">
          <Path>powershell Enable-AgentService</Path>
          <Order>1</Order>
          <Description>Enable Bosh Agent Service</Description>
        </RunAsynchronousCommand>
      </RunAsynchronous>
    </component>
  </settings>
  <settings pass="oobeSystem">
    <!-- Setting Location Information -->
    <component name="Microsoft-Windows-International-Core" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
      <InputLocale>en-us</InputLocale>
      <SystemLocale>en-us</SystemLocale>
      <UILanguage>en-us</UILanguage>
      <UserLocale>en-us</UserLocale>
    </component>
    <component name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
      <OOBE>
        <!-- Setting EULA -->
        <HideEULAPage>true</HideEULAPage>
        <!-- Setting network location to public -->
        <NetworkLocation>Other</NetworkLocation>
        <!-- Hide Wirelss setup -->
        <HideWirelessSetupInOOBE>true</HideWirelessSetupInOOBE>
        <ProtectYourPC>3</ProtectYourPC>
        <SkipMachineOOBE>true</SkipMachineOOBE>
        <SkipUserOOBE>true</SkipUserOOBE>
      </OOBE>
      <!-- Setting timezone to GMT -->
      <ShowWindowsLive>false</ShowWindowsLive>
      <TimeZone>UTC</TimeZone>
      <!--Setting OEM information -->
      <OEMInformation>
        <Manufacturer>Google Cloud Platform</Manufacturer>
        <Model>Google Compute Engine Virtual Machine</Model>
        <SupportURL>https://support.google.com/enterprisehelp/answer/142244?hl=en#cloud</SupportURL>
        <Logo>C:\Program Files\Google Compute Engine\sysprep\gcp.bmp</Logo>
      </OEMInformation>
    </component>
  </settings>
</unattend>
'@

  $UnattendPath = Join-Path $UnattendDestination "unattended.xml"

  Out-File -FilePath $UnattendPath -InputObject $UnattendXML -Encoding utf8 -Force
}

function Remove-WasPassProcessed {
  Param (
    [string]$AnswerFilePath
  )

  If (!$(Test-Path $AnswerFilePath)) {
    Throw "Answer file $AnswerFilePath does not exist"
  }

  Write-Log "Removing wasPassProcessed"

  $content = [xml](Get-Content $AnswerFilePath)

  foreach ($specializeBlock in $content.unattend.settings) {
    $specializeBlock.RemoveAttribute("wasPassProcessed")
  }

  $content.Save($AnswerFilePath)
}

function Remove-UserAccounts {
  Param (
    [string]$AnswerFilePath
  )

  If (!$(Test-Path $AnswerFilePath)) {
    Throw "Answer file $AnswerFilePath does not exist"
  }

  Write-Log "Removing UserAccounts block from Answer File"

  $content = [xml](Get-Content $AnswerFilePath)
  $mswShellSetup =  (($content.unattend.settings|where {$_.pass -eq 'oobeSystem'}).component|where {$_.name -eq "Microsoft-Windows-Shell-Setup"})

  if ($mswShellSetup -eq $Null) {
    Throw "Could not locate oobeSystem XML block. You may not be running this function on an answer file."
  }

  $userAccountsBlock = $mswShellSetup.UserAccounts

  if ($userAccountsBlock.Count -eq 0) {
    Return
  }

  $mswShellSetup.RemoveChild($userAccountsBlock)

  $content.Save($AnswerFilePath)
}

function Update-AWS2016Config
{
  $LaunchConfigJson = 'C:\ProgramData\Amazon\EC2-Windows\Launch\Config\LaunchConfig.json'
  $LaunchConfig = Get-Content $LaunchConfigJson -raw | ConvertFrom-Json
  $LaunchConfig.addDnsSuffixList = $False
  $LaunchConfig.extendBootVolumeSize = $False
  $LaunchConfig | ConvertTo-Json | Set-Content $LaunchConfigJson
}

function Create-Unattend-AWS
{
  $UnattendedXmlPath = 'C:\ProgramData\Amazon\EC2-Windows\Launch\Sysprep\Unattend.xml'
  $UnattendedContent = [xml](Get-Content $UnattendedXmlPath)
  $SpecializeSettings = ($UnattendedContent.unattend.settings | Where-Object { $_.pass -EQ "specialize" })
  $WindowsDeploymentComponent = ($SpecializeSettings.component | Where-Object { $_.name -EQ "Microsoft-Windows-Deployment" })
  $rynsync = $WindowsDeploymentComponent.RunSynchronous
  $runsynccommand = $UnattendedContent.CreateElement("RunSynchronousCommand", $UnattendedContent.unattend.xmlns)
  $rynsync.AppendChild($runsynccommand)
  $runsynccommand.SetAttribute("action", $WindowsDeploymentComponent.wcm, "add")
  $pathElement = $UnattendedContent.CreateElement("Path", $UnattendedContent.unattend.xmlns)
  $pathText = $UnattendedContent.CreateTextNode("powershell Enable-AgentService")
  $pathElement.AppendChild($pathText)
  $runsynccommand.AppendChild($pathElement)
  $orderElement = $UnattendedContent.CreateElement("Order", $UnattendedContent.unattend.xmlns)
  $orderText = $UnattendedContent.CreateTextNode("3")
  $orderElement.AppendChild($orderText)
  $runsynccommand.AppendChild($orderElement)

  $UnattendedContent.Save($UnattendedXmlPath)
}

function Enable-AWS2016Sysprep {
  # Enable sysprep
  cd 'C:\ProgramData\Amazon\EC2-Windows\Launch\Scripts'
  ./InitializeInstance.ps1 -Schedule
  ./SysprepInstance.ps1
}

<#
.Synopsis
  Sysprep Utilities
.Description
  This cmdlet runs Sysprep and generalizes a VM so it can be a BOSH stemcell
#>
function Invoke-Sysprep()
{
  Param (
    [string]$IaaS = $( Throw "Provide the IaaS this stemcell will be used for" ),
    [string]$NewPassword,
    [string]$ProductKey = "",
    [string]$Organization = "",
    [string]$Owner = "",
    [switch]$SkipLGPO,
    [switch]$EnableRDP
  )

  Write-Log "Invoking Sysprep for IaaS: ${IaaS}"

  $OsVersion = Get-OSVersion

  # WARN WARN: this should be removed when Microsoft fixes this bug
  # See tracker story https://www.pivotaltracker.com/story/show/150238324
  # Skip sysprep if using Windows Server 2016 insider build with UALSVC bug
  $RegPath = "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion"
  If ((Get-ItemProperty -Path $RegPath).CurrentBuildNumber -Eq '16278')
  {
    Stop-Computer
  }

  Allow-NTPSync

  if (-Not $SkipLGPO)
  {
    if (-Not (Test-Path "C:\Windows\LGPO.exe")) {
      Throw "Error: LGPO.exe is expected to be installed to C:\Windows\LGPO.exe"
    }

    switch ($OsVersion)
    {
      "windows2012R2" {
        Enable-LocalSecurityPolicy (Join-Path $PSScriptRoot "cis-merge-2012R2")
      }

      "windows1803" {
        Enable-LocalSecurityPolicy (Join-Path $PSScriptRoot "cis-merge-1803")
      }

      "windows2019" {
        Enable-LocalSecurityPolicy (Join-Path $PSScriptRoot "cis-merge-2019")
      }
    }
  }

  switch ($IaaS) {
    "aws" {
      Disable-AgentService
      Create-Unattend-AWS
      Update-AWS2016Config
      Enable-AWS2016Sysprep
    }
    "gcp" {
      Disable-AgentService
      Create-Unattend-GCP
      GCESysprep
    }
    "azure" {
      C:\Windows\System32\Sysprep\sysprep.exe /generalize /quiet /oobe /quit
    }
    "vsphere" {
      Disable-AgentService
      Create-Unattend -NewPassword $NewPassword -ProductKey $ProductKey `
        -Organization $Organization -Owner $Owner

      Invoke-Expression -Command 'C:/windows/system32/sysprep/sysprep.exe /generalize /oobe /unattend:"C:/Windows/Panther/Unattend/unattend.xml" /quiet /shutdown'
    }
    Default { Throw "Invalid IaaS '${IaaS}' supported platforms are: AWS, Azure, GCP and Vsphere" }
  }
}

function ModifyInfFile() {
  Param(
    [string]$InfFilePath = $(Throw "inf file path missing"),
    [string]$KeyName = $(Throw "keyname missing"),
    [string]$KeyValue = $(Throw "keyvalue missing")
  )

  $Regex = "^$KeyName"
  $TempFile = $InfFilePath + ".tmp"

  Get-Content $InfFilePath | ForEach-Object {
    $ValueToWrite=$_
    if($_ -match $Regex) {
      $ValueToWrite="$KeyName=$KeyValue"
    }
    $ValueToWrite | Out-File -Append $TempFile
  }

  Move-Item -Path $TempFile -Destination $InfFilePath -Force
}

function Allow-NTPSync() {
      Set-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Services\W32Time\Config" -Name 'MaxNegPhaseCorrection' -Value 0xFFFFFFFF -Type dword
      Set-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Services\W32Time\Config" -Name 'MaxPosPhaseCorrection' -Value 0xFFFFFFFF -Type dword
}
