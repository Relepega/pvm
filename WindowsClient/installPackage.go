package windowsClient

import (
	utils "AppUtils"
	pythonVersion "PythonVersion"

	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func UseVersion(version string) (*pythonVersion.PythonVersion, error) {
	client := NewClient()

	ver := strings.ToLower(version)

	var pv pythonVersion.PythonVersion

	found := false

	if ver == "latest" {
		length := len(client.PythonVersions.Stable)
		versionNumber := client.PythonVersions.Stable[length-1]
		pv = *client.PythonVersions.Classes[versionNumber]
		found = true
	} else {
		for _, v := range client.PythonVersions.All {
			if v == version {
				found = true
				pv = *client.PythonVersions.Classes[version]
			}
		}
	}

	client = nil

	if !found {
		e := fmt.Errorf("\"%s\" is not a valid python version", version)
		return nil, e
	}

	return &pv, nil

}

// func UseAlias(version string, alias string) string {
// 	if alias != "" {
// 		return path.Join(PythonRootContainer, alias)
// 	}

// 	return path.Join(PythonRootContainer, version)
// }

func (client *Client) InstallNewVersion(version *pythonVersion.PythonVersion, alias string) {
	client.fetchAllAvailableVersions()

	unpackedPythonPath := path.Join(PythonRootContainer, version.VersionNumber)
	offlineFilePath := path.Join(client.AppRoot, version.InstallerFilename)

	// https://stackoverflow.com/a/40624033
	if stat, err := os.Stat(unpackedPythonPath); err != nil || stat.IsDir() {
		fmt.Printf("Python %s is already installed. Please use the command \"reinstall\" instead.", version.VersionNumber)
	}

	// if _, err := os.Stat(unpackedPythonPath); !os.IsNotExist(err) {
	// 	err := os.RemoveAll(unpackedPythonPath)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	fmt.Printf(`Downloading "%s"... `, version.InstallerFilename)

	err := utils.DownloadFile(version.DownloadUrl, offlineFilePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done!")

	// use different install method based on major release
	var pythonPath string

	if version.VersionInfo.Major() == 2 {
		pythonPath = client.python2Install(version, unpackedPythonPath, offlineFilePath)
	} else {
		finalPythonPath := unpackedPythonPath
		unpackedPythonPath := unpackedPythonPath + "temp"
		pythonPath = client.python3Install(version, unpackedPythonPath, finalPythonPath, offlineFilePath)
	}

	if pythonPath == "" {
		log.Fatalln("Python version not installed correctly, try again...")
	}

	// // create "version" file to not mismatch the version with the parent folder name
	// versionFileName := filepath.Join(pythonPath, "version")
	// f, err := os.OpenFile(versionFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeDevice)

	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer f.Close()

	// f.WriteString(version.VersionNumber)

	fmt.Print(`Cleaning up... `)

	err = os.Remove(offlineFilePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done!")

	client.MakeSymlink(version.VersionNumber, pythonPath)

	fmt.Printf("Python %s installed successfully!\n", version.VersionNumber)
}

func (client *Client) python2Install(version *pythonVersion.PythonVersion, unpackedPythonPath string, offlineFilePath string) string {
	fmt.Print("Unpacking installer data... ")

	absOfflineFilePath, _ := filepath.Abs(offlineFilePath)
	absUnpackedPythonPath, _ := filepath.Abs(unpackedPythonPath)

	command := fmt.Sprintf(`msiexec /n /a %s /qn TARGETDIR=%s`, absOfflineFilePath, absUnpackedPythonPath)
	_, err := utils.RunCmd(command)

	if err != nil {
		fmt.Println(" ")
		log.Fatal(err)
		// fmt.Println("command: " + command)
		log.Fatal("Couldn't unpack the requested data. Aborting...")
	}

	fmt.Println("Done!")

	// move all the files into "DLLs" to unpackedPythonPath'
	fmt.Print("Sorting files... ")

	dllsPath := path.Join(unpackedPythonPath, "DLLs")
	files, err := os.ReadDir(dllsPath)
	if err != nil {
		fmt.Println(" ")
		panic(err)
	}

	for _, file := range files {
		oldPath := path.Join(dllsPath, file.Name())
		newPath := path.Join(unpackedPythonPath, file.Name())
		err := os.Rename(oldPath, newPath)
		if err != nil {
			fmt.Println(" ")
			panic(err)
		}
	}

	// delete 'DLLs' directory
	os.RemoveAll(dllsPath)

	fmt.Println("Done!")

	// enable pip functionality
	fmt.Println(`Installing "pip" package... `)

	pythonExe, _ := filepath.Abs(path.Join(unpackedPythonPath, "python.exe"))
	command = fmt.Sprintf(`%s -m ensurepip --default-pip && %s -m pip install --upgrade pip`, pythonExe, pythonExe)
	_, err = utils.RunCmd(command)

	fmt.Print(utils.CmdColors["Reset"])
	if err != nil {
		fmt.Println(" ")
		log.Fatal(err)
	}

	fmt.Println("Done!")

	return absUnpackedPythonPath
}

func (client *Client) python3Install(version *pythonVersion.PythonVersion, unpackedPythonPath string, finalPythonPath string, offlineFilePath string) string {
	pipInstallationScriptFilepath := filepath.Join(finalPythonPath, "Tools", version.PipVersion.Filename)
	pythonVersionBasename := fmt.Sprintf("python%d%d", version.VersionInfo.Major(), version.VersionInfo.Minor())

	fmt.Print("Sorting files and fixing bugs... ")

	utils.UnzipFile(offlineFilePath, unpackedPythonPath)

	// move 'tools' folder outside the temp dir and make it the final one
	os.Rename(filepath.Join(unpackedPythonPath, "tools"), finalPythonPath)

	// remove temp dir
	os.RemoveAll(unpackedPythonPath)

	// zip all 'Lib' content apart 'site-packages' to 'pythonXXX.zip'
	dirToZip := filepath.Join(finalPythonPath, "Lib")
	ouputFilePath := filepath.Join(finalPythonPath, pythonVersionBasename+".zip")
	excludedFiles := []string{"site-packages"}
	utils.ZipDirWithExclusions(dirToZip+"\\", ouputFilePath, excludedFiles)

	// remove all directories from 'Lib' apart 'site-packages'
	libDir := filepath.Join(finalPythonPath, "Lib")
	files, err := os.ReadDir(libDir)
	if err != nil {
		fmt.Println(" ")
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.Contains(f.Name(), "site-packages") {
			continue
		}

		fp := filepath.Join(libDir, f.Name())
		os.RemoveAll(fp)
	}

	// fix site-packages (https://stackoverflow.com/a/68891090)
	pthFileFilepath := filepath.Join(finalPythonPath, pythonVersionBasename+"._pth")
	f, err := os.Create(pthFileFilepath)

	if err != nil {
		fmt.Println(" ")
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString(pythonVersionBasename + ".zip" + "\n" + ".\n" + "\n" + "# Uncomment to run site.main() automatically\n" + "#import site\n" + "\n" + "Lib\\site-packages")

	// move all the 'DLLs' files to '{unpackedPythonPath}'
	dllsDir := filepath.Join(finalPythonPath, "DLLs")
	files, err = os.ReadDir(dllsDir)
	if err != nil {
		fmt.Println(" ")
		log.Fatal(err)
	}

	for _, f := range files {
		fp := filepath.Join(dllsDir, f.Name())
		newFp := filepath.Join(finalPythonPath, f.Name())
		os.Rename(fp, newFp)
	}

	// delete 'DLLs' directory
	os.RemoveAll(dllsDir)
	fmt.Println("Done!")

	// download "get-pip.py" if not already downloaded
	_, err = os.Stat(pipInstallationScriptFilepath)

	if err != nil {
		fmt.Printf("Downloading \"get-pip.py\" from \"%s\" ...\n", version.PipVersion.DownloadUrl)
		utils.DownloadFile(version.PipVersion.DownloadUrl, pipInstallationScriptFilepath)
	}

	// enable pip functionality and fixing the issue https://github.com/pypa/pip/issues/5292
	fmt.Println("Installing \"pip\" package... ")

	fmt.Print(utils.CmdColors["Orange"])

	pythonExe := filepath.Join(finalPythonPath, "python.exe")
	fmt.Println("Python.exe path: " + pythonExe)
	command := fmt.Sprintf(`%s %s`, pythonExe, pipInstallationScriptFilepath)

	malfunctioningVersions := []string{"3.5.2", "3.5.2.1", "3.5.2.2", "3.6.0"}

	for _, s := range malfunctioningVersions {
		if s == version.VersionNumber {
			command = fmt.Sprintf(`%s -m easy_install pip easy_install`, pythonExe)
			break
		}
	}

	command += fmt.Sprintf(` && %s -m pip install --upgrade pip`, pythonExe)

	_, err = utils.RunCmd(command)

	fmt.Print(utils.CmdColors["Reset"])

	if err != nil {
		fmt.Println(" ")
		log.Fatal(err)
	}

	fmt.Println("Done!")

	// return absolute path to newly installed python dir
	returnValue, _ := filepath.Abs(finalPythonPath)
	return returnValue
}
