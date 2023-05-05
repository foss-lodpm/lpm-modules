package builder

import (
	"fmt"
	"lpm_builder/pkg/common"
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

type ScriptType string

func PrepareScript(stage0Path string, scriptType ScriptType) string {
	script := filepath.Join(stage0Path, string(scriptType))
	preparedScript := fmt.Sprintf(`
		#!/bin/bash
		set -e

		%s

		target_script=$(<%s)
		eval "$target_script"
	`, BuiltInFunctions, script)

	return preparedScript
}

func executeStage0(ctx *BuilderCtx) {
	scripts := [][2]string{
		{Init, ctx.TmpPkgDir},
		{Build, ctx.TmpSrcDir},
		{InstallFiles, ctx.TmpPkgDir},
		{PostInstallFiles, ctx.TmpPkgDir},
	}

	for _, tuple := range scripts {
		script := tuple[0]
		dir := tuple[1]

		common.Logger.Printf("Executing stage0/%s script", script)

		cmd := exec.Command("/bin/bash", "-c", PrepareScript(ctx.Stage0ScriptsDir, ScriptType(script)))
		cmd.Dir = dir
		_, err := cmd.Output()
		common.FailOnError(err, "Couldn't execute "+script+" script from template directory.")
	}

}
