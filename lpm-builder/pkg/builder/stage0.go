package builder

import (
	"fmt"
	"lpm_builder/pkg/common"
	"os"
	"os/exec"
	"path/filepath"
)

var BuiltInFunctions string = `
function validate_checksum {
	# Get the file path and checksum from the arguments
	file_path="$1"
	expected_checksum="$2"

	# Calculate the actual checksum of the file
	actual_checksum="$(sha256sum "$file_path" | awk '{print $1}')"

	# Compare the actual and expected checksums
	if [[ "$actual_checksum" == "$expected_checksum" ]]; then
		echo "Checksum validation successful"
		return 0
	else
		echo "Checksum validation failed"
		rm "$1"
		exit 1
	fi
}

function install_to_package {
	src_file="$1"
	target="$2"

	install -D $src_file program/$target
}
`

const (
	Init             = "init"
	Build            = "build"
	InstallFiles     = "install_files"
	PostInstallFiles = "post_install_files"
)

func PrepareScript(stage0Path string, script string) string {
	s := filepath.Join(stage0Path, script)
	preparedScript := fmt.Sprintf(`
		#!/bin/bash
		set -e

		%s

		target_script=$(<%s)
		eval "$target_script"
	`, BuiltInFunctions, s)

	return preparedScript
}

func execute(scriptPath string, script string, executeIn string) {
	common.Logger.Printf("Executing stage0/%s script", script)

	cmd := exec.Command("/bin/bash", "-c", PrepareScript(scriptPath, script))
	cmd.Dir = executeIn
	_, err := cmd.Output()
	common.FailOnError(err, "Couldn't execute "+script+" script from template directory.")

}

func executeStage0(ctx *BuilderCtx) {
	execute(ctx.Stage0ScriptsDir, Init, ctx.TmpPkgDir)
	execute(ctx.Stage0ScriptsDir, Build, ctx.TmpSrcDir)
	execute(ctx.Stage0ScriptsDir, InstallFiles, ctx.TmpSrcDir)

	// since post_install_files is optional, check if it exists before the execution
	if _, err := os.Stat(filepath.Join(ctx.Stage0ScriptsDir, PostInstallFiles)); err == nil {
		execute(ctx.Stage0ScriptsDir, PostInstallFiles, ctx.TmpPkgDir)
	}
}
