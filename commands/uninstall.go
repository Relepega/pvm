package commands

import (
	windowsClient "WindowsClient"
	"fmt"
	"log"

	"os"
	"path"
	"path/filepath"
	"strings"
)

func Uninstall(userInput string) {
	if strings.ToLower(userInput) == "all" {
		os.RemoveAll(windowsClient.PythonRootContainer)
		return
	}

	if ver, err := windowsClient.UseVersion(userInput); err != nil {
		uninstallSingleVersion(ver.VersionNumber)
		return
	} else {
		fmt.Printf("Python %s is already installed. Please use the command \"reinstall\" instead.", ver.VersionNumber)
		log.Fatalln(err)
	}

	// TODO: Handle aliases

	aliasedPath := filepath.Join(windowsClient.PythonRootContainer, userInput)

	if stat, err := os.Stat(aliasedPath); err != nil && stat.IsDir() {
		os.RemoveAll(stat.Name())
	} else {
		fmt.Printf("No python installation aliased as \"%s\" has been found.\n", userInput)
		fmt.Println("Please use the command \"install\" instead.")
	}
}

func uninstallSingleVersion(version string) {
	// check if path exists and is a directory
	versionPath := path.Join(windowsClient.PythonRootContainer, version)
	fileInfo, err := os.Stat(versionPath)

	fmt.Printf("Uninstalling Python %s... ", version)

	if err != nil {
		fmt.Println(" ")
		log.Fatalln(err)
		os.Exit(1)
	}

	if !fileInfo.IsDir() {
		fmt.Println(" ")
		log.Fatalf(`"%s" is not a directory.`, versionPath)
		os.Exit(1)
	}

	os.RemoveAll(filepath.Clean(versionPath))

	fmt.Println("Done!")
}
