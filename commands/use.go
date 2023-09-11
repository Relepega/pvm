package commands

import (
	windowsClient "WindowsClient"
	"log"
	"os"
	"path"
	"path/filepath"
)

/**
 * @param {slug} either the version or the alias
 */
func Use(slug string) {
	client := windowsClient.NewClient()
	installationPath, _ := filepath.Abs(path.Join(client.AppRoot, windowsClient.PythonRootContainer, slug))
	stat, err := os.Stat(installationPath)

	if err != nil || !stat.IsDir() {
		log.Fatalf("\"%s\" is not a valid installation. Use the 'list' command to list the current installations", slug)
	}

	client.MakeSymlink(slug, installationPath)
}
