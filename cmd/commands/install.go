package Commands

import (
	"AppUtils"
	"PvmState"
	"WindowsClient"

	"fmt"
	"log"
	"strings"
)

func Install(version string, alias string) {
	if strings.ToLower(version) == "pvm" {
		PvmState.Handler("install")
		return
	}

	client := WindowsClient.NewClient()

	ver, err := WindowsClient.UseVersion(version)

	if err != nil {
		log.Fatalln(err)
	}

	name := alias

	if alias == "" {
		name = ver.VersionNumber
	} else if !AppUtils.IsValidFolderName(alias) {
		fmt.Println("Invalid alias or alias too short (minimum required lenght is 3 characters). Using no alias as fallback...")
	}

	client.InstallNewVersion(ver, name)
}
