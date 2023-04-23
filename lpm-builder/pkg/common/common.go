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

func FailOnError(err error, msg string) {
	if err != nil {
		fmt.Println("Error: ", err, msg)
		os.Exit(1)
	}
}

func SetReadableVersion(version *Version) {
	readable_format := fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)

	if version.Tag != nil {
		readable_format += fmt.Sprintf("-%s", *version.Tag)
	}

	version.ReadableFormat = readable_format

}
