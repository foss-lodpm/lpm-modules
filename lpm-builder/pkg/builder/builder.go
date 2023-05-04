package builder

import (
	"lpm_builder/pkg/common"
	"lpm_builder/pkg/template"
	"os"
	"path/filepath"
)

const (
	BUILD_ROOT_ENV_TAG = "BUILD_ROOT"
	SRC_ENV_TAG        = "SRC"
	TMP_ROOT_PATH      = "/tmp"
)

type BuilderCtx struct {
	TemplateDir      string
	TmpPkgDir        string
	TmpSrcDir        string
	Stage0ScriptsDir string
	Stage1ScriptsDir string
	TemplateFields   template.Template
}

func prepare(templateDir string) BuilderCtx {
	var ctx BuilderCtx
	ctx.TemplateDir = templateDir

	// TODO
	// validate deserialized template fields
	ctx.TemplateFields = template.DeserializeTemplate(ctx.TemplateDir)

	ctx.TmpPkgDir = filepath.Join(TMP_ROOT_PATH, ctx.TemplateFields.Name)
	ctx.TmpSrcDir = filepath.Join(ctx.TmpPkgDir, "src")
	ctx.Stage0ScriptsDir = filepath.Join(ctx.TemplateDir, "stage0")
	ctx.Stage1ScriptsDir = filepath.Join(ctx.TemplateDir, "stage1")

	err := os.MkdirAll(ctx.TmpSrcDir, os.ModePerm)
	common.FailOnError(err, "Couldn't create source directory for downloading/building package source.")

	err = os.Setenv(BUILD_ROOT_ENV_TAG, ctx.TmpPkgDir)
	common.FailOnError(err, "Couldn't set $"+BUILD_ROOT_ENV_TAG)

	err = os.Setenv(SRC_ENV_TAG, ctx.TmpSrcDir)
	common.FailOnError(err, "Couldn't set $"+SRC_ENV_TAG)

	return ctx
}

func cleanup(ctx BuilderCtx) {
	err := os.RemoveAll(ctx.TmpPkgDir)
	common.FailOnError(err, "Failed cleaning the temporary files of building lod package.")
}

func StartBuilding(templateDir string) {
	ctx := prepare(templateDir)

	executeStage0(ctx)

	// TODO
	// cleanup(ctx)
}
