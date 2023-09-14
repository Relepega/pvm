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

	name := alias

	if alias == "" {
		name = ver.VersionNumber
	} else if !utils.IsValidFolderName(alias) {
		fmt.Println("Invalid alias or alias too short (minimum required lenght is 3 characters). Using no alias as fallback...")
	}

	client.InstallNewVersion(ver, name)
}
