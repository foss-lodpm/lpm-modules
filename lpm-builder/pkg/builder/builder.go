package builder

import (
	"fmt"
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

	ctx.PkgFilesData = make([]File, 0)

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

	lpmExecutable, err := os.Executable()
	common.FailOnError(err, "Couldn't set $LPM")
	err = os.Setenv("LPM", lpmExecutable)

	err = os.Setenv(BUILD_ROOT_ENV_TAG, ctx.TmpPkgDir)
	common.FailOnError(err, "Couldn't set $"+BUILD_ROOT_ENV_TAG)

	err = os.Setenv(SRC_ENV_TAG, ctx.TmpSrcDir)
	common.FailOnError(err, "Couldn't set $"+SRC_ENV_TAG)

	return ctx
}

func marshalAndWriteSystemJson(ctx *BuilderCtx) {
	file, err := common.Utf8FriendlyJsonMarshal(ctx.System)
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

func StartBuilding(templateDir string, targetBuild string) {
	ctx := prepare(templateDir)

	build, targetExists := ctx.TemplateFields.Builds[""]
	common.Assert(targetExists, fmt.Sprintf("Target build '%s' not found in `builds`.", targetBuild))

	InstallBuildTimeDependencies(build)

	CopyProvidedStage1Scripts(&ctx)

	executeStage0(&ctx)
	computeChecksumsAndInstallSize(&ctx, build)
	generateMetaFiles(&ctx, build)
	marshalAndWriteSystemJson(&ctx)

	// Get the current working directory
	workingDir, err := os.Getwd()
	common.FailOnError(err, "Failed on getting current path.")

	// Get the absolute path of the current working directory
	workingDir, err = filepath.Abs(workingDir)
	common.FailOnError(err)

	common.Logger.Println("Packing and generating the package")
	pkgOutputName := fmt.Sprintf("%s-%s.lod", ctx.TemplateFields.Name, build.Version.ReadableFormat)
	lodOutputPath := filepath.Join(workingDir, pkgOutputName)
	PackLodFile(ctx.TmpPkgDir, lodOutputPath)

	common.Logger.Println("Writing index patch")
	GenerateIndexPatch(&ctx, build)

	cleanup(ctx)
}
