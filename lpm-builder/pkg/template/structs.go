package template

import (
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
