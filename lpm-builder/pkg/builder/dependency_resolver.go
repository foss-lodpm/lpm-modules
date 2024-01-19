package builder

import (
	"fmt"
	"lpm_builder/pkg/common"
	"lpm_builder/pkg/template"
	"os/exec"
)

func InstallBuildTimeDependencies(build *template.Build) {
	common.Logger.Printf("Resolving build time dependencies..")

	for _, value := range build.MandatoryDependencies.Build {
		pkg_with_version := fmt.Sprintf("%s@%s%s", value.Name, value.Version.Condition, value.Version.ReadableFormat)

		common.Logger.Printf("Installing build time dependency %s", pkg_with_version)
		cmd := exec.Command("sudo", "lpm", "--install", pkg_with_version)
		out, err := cmd.CombinedOutput()
		if err != nil {
			common.Logger.Print("\n\n")
			common.FailOnError(err, "Failed installing build time dependency.\n"+string(out))
		}
		common.Logger.Printf("%s", out)
	}
}
