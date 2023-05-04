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
	TmpMetaDir       string
	TmpProgramDir    string
	Stage0ScriptsDir string
	Stage1ScriptsDir string
	TemplateFields   template.Template
	PkgFilesData     []File
}

func prepare(templateDir string) BuilderCtx {
	var ctx BuilderCtx
	ctx.TemplateDir = templateDir

	// TODO
	// validate deserialized template fields
	ctx.TemplateFields = template.DeserializeTemplate(ctx.TemplateDir)

	ctx.TmpPkgDir = filepath.Join(TMP_ROOT_PATH, ctx.TemplateFields.Name)
	ctx.TmpSrcDir = filepath.Join(ctx.TmpPkgDir, "src")
	ctx.TmpMetaDir = filepath.Join(ctx.TmpPkgDir, "meta")
	ctx.TmpProgramDir = filepath.Join(ctx.TmpPkgDir, "program")
	ctx.Stage0ScriptsDir = filepath.Join(ctx.TemplateDir, "stage0")
	ctx.Stage1ScriptsDir = filepath.Join(ctx.TemplateDir, "stage1")

	for _, dir := range []string{ctx.TmpMetaDir, ctx.TmpProgramDir, ctx.TmpSrcDir} {
		err := os.MkdirAll(dir, os.ModePerm)
		common.FailOnError(err, "Couldn't create "+dir+" directory.")

	}

	err := os.Setenv(BUILD_ROOT_ENV_TAG, ctx.TmpPkgDir)
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

	executeStage0(&ctx)
	computeChecksums(&ctx)
	generateMetaFiles(&ctx)

	// TODO
	// cleanup(ctx)
}
