package builder

import (
	"fmt"
	"lpm_builder/pkg/common"
	"os/exec"
)

func InstallBuildTimeDependencies(ctx *BuilderCtx) {
	common.Logger.Printf("Resolving build time dependencies..")
	for _, value := range ctx.TemplateFields.MandatoryDependencies.Build {
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
