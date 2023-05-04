package template

import (
	"encoding/json"
	"io/ioutil"
	common "lpm_builder/pkg/common"
)

type Template struct {
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	Maintainer          string              `json:"maintainer"`
	Repository          string              `json:"repository"`
	Homepage            string              `json:"homepage"`
	Arch                string              `json:"arch"`
	Kind                string              `json:"kind"`
	FileChecksumAlgo    string              `json:"file_checksum_algo"`
	Tags                []string            `json:"tags"`
	Version             common.Version      `json:"version"`
	License             string              `json:"license"`
	RuntimeDependencies []common.Dependency `json:"runtime_dependencies"`
	RuntimeSuggestions  []common.Dependency `json:"runtime_suggestions"`
	BuildDependencies   []common.Dependency `json:"build_dependencies"`
}

func DeserializeTemplate(templateDirPath string) Template {
	const templateLeafPath = "/template"

	var template Template

	template_json_content, err := ioutil.ReadFile(templateDirPath + templateLeafPath)
	common.FailOnError(err, "Failed reading template json file")

	err = json.Unmarshal(template_json_content, &template)
	common.FailOnError(err, "Failed reading template json file")

	return template
}
