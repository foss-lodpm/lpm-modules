package template

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	common "lpm_builder/pkg/common"
	"net/url"
	"regexp"
)

type Template struct {
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	Maintainer       string            `json:"maintainer"`
	SourceRepository string            `json:"source_repository"`
	Homepage         string            `json:"homepage"`
	Kind             string            `json:"kind"`
	Tags             []string          `json:"tags"`
	License          string            `json:"license"`
	Builds           map[string]*Build `json:"builds"`
}

type Build struct {
	FileChecksumAlgo      *string        `json:"file_checksum_algo"`
	Version               common.Version `json:"version"`
	MandatoryDependencies Dependencies   `json:"mandatory_dependencies"`
	SuggestedDependencies Dependencies   `json:"suggested_dependencies"`
}

type Dependencies struct {
	Build   []common.Dependency `json:"build"`
	Runtime []common.Dependency `json:"runtime"`
}

func (template *Template) validate() error {
	common.Logger.Println("Validating template fields")

	// field lengths
	{
		if len(template.Description) > 200 {
			return errors.New("description exceeds maximum length of 155 characters")
		}

		if len(template.Tags) > 5 {
			return errors.New("tags exceeds maximum length of 5 items")
		}
	}

	// Package name
	{
		pkgNameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

		if !pkgNameRegex.MatchString(template.Name) {
			return errors.New("Package name can only contain English alphabets, numbers, '-' and '_' characters")
		}
	}

	// Builds
	{
		common.Assert(len(template.Builds) > 0, "You must define at least one build in `builds` field.")
		for build_name, build := range template.Builds {
			var supportedAlgorithms []string = []string{
				"md5",
				"sha256",
				"sha512",
			}

			if build_name == "source" {
				common.Assert(build.FileChecksumAlgo == nil, "`file_checksum_algo` is not supported for source packages.")
			} else {
				common.Assert(common.Contains(supportedAlgorithms, *build.FileChecksumAlgo), fmt.Sprintf("Unsupported checksum algorithm used for '%s' build. Supported algorithms: %v", build_name, supportedAlgorithms))
			}

			var supportedBuilds []string = []string{
				"amd64",
				"source",
				"noarch",
			}

			common.Assert(common.Contains(supportedBuilds, build_name), fmt.Sprintf("Unsupported build '%s'. Supported builds (please contact with maintainers to support more builds): %v", build_name, supportedBuilds))

			if build_name == "source" {
				common.Assert(len(build.MandatoryDependencies.Build) == 0, "source packages can not contain build time dependencies")
				common.Assert(len(build.SuggestedDependencies.Build) == 0, "source packages can not contain build time dependencies")
			}
		}
	}

	// URLs
	{
		if len(template.SourceRepository) > 0 {
			_, err := url.ParseRequestURI(template.SourceRepository)
			common.FailOnError(err, "Invalid URL in source_repository.")
		}

		_, err := url.ParseRequestURI(template.Homepage)
		common.FailOnError(err, "Invalid URL in homepage.")
	}

	return nil
}

func DeserializeTemplate(templateDirPath string) Template {
	common.Logger.Println("Deserializing template fields")

	const templateLeafPath = "/template"

	var template = Template{
		Builds: make(map[string]*Build),
	}

	template_json_content, err := ioutil.ReadFile(templateDirPath + templateLeafPath)
	common.FailOnError(err, "Failed reading template json file")

	err = json.Unmarshal(template_json_content, &template)
	common.FailOnError(err, "Failed reading template json file")

	err = template.validate()
	common.FailOnError(err)

	for _, build := range template.Builds {
		for i := range build.MandatoryDependencies.Runtime {
			if len(build.MandatoryDependencies.Runtime[i].Version.Condition) == 0 {
				build.MandatoryDependencies.Runtime[i].Version.Condition = ">="
			}
		}

		for i := range build.MandatoryDependencies.Build {
			if len(build.MandatoryDependencies.Build[i].Version.Condition) == 0 {
				build.MandatoryDependencies.Build[i].Version.Condition = ">="
			}
		}

		for i := range build.SuggestedDependencies.Runtime {
			if len(build.SuggestedDependencies.Runtime[i].Version.Condition) == 0 {
				build.SuggestedDependencies.Runtime[i].Version.Condition = ">="
			}
		}

		for i := range build.SuggestedDependencies.Build {
			if len(build.SuggestedDependencies.Build[i].Version.Condition) == 0 {
				build.SuggestedDependencies.Build[i].Version.Condition = ">="
			}
		}
	}

	return template
}
