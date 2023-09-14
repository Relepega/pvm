package commands

import (
	"WindowsClient"

	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Uninstall(folderName string) {
	errorNotFound := fmt.Sprintf("No python installation named \"%s\" has been found.\nPlease install it first.", folderName)

	if strings.ToLower(folderName) == "all" {
		fmt.Println("Uninstalling all installations... ")
		os.RemoveAll(WindowsClient.PythonInstallDirname)
		fmt.Println("Done!")
		return
	}

	client := WindowsClient.NewClient()
	installationFolderPath := filepath.Join(client.InstallDir, folderName)

	stat, err := os.Stat(installationFolderPath)

	if err != nil {
		fmt.Println(errorNotFound)
		return
	}

	if stat.IsDir() {
		fmt.Printf("Uninstalling \"%s\" installation... ", folderName)
		os.RemoveAll(installationFolderPath)
		fmt.Println("Done!")
		return
	}

	fmt.Println(errorNotFound)
}
