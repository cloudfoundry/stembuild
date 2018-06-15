Push-Location $PSScriptRoot

. ./AutomationHelpers.ps1

#clean-up runonce
Unregister-ScheduledTask -TaskName Task01 -Confirm:$false

CleanUpVM
SysprepVM

Pop-Location