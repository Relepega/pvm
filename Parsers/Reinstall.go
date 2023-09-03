package Parsers

import (
	utils "Utils"
	windowsClient "WindowsClient"
	"strings"
)

func ReinstallParserHandler(version string, client *windowsClient.Client) {
	v := strings.ToLower(version)
	utils.IsValidPythonVersion(v)

	uninstallSingleVersion(v)
	InstallParserHandler(v, "", client)
}
