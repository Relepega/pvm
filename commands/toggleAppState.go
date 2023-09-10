package commands

import (
	windowsClient "WindowsClient"

	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func ToggleAppState(mode string) {
	absSymlinkPath := strings.ReplaceAll(windowsClient.SymlinkDest, "%localappdata%", os.Getenv("localappdata"))

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
