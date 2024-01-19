package builder

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"lpm_builder/pkg/common"
	"lpm_builder/pkg/template"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Path              string `json:"path"`
	ChecksumAlgorithm string `json:"checksum_algorithm"`
	Checksum          string `json:"checksum"`
}

type Meta struct {
	Name          string              `json:"name"`
	InstalledSize uint                `json:"installed_size"`
	Version       common.Version      `json:"version"`
	Dependencies  []common.Dependency `json:"dependencies"`
	Suggestions   []common.Dependency `json:"suggestions"`
}

func getHashOfFile(filePath string, hashAlgorithm string) string {
	file, err := os.Open(filePath)
	common.FailOnError(err)
	defer file.Close()

	var hash hash.Hash

	switch hashAlgorithm {
	case "md5":
		hash = md5.New()
	case "sha256":
		hash = sha256.New()
	case "sha512":
		hash = sha512.New()
	default:
		common.FatalError("Unexpected checksum algorithm." + hashAlgorithm + "  is not supported")
	}

	_, err = io.Copy(hash, file)
	common.FailOnError(err)

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func computeChecksumsAndInstallSize(ctx *BuilderCtx, build *template.Build) {
	err := filepath.Walk(ctx.TmpProgramDir, func(path string, info os.FileInfo, err error) error {
		common.FailOnError(err, "Failed while searching files in "+ctx.TmpProgramDir)

		if !info.IsDir() {
			ctx.InstallSize += uint(info.Size() / 1024)

			file := File{
				Path:              strings.Split(path, "/program/")[1],
				ChecksumAlgorithm: *build.FileChecksumAlgo,
				Checksum:          getHashOfFile(path, *build.FileChecksumAlgo),
			}

			common.Logger.Printf("computed %s as %s checksum for file %s", file.Checksum, file.ChecksumAlgorithm, file.Path)

			ctx.PkgFilesData = append(ctx.PkgFilesData, file)
		}

		return nil
	})

	common.FailOnError(err, "filepath.Walk failed for "+ctx.TmpProgramDir)
}

func genMetaFromTemplateFields(ctx *BuilderCtx, build *template.Build) Meta {
	var meta Meta

	meta.Name = ctx.TemplateFields.Name
	meta.InstalledSize = ctx.InstallSize
	meta.Version = build.Version
	meta.Dependencies = build.MandatoryDependencies.Runtime
	meta.Suggestions = build.SuggestedDependencies.Runtime

	return meta
}

func marshalAndWriteFilesJson(ctx *BuilderCtx) {
	file, err := common.Utf8FriendlyJsonMarshal(ctx.PkgFilesData)
	common.FailOnError(err, "Failed on serializing ctx.PkgFilesData")

	common.Logger.Println("Writing meta/files.json")
	filesJsonPath := filepath.Join(ctx.TmpMetaDir, "files.json")
	err = ioutil.WriteFile(filesJsonPath, file, 0644)
	common.FailOnError(err)
}

func marshalAndWriteMetaJson(meta Meta, metaDir string) {
	file, err := common.Utf8FriendlyJsonMarshal(meta)
	common.FailOnError(err, "Failed on serializing Meta fields")

	common.Logger.Println("Writing meta/meta.json")
	filesJsonPath := filepath.Join(metaDir, "meta.json")
	err = ioutil.WriteFile(filesJsonPath, file, 0644)
	common.FailOnError(err)
}

func generateMetaFiles(ctx *BuilderCtx, build *template.Build) {
	marshalAndWriteFilesJson(ctx)

	meta := genMetaFromTemplateFields(ctx, build)
	marshalAndWriteMetaJson(meta, ctx.TmpMetaDir)
}
