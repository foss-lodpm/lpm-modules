package template

import (
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
		Kind:             "",
		Tags:             []string{},
		License:          "",
		Builds:           make(map[string]*Build),
	}

	for _, build := range template.Builds {
		common.SetReadableVersion(&build.Version)
	}

	template_json, err := common.Utf8FriendlyJsonMarshal(template)
	common.FailOnError(err, "Failed on serializing template to json string")

	err = os.WriteFile(out_path+"template.json", template_json, 0644)
	common.FailOnError(err, "Failed on writing to file")
}
