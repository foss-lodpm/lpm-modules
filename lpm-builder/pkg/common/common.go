package common

import (
	"fmt"
	"log"
)

var Logger *log.Logger

type System struct {
	BuilderVersion         Version `json:"builder_version"`
	MinSupportedLpmVersion Version `json:"min_supported_lpm_version"`
}

type Version struct {
	ReadableFormat string  `json:"readable_format"`
	Major          uint    `json:"major"`
	Minor          uint    `json:"minor"`
	Patch          uint    `json:"patch"`
	Tag            *string `json:"tag"`
}

type Dependency struct {
	Name    string  `json:"name"`
	Version Version `json:"version"`
}

var BuilderVersion = Version{
	ReadableFormat: "1.0.0-beta",
	Major:          1,
	Minor:          0,
	Patch:          0,
	Tag:            StringPtr("beta"),
}

var MinSupportedLpmVersion = Version{
	ReadableFormat: "0.0.1-alpha",
	Major:          0,
	Minor:          0,
	Patch:          1,
	Tag:            StringPtr("alpha"),
}

func StringPtr(s string) *string {
	return &s
}

func FailOnError(err error, v ...any) {
	if err != nil {
		log.Fatal("Error: ", err, fmt.Sprint(v...))
	}
}

func FatalError(v ...any) {
	log.Fatal(fmt.Sprint(v...))
}

func SetReadableVersion(version *Version) {
	readable_format := fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)

	if version.Tag != nil {
		readable_format += fmt.Sprintf("-%s", *version.Tag)
	}

	version.ReadableFormat = readable_format

}

func Contains[T comparable](s []T, e T) bool {
    for _, v := range s {
        if v == e {
            return true
        }
    }
    return false
}
