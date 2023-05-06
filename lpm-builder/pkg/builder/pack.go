package builder

import (
	"context"
	"os"
	"path/filepath"

	"github.com/mholt/archiver/v4"
)

func PackLodFile(rootDir string, destFile string) error {
	// absolute path -- package level path
	filesToPack := map[string]string{
		filepath.Join(rootDir, "meta"):        "meta",
		filepath.Join(rootDir, "program"):     "program",
		filepath.Join(rootDir, "scripts"):     "scripts",
		filepath.Join(rootDir, "system.json"): "system.json",
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
