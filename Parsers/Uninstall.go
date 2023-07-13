package Parsers

import (
	windowsClient "WindowsClient"
	"fmt"
	"log"

	"os"
	"path"
	"path/filepath"
	"strings"
)

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

func UninstallParserHandler(version string) {
	switch strings.ToLower(version) {
	case "all":
		os.RemoveAll(windowsClient.PythonRootContainer)
	default:
		uninstallSingleVersion(version)
	}
}
