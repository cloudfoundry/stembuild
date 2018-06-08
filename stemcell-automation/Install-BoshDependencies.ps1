#. ./AutomationHelpers.ps1

#Install Bosh Agent
InstallBoshAgent

#Install SSH Deamon
InstallOpenSSH

#Install Bosh Powershell Modules
CopyPSModules
InstallCFCell
InstallCFFeatures