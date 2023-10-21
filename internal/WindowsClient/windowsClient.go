package WindowsClient

import (
	"AppUtils"
	"PythonVersion"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	SymlinkDest          = "%localappdata%\\Python"
	PythonInstallDirname = "Python"
)

type PythonVersions struct {
	All          []string
	Stable       []string
	Unstable     []string
	Classes      map[string]*PythonVersion.PythonVersion
	CreationDate int64
	ExpiryDate   int64
}

type Client struct {
	AppRoot          string
	InstallDir       string
	Arch             string
	HttpClient       http.Client
	PythonVersions   PythonVersions
	CachedDataExists bool
}

func NewClient() *Client {
	appRoot := filepath.Clean(AppUtils.GetWorkingDir())

	arch := "amd64"

	if runtime.GOARCH == "386" {
		arch = "win32"
	}

	httpClient := http.Client{
		Timeout: time.Duration(10 * time.Second),
		Transport: &http.Transport{
			IdleConnTimeout: 60 * time.Second,
		},
	}

	return &Client{
		AppRoot:    appRoot,
		InstallDir: filepath.Join(appRoot, PythonInstallDirname),
		Arch:       arch,
		HttpClient: httpClient,
		PythonVersions: PythonVersions{
			All:          []string{},
			Stable:       []string{},
			Unstable:     []string{},
			Classes:      make(map[string]*PythonVersion.PythonVersion),
			CreationDate: 0,
			ExpiryDate:   0,
		},
		CachedDataExists: false,
	}
}

func (client *Client) ClientInfo() string {
	return fmt.Sprintf("Windows client (%s)", client.Arch)
}

func (client *Client) MakeSymlink(slug string, srcPath string) bool {
	fileInfo, err := os.Stat(srcPath)

	if err != nil && !fileInfo.IsDir() {
		log.Fatalf(`Python installation "%s" hasn't been found. Try installing it first...`, slug)
		return false
	}

	fmt.Print("Making symlink... ")

	err = AppUtils.CreateSymlink(srcPath, SymlinkDest)
	if err != nil {
		print("Couldn't create the symlink, exiting...")
		log.Fatalln(err)
	}

	fmt.Println("Done!")

	return true
}
