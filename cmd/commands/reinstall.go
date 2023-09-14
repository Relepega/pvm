package commands

import (
	windowsClient "WindowsClient"

	"log"
	"os"
	"path"
	"path/filepath"
)

func Reinstall(folderName string) {
	client := windowsClient.NewClient()
	installationFolderPath := path.Join(client.InstallDir, folderName)

	stat, err := os.Stat(installationFolderPath)
	if err != nil || !stat.IsDir() {
		log.Fatalf("No installation named \"%s\" found", folderName)
	}

	data, err := os.ReadFile(filepath.Join(installationFolderPath, "version"))
	if err != nil {
		log.Fatalln(err)
	}

	ver, err := windowsClient.UseVersion(string(data))

	if err != nil {
		log.Fatalln(err)
	}

	Uninstall(folderName)
	client.InstallNewVersion(ver, folderName)
}
