package commands

import (
	"PvmState"
	"WindowsClient"

	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Uninstall(name string) {
	if strings.ToLower(name) == "pvm" {
		PvmState.Uninstall()
		return
	}

	errorNotFound := fmt.Sprintf("No python installation named \"%s\" has been found.\nPlease install it first.", name)

	if strings.ToLower(name) == "all" {
		fmt.Println("Uninstalling all installations... ")
		os.RemoveAll(WindowsClient.PythonInstallDirname)
		fmt.Println("Done!")
		return
	}

	client := WindowsClient.NewClient()
	installationFolderPath := filepath.Join(client.InstallDir, name)

	stat, err := os.Stat(installationFolderPath)

	if err != nil {
		fmt.Println(errorNotFound)
		return
	}

	if stat.IsDir() {
		fmt.Printf("Uninstalling \"%s\" installation... ", name)
		os.RemoveAll(installationFolderPath)
		fmt.Println("Done!")
		return
	}

	fmt.Println(errorNotFound)
}
