package commands

import (
	windowsClient "WindowsClient"
	"log"

	"strings"
)

func Reinstall(version string) {
	client := windowsClient.NewClient()

	ver, err := windowsClient.UseVersion(strings.ToLower(version))

	if err != nil {
		log.Fatalln(err)
	}

	uninstallSingleVersion(ver.VersionNumber)

	client.InstallNewVersion(ver, "") // disabling alias by hardcoding
}
