package builder

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"lpm_builder/pkg/common"
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
	Description   string              `json:"description"`
	Maintainer    string              `json:"maintainer"`
	Repository    string              `json:"repository"`
	PkgRepository string              `json:"pkg_repository"`
	Homepage      string              `json:"homepage"`
	Arch          string              `json:"arch"`
	Kind          string              `json:"kind"`
	InstalledSize uint                `json:"installed_size"`
	Tags          []string            `json:"tags"`
	Version       common.Version      `json:"version"`
	License       string              `json:"license"`
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

func computeChecksumsAndInstallSize(ctx *BuilderCtx) {

	common.Logger.Println(ctx.TmpProgramDir)
	err := filepath.Walk(ctx.TmpProgramDir, func(path string, info os.FileInfo, err error) error {
		common.FailOnError(err, "Failed while searching files in "+ctx.TmpProgramDir)

		if !info.IsDir() {
			ctx.InstallSize += uint(info.Size() / 1024)

			file := File{
				Path:              strings.Split(path, "/program/")[1],
				ChecksumAlgorithm: ctx.TemplateFields.FileChecksumAlgo,
				Checksum:          getHashOfFile(path, ctx.TemplateFields.FileChecksumAlgo),
			}

			common.Logger.Printf("computed %s as %s checksum for file %s", file.Checksum, file.ChecksumAlgorithm, file.Path)

			ctx.PkgFilesData = append(ctx.PkgFilesData, file)
		}

		return nil
	})

	common.FailOnError(err, "filepath.Walk failed for "+ctx.TmpProgramDir)
}

func genMetaFromTemplateFields(ctx *BuilderCtx) Meta {
	var meta Meta

	meta.Name = ctx.TemplateFields.Name
	meta.Description = ctx.TemplateFields.Description
	meta.Maintainer = ctx.TemplateFields.Maintainer
	meta.Repository = ctx.TemplateFields.Repository
	meta.PkgRepository = ctx.TemplateFields.PkgRepository
	meta.Homepage = ctx.TemplateFields.Homepage
	meta.Arch = ctx.TemplateFields.Arch
	meta.Kind = ctx.TemplateFields.Kind
	meta.InstalledSize = ctx.InstallSize
	meta.Tags = ctx.TemplateFields.Tags
	meta.Version = ctx.TemplateFields.Version
	meta.License = ctx.TemplateFields.License
	meta.Dependencies = ctx.TemplateFields.MandatoryDependencies.Runtime
	meta.Suggestions = ctx.TemplateFields.SuggestedDependencies.Runtime

	return meta
}

func marshalAndWriteFilesJson(ctx *BuilderCtx) {
	file, err := json.MarshalIndent(ctx.PkgFilesData, "", "\t")
	common.FailOnError(err, "Failed on serializing ctx.PkgFilesData")

	common.Logger.Println("Writing meta/files.json")
	filesJsonPath := filepath.Join(ctx.TmpMetaDir, "files.json")
	err = ioutil.WriteFile(filesJsonPath, file, 0644)
	common.FailOnError(err)
}

func marshalAndWriteMetaJson(meta Meta, metaDir string) {
	file, err := json.MarshalIndent(meta, "", "\t")
	common.FailOnError(err, "Failed on serializing Meta fields")

	common.Logger.Println("Writing meta/meta.json")
	filesJsonPath := filepath.Join(metaDir, "meta.json")
	err = ioutil.WriteFile(filesJsonPath, file, 0644)
	common.FailOnError(err)
}

func generateMetaFiles(ctx *BuilderCtx) {
	marshalAndWriteFilesJson(ctx)

	meta := genMetaFromTemplateFields(ctx)
	marshalAndWriteMetaJson(meta, ctx.TmpMetaDir)
}
