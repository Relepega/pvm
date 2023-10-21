package WindowsClient

import (
	"AppUtils"
	"PythonVersion"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func UseVersion(version string) (*PythonVersion.PythonVersion, error) {
	client := NewClient()

	client.fetchAllAvailableVersions()

	ver := strings.ToLower(version)

	var pv PythonVersion.PythonVersion

	found := false

	if ver == "latest" {
		versionNumber := client.PythonVersions.Stable[0]
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

func UseAlias(installDir string, version *PythonVersion.PythonVersion, alias string) string {
	if strings.ContainsRune(alias, ' ') {
		s1 := "Alias cannot contain whitespaces. use '-' instead."
		s2 := fmt.Sprintf("Example: pvm install 3.11.0 \"%s\"", strings.ReplaceAll(alias, " ", "-"))
		log.Fatalf("%s\n%s", s1, s2)
	}

	if stat, err := os.Stat(installDir); err != nil || !stat.IsDir() {
		os.MkdirAll(installDir, os.FileMode(os.O_RDWR))
	}

	if alias != "" {
		return path.Join(installDir, alias)
	}

	return path.Join(installDir, version.VersionNumber)
}

func (client *Client) InstallNewVersion(version *PythonVersion.PythonVersion, alias string) {
	installPath := UseAlias(client.InstallDir, version, alias)
	installerFp := path.Join(client.AppRoot, version.InstallerFilename)

	stat, err := os.Stat(installPath)
	if err == nil && stat.IsDir() {
		if alias == "" {
			fmt.Printf("Python %s is already installed. Please use the command \"reinstall\" instead.\n", version.VersionNumber)
		} else {
			fmt.Printf("Python installation \"%s\" (aliased as \"%s\") is already installed. Please use the command \"reinstall\" instead.\n", version.VersionNumber, alias)
		}
		return
	}

	fmt.Printf(`Downloading "%s"... `, version.InstallerFilename)

	err = AppUtils.DownloadFile(version.DownloadUrl, installerFp)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Done!")

	// use different install method based on major release
	var pythonPath string

	os.MkdirAll(installPath, os.FileMode(os.O_RDWR))

	if version.VersionInfo.Major() == 2 {
		pythonPath = client.python2Install(version, installPath, installerFp)
	} else {
		pythonPath = client.python3Install(version, installPath, installerFp)
	}

	if pythonPath == "" {
		log.Fatalf("Python version not installed correctly, try again...")
	}

	// create "version" file to not mismatch the version with the parent folder name
	versionFileName := filepath.Join(pythonPath, "version")
	f, err := os.OpenFile(versionFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeDevice)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	f.WriteString(version.VersionNumber)

	fmt.Print(`Cleaning up... `)

	err = os.Remove(installerFp)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Done!")

	client.MakeSymlink(version.VersionNumber, pythonPath)

	fmt.Printf("Python %s installed successfully!\n", version.VersionNumber)
}

func (client *Client) python3Install(version *PythonVersion.PythonVersion, installPath string, installerFp string) string {
	pipInstallationScriptFilepath := filepath.Join(installPath, version.PipVersion.Filename)
	pythonVersionBasename := fmt.Sprintf("python%d%d", version.VersionInfo.Major(), version.VersionInfo.Minor())

	// fmt.Print("Sorting files and fixing bugs... ")

	AppUtils.UnzipFile(installerFp, installPath)

	// zip all 'Lib' content apart 'site-packages' to 'pythonXXX.zip'
	dirToZip := filepath.Join(installPath, "Lib")
	ouputFilePath := filepath.Join(installPath, pythonVersionBasename+".zip")
	excludedFiles := []string{"site-packages"}
	AppUtils.ZipDirWithExclusions(dirToZip+"\\", ouputFilePath, excludedFiles)

	// remove all directories from 'Lib' apart 'site-packages'
	libDir := filepath.Join(installPath, "Lib")
	files, err := os.ReadDir(libDir)
	if err != nil {
		fmt.Println(" ")
		log.Fatalln(err)
	}

	for _, f := range files {
		if strings.Contains(f.Name(), "site-packages") {
			continue
		}

		fp := filepath.Join(libDir, f.Name())
		os.RemoveAll(fp)
	}

	// fix site-packages (https://stackoverflow.com/a/68891090)
	pthFileFilepath := filepath.Join(installPath, pythonVersionBasename+"._pth")
	f, err := os.Create(pthFileFilepath)
	if err != nil {
		fmt.Println(" ")
		log.Fatalln(err)
	}
	defer f.Close()

	f.WriteString(pythonVersionBasename + ".zip" + "\n" + ".\n" + "\n" + "# Uncomment to run site.main() automatically\n" + "#import site\n" + "\n" + "Lib\\site-packages")

	// move all the files from 'DLLs' to '{installPath}'
	dllsDir := filepath.Join(installPath, "DLLs")
	files, err = os.ReadDir(dllsDir)
	if err != nil {
		fmt.Println(" ")
		log.Fatalln(err)
	}

	for _, f := range files {
		fp := filepath.Join(dllsDir, f.Name())
		newFp := filepath.Join(installPath, f.Name())
		os.Rename(fp, newFp)
	}

	// delete 'DLLs' directory
	os.RemoveAll(dllsDir)

	// download "get-pip.py" if not already downloaded
	_, err = os.Stat(pipInstallationScriptFilepath)

	if err != nil {
		fmt.Printf("Downloading \"get-pip.py\" from \"%s\" ...\n", version.PipVersion.DownloadUrl)
		AppUtils.DownloadFile(version.PipVersion.DownloadUrl, pipInstallationScriptFilepath)
	}

	// enable pip functionality and fixing the issue https://github.com/pypa/pip/issues/5292
	fmt.Println("Installing \"pip\" package... ")

	fmt.Print(AppUtils.CmdColors["Orange"])

	pythonExe := filepath.Join(installPath, "python.exe")
	fmt.Println("Python.exe path: " + pythonExe)
	command := fmt.Sprintf(`%s %s --no-warn-script-location`, pythonExe, pipInstallationScriptFilepath)

	malfunctioningVersions := []string{"3.5.2", "3.5.2.1", "3.5.2.2", "3.6.0"}

	for _, s := range malfunctioningVersions {
		if s == version.VersionNumber {
			command = fmt.Sprintf(`%s -m easy_install pip easy_install`, pythonExe)
			break
		}
	}

	command += fmt.Sprintf(` && %s -m pip install --upgrade pip`, pythonExe)

	_, err = AppUtils.RunCmd(command)

	fmt.Print(AppUtils.CmdColors["Reset"])

	if err != nil {
		fmt.Println(" ")
		log.Fatalln(err)
	}

	fmt.Println("Done!")

	// return absolute path to newly installed python dir
	returnValue, _ := filepath.Abs(installPath)
	return returnValue
}

func (client *Client) python2Install(version *PythonVersion.PythonVersion, installPath string, installerFp string) string {
	fmt.Print("Unpacking installer data... ")

	absinstallerFp, _ := filepath.Abs(installerFp)
	absinstallPath, _ := filepath.Abs(installPath)

	command := fmt.Sprintf(`msiexec /n /a %s /qn TARGETDIR=%s`, absinstallerFp, absinstallPath)
	_, err := AppUtils.RunCmd(command)
	if err != nil {
		fmt.Println(" ")
		log.Fatalln(err)
		// fmt.Println("command: " + command)
		log.Fatalln("Couldn't unpack the requested data. Aborting...")
	}

	fmt.Println("Done!")

	// move all the files into "DLLs" to {installPath}
	fmt.Print("Sorting files... ")

	dllsPath := path.Join(installPath, "DLLs")
	files, err := os.ReadDir(dllsPath)
	if err != nil {
		fmt.Println(" ")
		panic(err)
	}

	for _, file := range files {
		oldPath := path.Join(dllsPath, file.Name())
		newPath := path.Join(installPath, file.Name())
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

	pythonExe, _ := filepath.Abs(path.Join(installPath, "python.exe"))
	command = fmt.Sprintf(`%s -m ensurepip --default-pip && %s -m pip install --upgrade pip`, pythonExe, pythonExe)
	_, _ = AppUtils.RunCmd(command)

	fmt.Print(AppUtils.CmdColors["Reset"])
	// if err != nil {
	// 	fmt.Println(" ")
	// 	log.Fatalln(err)
	// }

	fmt.Println("Done!")

	return absinstallPath
}
