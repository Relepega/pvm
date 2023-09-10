package windowsClient

import (
	utils "AppUtils"
	pythonVersion "PythonVersion"

	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/vmihailenco/msgpack/v5"
)

func (client *Client) parsePythonPackages(items []pythonVersion.PackagesCatalog) {
	for _, version := range items {
		versionNumber := version.PackageEntry.Version
		packageContent := version.PackageEntry.PackageContent
		releaseDate, _ := time.Parse(time.RFC3339, version.PackageEntry.Published)

		if len(strings.Split(versionNumber, ".")) > 3 {
			return
		}

		pv, err := semver.NewVersion(versionNumber)

		if err != nil {
			panic(err.Error() + ": " + versionNumber)
		}

		client.PythonVersions.All = append([]string{versionNumber}, client.PythonVersions.All...)

		if pv.Prerelease() == "" {
			client.PythonVersions.Stable = append([]string{versionNumber}, client.PythonVersions.Stable...)
		} else {
			client.PythonVersions.Unstable = append([]string{versionNumber}, client.PythonVersions.Unstable...)
		}

		vnNoDash := strings.ReplaceAll(versionNumber, "-", "")
		baseVersion := fmt.Sprintf("%d.%d.%d", pv.Major(), pv.Minor(), pv.Patch())

		arch := ".amd64"
		if runtime.GOARCH == "386" {
			arch = ""
		}

		installerFilename := fmt.Sprintf("python-%s%s.msi", vnNoDash, arch)
		splittedPackageContentStr := strings.Split(packageContent, "/")
		if pv.Major() >= 3 {
			installerFilename = fmt.Sprintf("%s.zip", splittedPackageContentStr[len(splittedPackageContentStr)-1])
		}

		downloadUrl := fmt.Sprintf("https://www.python.org/ftp/python/%s/%s", baseVersion, installerFilename)
		if pv.Major() >= 3 {
			downloadUrl = packageContent
		}

		y, m, d := releaseDate.Date()

		client.PythonVersions.Classes[versionNumber] = pythonVersion.NewPythonVersion(
			versionNumber,
			fmt.Sprintf("%d/%d/%d", m, d, y),
			releaseDate.Unix(),
			downloadUrl,
			installerFilename,
		)
	}
}

func (client *Client) fetchAllAvailableVersions() {
	if client.CachedDataExists {
		return
	}

	cacheFile := filepath.Join(client.AppRoot, "cache_"+client.Arch)

	if f, err := os.Stat(cacheFile); err == nil {
		if f.IsDir() {
			log.Fatalf("Cache file can't be a directory.")
		}

		data, err := os.ReadFile(cacheFile)
		if err != nil {
			log.Fatal(err)
		}

		var cacheData PythonVersions
		err = msgpack.Unmarshal(data, &cacheData)
		if err != nil {
			log.Fatal(err)
		}

		var expiry int64 = 2 * 60 * 60 * 1000 // 2 hours
		if client.ExpiryDate+expiry < time.Now().Unix() {
			client.PythonVersions = cacheData
			return
		}
	}

	nugetPackages := []string{"python2", "python"}
	if client.Arch == "win32" {
		nugetPackages = []string{"python2x86", "pythonx86"}
	}

	var wg sync.WaitGroup
	for _, packageID := range nugetPackages {
		wg.Add(1)
		go func(packageID string) {
			defer wg.Done()

			url := fmt.Sprintf("https://api.nuget.org/v3/registration5-semver1/%s/index.json", packageID)
			body := utils.FetchJson(url, client.HttpClient)

			if strings.Contains(packageID, "python2") {
				var data pythonVersion.Python2ApiSchema

				err := json.Unmarshal(body, &data)
				if err != nil {
					log.Fatal(err)
				}

				catalogs := data.Items
				client.parsePythonPackages(catalogs[0].Items)
			} else {
				var data pythonVersion.Python3ApiSchema
				err := json.Unmarshal(body, &data)
				if err != nil {
					log.Fatal(err)
				}

				for _, item := range data.Items {
					paginationContainer := utils.FetchJson(item.ID, client.HttpClient)

					// var paginationElements map[string]interface{}
					var paginationElements pythonVersion.Python2CatalogItem
					err := json.Unmarshal(paginationContainer, &paginationElements)
					if err != nil {
						log.Fatal(err)
					}

					client.parsePythonPackages(paginationElements.Items)
				}
			}
		}(packageID)
	}
	wg.Wait()

	client.PythonVersions.CreationDate = time.Now().Unix()

	// create file if not exists
	file, err := os.OpenFile(cacheFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	// convert data to []byte
	encoded, err := msgpack.Marshal(client.PythonVersions)
	if err != nil {
		log.Fatalf("Error while encoding data to bytes: %v", err)
	}

	// write bytes to file
	_, err = file.Write(encoded)
	if err != nil {
		log.Fatalf("Error writing file: %v", err)
	}

	client.CachedDataExists = true
}
