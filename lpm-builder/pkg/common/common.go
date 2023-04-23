package common

import (
	"fmt"
	"os"
)

type Version struct {
	ReadableFormat string  `json:"readable_format"`
	Major          int     `json:"major"`
	Minor          int     `json:"minor"`
	Patch          int     `json:"patch"`
	Tag            *string `json:"tag"`
}

type Dependency struct {
	Name    string  `json:"name"`
	Version Version `json:"version"`
}

type Template struct {
	Name                string       `json:"name"`
	Description         string       `json:"description"`
	Maintainer          string       `json:"maintainer"`
	SourcePkg           string       `json:"source_pkg"`
	Repository          string       `json:"repository"`
	Homepage            string       `json:"homepage"`
	Arch                string       `json:"arch"`
	Kind                string       `json:"kind"`
	Tags                []string     `json:"tags"`
	Version             Version      `json:"version"`
	License             string       `json:"license"`
	RuntimeDependencies []Dependency `json:"runtime_dependencies"`
	RuntimeSuggestions  []Dependency `json:"runtime_suggestions"`
	BuildDependencies   []Dependency `json:"build_dependencies"`
}

func SetReadableVersion(version *Version) {
	readable_format := fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)

	if version.Tag != nil {
		readable_format += fmt.Sprintf("-%s", *version.Tag)
	}

	version.ReadableFormat = readable_format

}

func FailOnError(err error, msg string) {
	if err != nil {
		fmt.Println("Error: ", err, msg)
		os.Exit(1)
	}
}
