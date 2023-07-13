package PythonVersion

/*

	+---------------------------+
	|    PYTHON 2 API SCHEMA    |
	+---------------------------+

*/

type Python2ApiSchema struct {
	ID              string               `json:"@id"`
	Type            []string             `json:"@type"`
	CommitID        string               `json:"commitId"`
	CommitTimeStamp string               `json:"commitTimeStamp"`
	Count           int                  `json:"count"`
	Items           []Python2CatalogItem `json:"items"`
	Context         interface{}          `json:"@context"`
}

type Python2CatalogItem struct {
	ID              string            `json:"@id"`
	Type            string            `json:"@type"`
	CommitID        string            `json:"commitId"`
	CommitTimeStamp string            `json:"commitTimeStamp"`
	Count           int               `json:"count"`
	Items           []PackagesCatalog `json:"items"`
	Parent          string            `json:"parent"`
	Lower           string            `json:"lower"`
	Upper           string            `json:"upper"`
}

type PackageEntry struct {
	ID                       string   `json:"@id"`
	Type                     string   `json:"@type"`
	Authors                  string   `json:"authors"`
	Description              string   `json:"description"`
	IconUrl                  string   `json:"iconUrl"`
	NugetPackageID           string   `json:"id"`
	Language                 string   `json:"language"`
	LicenseExpression        string   `json:"licenseExpression"`
	LicenseUrl               string   `json:"licenseUrl"`
	Listed                   bool     `json:"listed"`
	MinClientVersion         string   `json:"minClientVersion"`
	PackageContent           string   `json:"packageContent"`
	ProjectUrl               string   `json:"projectUrl"`
	Published                string   `json:"published"`
	RequireLicenseAcceptance bool     `json:"requireLicenseAcceptance"`
	Summary                  string   `json:"summary"`
	Tags                     []string `json:"tags"`
	Title                    string   `json:"title"`
	Version                  string   `json:"version"`
}

/*

	+---------------------------+
	|    PYTHON 3 API SCHEMA    |
	+---------------------------+

*/

type Python3ApiSchema struct {
	ID              string               `json:"@id"`
	Type            []string             `json:"@type"`
	CommitID        string               `json:"commitId"`
	CommitTimeStamp string               `json:"commitTimeStamp"`
	Count           int                  `json:"count"`
	Items           []Python3CatalogItem `json:"items"`
	Context         interface{}          `json:"@context"`
}

type Python3CatalogItem struct {
	ID              string `json:"@id"`
	Type            string `json:"@type"`
	CommitID        string `json:"commitId"`
	CommitTimeStamp string `json:"commitTimeStamp"`
	Count           int    `json:"count"`
	Lower           string `json:"lower"`
	Upper           string `json:"upper"`
}

/*

	+--------------------------------------+
	|    PYTHON 3 PAGINATION API SCHEMA    |
	+--------------------------------------+

*/

type Python3PaginationSchema struct {
	ID              string            `json:"@id"`
	Type            string            `json:"@type"`
	CommitID        string            `json:"commitId"`
	CommitTimeStamp string            `json:"commitTimeStamp"`
	Count           int               `json:"count"`
	Items           []PackagesCatalog `json:"catalogEntry"`
	Parent          string            `json:"parent"`
	Lower           string            `json:"lower"`
	Upper           string            `json:"upper"`
	Context         interface{}       `json:"@context"`
}

/*

	+----------------+
	|    GENERICS    |
	+----------------+

*/

type PackagesCatalog struct {
	ID              string       `json:"@id"`
	Type            string       `json:"@type"`
	CommitID        string       `json:"commitId"`
	CommitTimeStamp string       `json:"commitTimeStamp"`
	PackageEntry    PackageEntry `json:"catalogEntry"`
	PackageContent  string       `json:"packageContent"`
	Registration    string       `json:"registration"`
}
