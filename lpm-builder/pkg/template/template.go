package template

import (
	"encoding/json"
	common "lpm_builder/pkg/common"
	"os"

	"github.com/urfave/cli/v2"
)

func CreateDefault(c *cli.Context, template_name string, out_path string) {
	template := &Template{
		Name:             template_name,
		Description:      "",
		Maintainer:       "",
		SourceRepository: "",
		Homepage:         "",
		Arch:             "",
		Kind:             "",
		FileChecksumAlgo: "",
		Tags:             []string{},
		Version:          common.Version{},
		License:          "",
		MandatoryDependencies: Dependencies{
			Build:   []common.Dependency{},
			Runtime: []common.Dependency{},
		},
		SuggestedDependencies: Dependencies{
			Build:   []common.Dependency{},
			Runtime: []common.Dependency{},
		},
	}

	common.SetReadableVersion(&template.Version)

	template_json, err := json.MarshalIndent(template, "", "\t")
	common.FailOnError(err, "Failed on serializing template to json string")

	err = os.WriteFile(out_path+"template.json", template_json, 0644)
	common.FailOnError(err, "Failed on writing to file")
}
