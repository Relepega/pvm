package main

import (
	"Parsers"
	"Utils"
	WindowsClient "WindowsClient"

	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const PvmVersion = "1.0.0"

func help() {
	fmt.Println("\nPVM (Python Version Manager) for Windows")
	fmt.Println("Running version " + PvmVersion + ".")
	fmt.Println("\nUsage:")
	fmt.Println(" ")
	fmt.Println(`  pvm [--]install <version>        : The version can be a specific version or "latest" for the`)
	fmt.Println(`                                         current latest stable version. Aliased as [-]i.`)
	fmt.Println(`  pvm [--]reinstall <version>      : The version must be a specific version. Aliased as [-]r.`)
	fmt.Println(`  pvm [--]uninstall <version>      : The version can either be a specific version or "all" for`)
	fmt.Println(`                                         uninstalling all the currently installed versions. Aliased as [-]u.`)
	fmt.Println("  pvm [--]use <version>            : Switch to use the specified version.")
	fmt.Println(`  pvm [--]list <mode>              : Type "All" for listing all stable and unstable versions released;`)
	fmt.Println(`                                         "Installed" for listing all the installed versions`)
	fmt.Println(`                                         and "latest" for the latest 5 version for each major python`)
	fmt.Println(`                                         version (python 2, python 3, etc...). Aliased as [-]l.`)
	fmt.Println("  pvm on                           : Enable python version management.")
	fmt.Println("  pvm off                          : Disable python version management.")
	// fmt.Println("  pvm update                       : Automatically update pvm to the latest version.") // todo
	fmt.Println("  pvm [--]help                     : Displays this help message. Aliased as [-]h.")
	fmt.Println("  pvm [--]version                  : Displays the current running version of pvm for Windows. Aliased as [-]v.")
	fmt.Println(" ")
}

func version(client *WindowsClient.Client) {
	fmt.Println("PVM (Python Version Manager) for Windows")
	fmt.Println("----------------------------------------")
	fmt.Println("Version: " + PvmVersion)
	fmt.Println("Arch:    " + client.Arch)
	fmt.Println("AppRoot: " + Utils.GetWorkingDir())
}

func pvmOnOff(mode string) {
	absSymlinkPath := strings.ReplaceAll(WindowsClient.SymlinkDest, "%localappdata%", os.Getenv("localappdata"))

	// fetch MACHINE-scope PATH enviroment variable and split it
	cmd := exec.Command("powershell.exe", "-Command", `[System.Environment]::GetEnvironmentVariable("Path", "Machine")`)

	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	pathElements := strings.Split(string(out), ";")

	// set psCmd and message based on the case
	var psCmd string
	var message string

	if mode == "enable" {
		pathsToCheck := []string{absSymlinkPath + "\\", absSymlinkPath + "\\Scripts\\"}
		flags := 0

		// check if symlink is already in path
		for _, s := range pathElements {
			if s == pathsToCheck[0] || s == pathsToCheck[1] {
				flags++
			}

			if flags == 2 {
				return
			}
		}

		psCmd = fmt.Sprintf(`[Environment]::SetEnvironmentVariable("PATH", $Env:PATH + ";%s;", [EnvironmentVariableTarget]::Machine)`, strings.Join(pathsToCheck, ";"))
		message = "Python version management is now ENABLED. Restart your PC to complete the process."
	} else {
		/*
			create new filtered path string,
			iterate through all the items
			and filter out pvm paths and newlines
		*/
		var filteredPathElements []string

		for i, s := range pathElements {
			if strings.Contains(s, absSymlinkPath) {
				continue
			}

			if i+1 == len(pathElements) {
				continue
			}

			filteredPathElements = append(filteredPathElements, s)
		}

		newPath := strings.TrimSpace(strings.Join(filteredPathElements, ";"))
		psCmd = fmt.Sprintf(`[Environment]::SetEnvironmentVariable("PATH", "%s" , [EnvironmentVariableTarget]::Machine)`, newPath)
		message = "Python version management is now DISABLED. Restart your PC to complete the process."
	}

	/*
		Create a temporary ps1 script file for disabling symlinks
		functionality through filtering their path in $PATH machine env var
	*/
	file, err := os.CreateTemp("", "pvm_temp-*.ps1")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString(psCmd)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	psCmd = `Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope LocalMachine; ` + file.Name() + `; Set-ExecutionPolicy -ExecutionPolicy Restricted -Scope LocalMachine`
	cmd = exec.Command("powershell.exe", "-noprofile", "Start-Process", "-WindowStyle", "hidden", "-Verb", "RunAs", "-Wait", "powershell.exe", `-Args "`+psCmd+`"`)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(message)
}

func main() {
	args := os.Args
	windowsClient := WindowsClient.NewClient()

	if len(args) == 1 {
		fmt.Println("Not enough arguments, please use \"pvm --help\" for more details...")
		os.Exit(1)
	}

	if len(args) == 4 && args[1] == "install" {
		println("suca")
		Parsers.InstallParserHandler(args[2], args[3], windowsClient)
		return
	}

	switch args[1] {
	case "install", "--install", "i", "-i":
		Parsers.InstallParserHandler(args[2], "", windowsClient)
	case "reinstall", "--reinstall", "r", "-r":
		Parsers.ReinstallParserHandler(args[2], windowsClient)
	case "uninstall", "--uninstall", "u", "-u":
		Parsers.UninstallParserHandler(args[2])
	case "use", "--use":
		Parsers.UseParserHandler(args[2], windowsClient)
	case "list", "--list", "l", "-l":
		Parsers.ListParserHandler(args[2], windowsClient)
	case "on":
		pvmOnOff("enable")
	case "off":
		pvmOnOff("disable")
	case "help", "--help", "h", "-h":
		help()
	case "version", "--version", "v", "-v":
		version(windowsClient)
	default:
		help()
	}

}
