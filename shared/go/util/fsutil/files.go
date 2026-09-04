package fsutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopyFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() {
		_ = srcFile.Close()
	}()

	dstDir := filepath.Dir(dstPath)

	// create all dirs to file
	err = os.MkdirAll(dstDir, 0755)
	if err != nil {
		return fmt.Errorf("failed create fir for copy file %s: %w", dstDir, err)
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		_ = dstFile.Close()
	}()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	return nil
}

func CreateDirWithEmptyFiles(dirPath string, files []string) error {
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	for _, fileName := range files {
		f, err := os.Create(dirPath + "/" + fileName)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		_ = f.Close()
	}
	return nil
}

// returns all files in the directory with relative paths
func AllFilesInDir(dirPath string, recursive bool) (fileNames []string, err error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %w", err)
	}

	if len(entries) == 0 {
		return []string{}, nil
	}

	dirs := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, dirPath+"/"+entry.Name())
			continue
		}
		fileNames = append(fileNames, dirPath+"/"+entry.Name())
	}

	if !recursive {
		return fileNames, nil
	}

	for _, dir := range dirs {
		subFiles, err := AllFilesInDir(dir, true)
		if err != nil {
			return nil, fmt.Errorf("failed to read subdirectory \"%s\": %w", dir, err)
		}
		fileNames = append(fileNames, subFiles...)
	}
	return fileNames, nil
}
