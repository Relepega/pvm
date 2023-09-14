package PythonVersion

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

type PipVersion struct {
	Version     string
	Filename    string
	DownloadUrl string
}

type PythonVersion struct {
	VersionNumber     string
	VersionInfo       *semver.Version
	IsStable          bool
	ReleaseDate       string
	ReleaseDateInt    int64
	DownloadUrl       string
	InstallerFilename string
	PipVersion        PipVersion
}

func NewPythonVersion(versionNumber string, releaseDate string, releaseDateInt int64, downloadUrl string, installerFilename string) *PythonVersion {
	v, _ := semver.NewVersion(versionNumber)

	isStableStr := v.Prerelease()
	isStable := false

	if isStableStr == "" {
		isStable = true
	}

	pipVersion := getPipVersion(v)

	return &PythonVersion{
		VersionNumber:     versionNumber,
		VersionInfo:       v,
		IsStable:          isStable,
		ReleaseDate:       releaseDate,
		ReleaseDateInt:    releaseDateInt,
		DownloadUrl:       downloadUrl,
		InstallerFilename: installerFilename,
		PipVersion:        pipVersion,
	}
}

func (pythonVersion *PythonVersion) String() string {
	return fmt.Sprintf("Python %d, version %s, released in date %s", pythonVersion.VersionInfo.Major(), pythonVersion.VersionNumber, pythonVersion.ReleaseDate)
}

func (pythonVersion *PythonVersion) GetPipVersion() PipVersion {
	return pythonVersion.PipVersion
}

func getPipVersion(v *semver.Version) PipVersion {
	s := semver.MustParse("3.7.0")
	if v.LessThan(s) {
		majorMinor := fmt.Sprintf("%d.%d", v.Major(), v.Minor())
		return PipVersion{
			Version:     majorMinor,
			Filename:    fmt.Sprintf("get-pip-v%s.py", majorMinor),
			DownloadUrl: fmt.Sprintf("https://bootstrap.pypa.io/pip/%s/get-pip.py", majorMinor),
		}
	} else {
		return PipVersion{
			Version:     "latest",
			Filename:    "get-pip-latest.py",
			DownloadUrl: "https://bootstrap.pypa.io/pip/get-pip.py",
		}
	}
}
