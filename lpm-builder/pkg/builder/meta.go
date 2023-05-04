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

func computeChecksums(ctx *BuilderCtx) {

	err := filepath.Walk(ctx.TmpProgramDir, func(path string, info os.FileInfo, err error) error {
		common.FailOnError(err, "Failed while searching files in "+ctx.TmpProgramDir)

		if !info.IsDir() {
			file := File{
				Path:              strings.Split(path, "/program/")[1],
				ChecksumAlgorithm: ctx.TemplateFields.FileChecksumAlgo,
				Checksum:          getHashOfFile(path, ctx.TemplateFields.FileChecksumAlgo),
			}

			ctx.PkgFilesData = append(ctx.PkgFilesData, file)
		}

		return nil
	})

	common.FailOnError(err, "filepath.Walk failed for "+ctx.TmpProgramDir)
}

func marshalAndWriteFilesJson(ctx *BuilderCtx) {
	file, err := json.MarshalIndent(ctx.PkgFilesData, "", " ")
	common.FailOnError(err, "Failed on serializing ctx.PkgFilesData")

	filesJsonPath := filepath.Join(ctx.TmpMetaDir, "files.json")
	err = ioutil.WriteFile(filesJsonPath, file, 0644)
	common.FailOnError(err)
}

func generateMetaFiles(ctx *BuilderCtx) {
	marshalAndWriteFilesJson(ctx)

	// TODO
	// marshalAndWriteMetaJson(ctx)
}
