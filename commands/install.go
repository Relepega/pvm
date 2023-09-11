package commands

import (
	windowsClient "WindowsClient"
	"log"
)

func Install(version string, alias string) {
	client := windowsClient.NewClient()

	ver, err := windowsClient.UseVersion(version)

	if err != nil {
		log.Fatalln(err)
	}

	client.InstallNewVersion(ver, alias)
}
