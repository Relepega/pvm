# https://www.delftstack.com/howto/powershell/powershell-run-as-administrator/#run-powershell-script-with-arguments-as-administrator
# Self-elevate the script if required
if (-Not ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')) {
    if ([int](Get-CimInstance -Class Win32_OperatingSystem | Select-Object -ExpandProperty BuildNumber) -ge 6000) {
        $Command = "-File `"" + $MyInvocation.MyCommand.Path + "`" " + $MyInvocation.UnboundArguments
        Start-Process -FilePath PowerShell.exe -Verb RunAs -ArgumentList $Command
        Exit
 }
}

$localappdata = [Environment]::GetEnvironmentVariable('localappdata')
$target = "$localappdata\Python"
$pvm_home = Split-Path $PSScriptRoot -Parent

[Environment]::SetEnvironmentVariable("PATH", $Env:PATH + ";$pvm_home\;", [EnvironmentVariableTarget]::Machine)
[Environment]::SetEnvironmentVariable("PATH", $Env:PATH + ";$target\;$target\Scripts\;", [EnvironmentVariableTarget]::User)

Remove-Item $env:LOCALAPPDATA\Microsoft\WindowsApps\python.exe
Remove-Item $env:LOCALAPPDATA\Microsoft\WindowsApps\python3.exe

# Read-Host -Prompt "Press any key to continue..."