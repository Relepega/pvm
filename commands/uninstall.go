package commands

import (
	windowsClient "WindowsClient"
	"fmt"
	"path/filepath"
	"strings"

	"os"
)

func Uninstall(folderName string) {
	errorNotFound := fmt.Sprintf("No python installation named \"%s\" has been found.\nPlease install it first.", folderName)

	if strings.ToLower(folderName) == "all" {
		fmt.Println("Uninstalling all installations... ")
		os.RemoveAll(windowsClient.PythonRootContainer)
		fmt.Println("Done!")
		return
	}

	client := windowsClient.NewClient()
	installationFolderPath := filepath.Clean(filepath.Join(client.AppRoot, windowsClient.PythonRootContainer, folderName))

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
