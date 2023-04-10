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

# Items to remove
$basedir = '%localappdata%\Python'
$targets = "$basedir\Scripts", "$basedir\", "$pvm_home\"


# Get current path
$Environment = [System.Environment]::GetEnvironmentVariable("Path", "Machine")

# Remove items
foreach($target in $targets){
	foreach ($path in ($Environment).Split(";")) {
		if ($path -like "*$target*") {
			$Environment = $Environment.Replace($Path ,"")
		}
	}
}

# Save updated path
[System.Environment]::SetEnvironmentVariable("Path", $Environment, "Machine")

New-Item $env:LOCALAPPDATA\Microsoft\WindowsApps\python.exe
New-Item $env:LOCALAPPDATA\Microsoft\WindowsApps\python3.exe

# Read-Host -Prompt "Press any key to continue..."