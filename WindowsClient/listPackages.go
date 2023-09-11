package windowsClient

import (
	appUtils "AppUtils"
	pythonVersion "PythonVersion"

	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func (client *Client) ListLatest() {
	client.fetchAllAvailableVersions()

	fmt.Println("Latest python versions:")
	fmt.Println("(first 5 for each major version)")
	fmt.Println(" ")

	stableVersions := make([]pythonVersion.PythonVersion, 0)
	for _, version := range client.PythonVersions.Classes {
		if version.IsStable {
			stableVersions = append(stableVersions, *version)
		}
	}

	// Sort by release date
	sort.SliceStable(stableVersions, func(i, j int) bool {
		return stableVersions[i].ReleaseDateInt >= stableVersions[j].ReleaseDateInt
	})

	// Only latest 5 releases for each major version will be printed
	python3 := make([]pythonVersion.PythonVersion, 0)
	python2 := make([]pythonVersion.PythonVersion, 0)

	for _, version := range stableVersions {
		if version.VersionInfo.Major() == 3 {
			if len(python3) == 5 {
				continue
			}
			python3 = append(python3, version)
		} else if version.VersionInfo.Major() == 2 {
			if len(python2) == 5 {
				continue
			}
			python2 = append(python2, version)
		}
	}

	// Print Python 3 and Python 2 versions
	for _, p := range python3 {
		fmt.Println(p.String())
	}

	for _, p := range python2 {
		fmt.Println(p.String())
	}

	fmt.Println("\nThis is a partial list. For a complete list, visit https://www.python.org/downloads/")
}

func (client *Client) ListAll() {
	client.fetchAllAvailableVersions()

	fmt.Println("All python versions:")
	fmt.Println("(First 20 of the list)")
	fmt.Println(" ")

	// limit := 20
	limit := len(client.PythonVersions.All) // or 20 if need to be more concise

	for i := 0; i < limit; i++ {
		current := client.PythonVersions.All[i]
		fmt.Println(client.PythonVersions.Classes[current].String())
	}

	fmt.Println("\nThis is a partial list. For a complete list, visit https://www.python.org/downloads/")
}

func (client *Client) ListInstalled() {
	var installed int

	installationDir := filepath.Join(client.AppRoot, PythonRootContainer)
	versionInUse := appUtils.GetPythonVersionInUse()

	entries, err := os.ReadDir(installationDir)

	for _, entry := range entries {
		if entry.IsDir() {
			dirname := entry.Name()

			data, err := os.ReadFile(filepath.Clean(filepath.Join(installationDir, dirname, "version")))
			if err != nil {
				log.Fatalln(err)
			}

			if strings.Contains(versionInUse, string(data)) {
				fmt.Println(dirname + " (in use)")
			} else {
				fmt.Println(dirname)
			}

			installed++
		}
	}

	if err != nil || installed == 0 {
		fmt.Println("No installed version found.")
	}
}
