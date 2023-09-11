package AppUtils

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

func RunCmd(command string) (cmd *exec.Cmd, err error) {
	osShell := []string{"bash", "-c"}

	if runtime.GOOS == "windows" {
		osShell = []string{"cmd", "/C"}
	}

	cmd = exec.Command(osShell[0], osShell[1], command)
	// cmd.Stdout = os.Stdout
	cmd.Stdout = nil
	cmd.Stderr = os.Stdout
	err = cmd.Run()

	if err != nil {
		return cmd, err
	}

	return cmd, nil
}

func getModuleFileName(hModule windows.Handle) (string, error) {
	var buffer [windows.MAX_PATH]uint16
	_, err := windows.GetModuleFileName(hModule, &buffer[0], windows.MAX_PATH)
	if err != nil {
		return "", err
	}
	return syscall.UTF16ToString(buffer[:]), nil
}

func GetWorkingDir() string {
	var appDir string

	if runtime.GOOS == "windows" {
		// Use syscall to get executable path on Windows
		exPath, err := getModuleFileName(0)
		if err != nil {
			panic(err)
		}
		appDir = filepath.Dir(exPath)
	} else {
		// Use os.Args on other platforms
		exPath := os.Args[0]
		appDir = filepath.Dir(exPath)
	}

	return appDir
}

func IsValidFolderName(input string) bool {
	if input == "" {
		return false
	}

	if len(input) < 3 {
		return false
	}

	pattern := `^[a-zA-Z0-9_.\s-]{1,255}$`
	regex, _ := regexp.Compile(pattern)

	return regex.MatchString(input)
}

func FetchJson(url string, httpClient http.Client) []byte {
	res, err := httpClient.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}

func GetPythonVersionInUse() string {
	cmd := exec.Command("python", "-V")
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatalln(err)
		return ""
	}

	output := strings.TrimSpace(string(out))
	res := strings.Split(output, " ")[1]

	return res
}

func DownloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func UnzipFile(src string, dest string) error {
	const filerPrefix string = "tools/"

	stat, err := os.Stat(dest)
	if err != nil || !stat.IsDir() {
		os.Mkdir(dest, os.ModeDir)
	}

	archive, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, f := range archive.File {
		fp := filepath.Join(dest, f.Name)

		// Zip Slip vuln check
		if !strings.HasPrefix(fp, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("\"%s\": illegal file path", fp)
		}

		// filter unwanted files to be extracted
		if !strings.HasPrefix(f.Name, filerPrefix) {
			continue
		}

		fp = filepath.Join(dest, strings.ReplaceAll(f.Name, filerPrefix, ""))

		if f.FileInfo().IsDir() {
			os.MkdirAll(fp, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fp), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(filepath.Clean(fp), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		_, err = io.Copy(outFile, rc)

		if err != nil {
			return err
		}
	}

	return nil
}

func ZipDirWithExclusions(inputDir, outputZip string, exclusions []string) error {
	file, err := os.Create(outputZip)
	if err != nil {
		return err
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		for _, exclusion := range exclusions {
			if strings.Contains(path, exclusion) {
				return nil
			}
		}

		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		f, err := w.Create(strings.ReplaceAll(path, inputDir, ""))
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}

	err = filepath.Walk(inputDir, walker)
	if err != nil {
		return err
	}

	return nil
}
