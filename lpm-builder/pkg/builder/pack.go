package builder

import (
	"archive/tar"
	"github.com/pierrec/lz4"
	"io"
	"os"
	"path/filepath"
)

func PackLodFile(rootDir string, srcPaths []string, destFile string) error {
	// Change to the source directory
	// This is for preventing compressing the whole file/dir paths under the lod file
	if err := os.Chdir(rootDir); err != nil {
		return err
	}

	// Create the output file
	outFile, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Create a new LZ4 writer on top of the output file
	lz4Writer := lz4.NewWriter(outFile)
	defer lz4Writer.Close()

	// Create a new tar writer on top of the LZ4 writer
	tarWriter := tar.NewWriter(lz4Writer)
	defer tarWriter.Close()

	// Walk through the source paths and add files to the tar archive
	for _, srcPath := range srcPaths {
		if err := filepath.Walk(srcPath, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Create a new tar header for the file or directory
			header := &tar.Header{
				Name:    filePath,
				Mode:    int64(info.Mode()),
				ModTime: info.ModTime(),
				Size:    info.Size(),
			}

			if info.IsDir() {
				header.Typeflag = tar.TypeDir
			}

			// Write the header to the tar archive
			if err := tarWriter.WriteHeader(header); err != nil {
				return err
			}

			if !info.IsDir() {
				// Open the source file
				inFile, err := os.Open(filePath)
				if err != nil {
					return err
				}
				defer inFile.Close()

				// Copy the file contents to the tar archive
				if _, err := io.Copy(tarWriter, inFile); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}
