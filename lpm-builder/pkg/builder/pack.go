package builder

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mholt/archiver/v4"
)

func PackLodFile(rootDir string, destFile string) error {
	filesToPack := map[string]string{
		filepath.Join(rootDir, "meta"):        "meta",
		filepath.Join(rootDir, "scripts"):     "scripts",
		filepath.Join(rootDir, "system.json"): "system.json",
	}

	programPath := filepath.Join(rootDir, "program")
	if content, _ := ioutil.ReadDir(programPath); len(content) != 0 {
		filesToPack[programPath] = "program"
    } else {
		filesToPack[filepath.Join(rootDir, "src") ] = "src"
	}

	files, err := archiver.FilesFromDisk(nil, filesToPack)

	if err != nil {
		return err
	}

	// create the output file we'll write to
	out, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer out.Close()

	format := archiver.CompressedArchive{
		Compression: archiver.Lz4{},
		Archival:    archiver.Tar{},
	}

	err = format.Archive(context.Background(), out, files)
	if err != nil {
		return err
	}

	return nil
}
