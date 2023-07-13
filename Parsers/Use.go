package Parsers

import (
	utils "Utils"
	windowsClient "WindowsClient"
	"path/filepath"

	"os"
	"path"
)

func UseParserHandler(version string, client *windowsClient.Client) {
	utils.IsValidPythonVersion(version)

	// check if path exists and is a directory
	versionPath, _ := filepath.Abs(path.Join(windowsClient.PythonRootContainer, version))
	fileInfo, err := os.Stat(versionPath)

	if err != nil && !fileInfo.IsDir() {
		os.Exit(1)
	}

	client.MakeSymlink(version, versionPath)
}
