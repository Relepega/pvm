package commands

import (
	utils "AppUtils"
	windowsClient "WindowsClient"
	"fmt"

	"log"
)

func Install(version string, alias string) {
	client := windowsClient.NewClient()

	ver, err := windowsClient.UseVersion(version)

	if err != nil {
		log.Fatalln(err)
	}

	if utils.IsValidFolderName(alias) {
		client.InstallNewVersion(ver, alias)
	} else {
		fmt.Println("Invalid alias or alias too short (minimum required lenght is 3 characters). Using no alias as fallback...")
		client.InstallNewVersion(ver, "")
	}
}
