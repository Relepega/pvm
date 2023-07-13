package Parsers

import (
	utils "Utils"
	windowsClient "WindowsClient"

	"fmt"
	"os"
	"path"
	"strings"
)

func InstallParserHandler(version string, client *windowsClient.Client) {
	v := strings.ToLower(version)

	if v == "latest" {
		length := len(client.PythonVersions.Stable)
		v = client.PythonVersions.Stable[length-1]
	} else {
		utils.IsValidPythonVersion(v)
	}

	versionPath := path.Join(windowsClient.PythonRootContainer, v)

	// https://stackoverflow.com/a/40624033
	if stat, err := os.Stat(versionPath); err == nil && stat.IsDir() {
		fmt.Printf("Python %s is already installed. Please use the command \"reinstall\" instead.", v)
	} else {
		client.InstallNewVersion(v)
	}
}
