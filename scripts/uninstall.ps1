# https://www.delftstack.com/howto/powershell/powershell-run-as-administrator/#run-powershell-script-with-arguments-as-administrator
# Self-elevate the script if required
if (-Not ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')) {
	if ([int](Get-CimInstance -Class Win32_OperatingSystem | Select-Object -ExpandProperty BuildNumber) -ge 6000) {
		$Command = "-File `"" + $MyInvocation.MyCommand.Path + "`" " + $MyInvocation.UnboundArguments
		Start-Process -FilePath PowerShell.exe -Verb RunAs -ArgumentList $Command
		Exit
 }
}

$pvm_home = Split-Path $PSScriptRoot -Parent
$localappdata = [Environment]::GetEnvironmentVariable('localappdata')

$basedir = "$localappdata\Python"

# Items to remove
$user_targets = "$basedir\Scripts\", "$basedir\"
$machine_targets = "$pvm_home\"


# Get current path from Machine scope
$Environment = [System.Environment]::GetEnvironmentVariable("Path", "Machine")

# Remove items
foreach($target in $user_targets){
	foreach ($path in ($Environment).Split(";")) {
		if ($path -like "*$target*") {
			$Environment = $Environment.Replace($Path ,"")
		}
	}
}

# Save updated path to Machine scope
[System.Environment]::SetEnvironmentVariable("Path", $Environment, "Machine")

# Get current path from User scope
$Environment = [System.Environment]::GetEnvironmentVariable("Path", "User")

# Remove items
foreach($target in $machine_targets){
	foreach ($path in ($Environment).Split(";")) {
		if ($path -like "*$target*") {
			$Environment = $Environment.Replace($Path ,"")
		}
	}
}

# Save updated path to User scope
[System.Environment]::SetEnvironmentVariable("Path", $Environment, "User")

New-Item $env:LOCALAPPDATA\Microsoft\WindowsApps\python.exe
New-Item $env:LOCALAPPDATA\Microsoft\WindowsApps\python3.exe

# Read-Host -Prompt "Press any key to continue..."