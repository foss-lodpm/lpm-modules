package builder

import (
	"encoding/json"
	"io/ioutil"
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
	TmpScriptsDir    string
	Stage0ScriptsDir string
	Stage1ScriptsDir string
	TemplateFields   template.Template
	PkgFilesData     []File
	System           common.System
	InstallSize      uint
}

func prepare(templateDir string) BuilderCtx {
	common.Logger.Println("Preparing for building")

	var ctx BuilderCtx
	ctx.TemplateDir = templateDir

	ctx.TemplateFields = template.DeserializeTemplate(ctx.TemplateDir)

	ctx.System.BuilderVersion = common.BuilderVersion
	ctx.System.MinSupportedLpmVersion = common.MinSupportedLpmVersion

	ctx.TmpPkgDir = filepath.Join(TMP_ROOT_PATH, ctx.TemplateFields.Name)
	ctx.TmpSrcDir = filepath.Join(ctx.TmpPkgDir, "src")
	ctx.TmpMetaDir = filepath.Join(ctx.TmpPkgDir, "meta")
	ctx.TmpProgramDir = filepath.Join(ctx.TmpPkgDir, "program")
	ctx.TmpScriptsDir = filepath.Join(ctx.TmpPkgDir, "scripts")
	ctx.Stage0ScriptsDir = filepath.Join(ctx.TemplateDir, "stage0")
	ctx.Stage1ScriptsDir = filepath.Join(ctx.TemplateDir, "stage1")

	for _, dir := range []string{ctx.TmpMetaDir, ctx.TmpProgramDir, ctx.TmpSrcDir, ctx.TmpScriptsDir} {
		err := os.MkdirAll(dir, os.ModePerm)
		common.FailOnError(err, "Couldn't create "+dir+" directory.")

	}

	err := os.Setenv(BUILD_ROOT_ENV_TAG, ctx.TmpPkgDir)
	common.FailOnError(err, "Couldn't set $"+BUILD_ROOT_ENV_TAG)

	err = os.Setenv(SRC_ENV_TAG, ctx.TmpSrcDir)
	common.FailOnError(err, "Couldn't set $"+SRC_ENV_TAG)

	return ctx
}

func marshalAndWriteSystemJson(ctx *BuilderCtx) {
	file, err := json.MarshalIndent(ctx.System, "", "\t")
	common.FailOnError(err, "Failed on serializing ctx.System")

	common.Logger.Println("Writing system.json")
	filesJsonPath := filepath.Join(ctx.TmpPkgDir, "system.json")
	err = ioutil.WriteFile(filesJsonPath, file, 0644)
	common.FailOnError(err)
}

func cleanup(ctx BuilderCtx) {
	err := os.RemoveAll(ctx.TmpPkgDir)
	common.FailOnError(err, "Failed cleaning the temporary files of building lod package.")
}

func StartBuilding(templateDir string) {
	ctx := prepare(templateDir)

	CopyProvidedStage1Scripts(&ctx)

	executeStage0(&ctx)
	computeChecksumsAndInstallSize(&ctx)
	generateMetaFiles(&ctx)
	marshalAndWriteSystemJson(&ctx)
	// TODO
	// Pack .lod file to current dir
	// cleanup(ctx)
}
