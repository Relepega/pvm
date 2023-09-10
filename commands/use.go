package commands

import (
	windowsClient "WindowsClient"
	"log"
	"path/filepath"

	"os"
	"path"
)

/**
 * @param {slug} either the version or the alias
 */
func Use(slug string) {
	client := windowsClient.NewClient()

	//! assuming that slug can only be the version number for now
	version, err := windowsClient.UseVersion(slug)

	if err != nil {
		log.Fatalln(err)
	}

	// check if path exists and is a directory
	versionPath, _ := filepath.Abs(path.Join(windowsClient.PythonRootContainer, version.VersionNumber))
	fileInfo, err := os.Stat(versionPath)

	if err != nil && !fileInfo.IsDir() {
		os.Exit(1)
	}

	client.MakeSymlink(version.VersionNumber, versionPath)
}
