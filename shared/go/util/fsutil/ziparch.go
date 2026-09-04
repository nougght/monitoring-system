package fsutil

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func handleZipFileName(zipName string) string {
	if !strings.HasSuffix(zipName, ".zip") {
		zipName = zipName + ".zip"
	}
	return zipName
}

// archive folder and save it to zip file
func ZipFolder(zipFilePath, folderPath string) error {

	arch, err := os.Create(handleZipFileName(zipFilePath))
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}

	defer func() {
		_ = arch.Close()
	}()

	zipBytes, err := ZipFolderRaw(folderPath)
	if err != nil {
		return fmt.Errorf("failed to zip folder: %w", err)
	}
	_, err = arch.Write(zipBytes)
	if err != nil {
		return fmt.Errorf("failed to write zip bytes to file: %w", err)
	}
	return nil
}

// archive folder and return zip bytes
func ZipFolderRaw(folderPath string) ([]byte, error) {

	allFiles, err := AllFilesInDir(folderPath, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get all files in directory: %w", err)
	}
	filesData := make(map[string][]byte, len(allFiles))

	for _, filePath := range allFiles {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file to zip with path \"%s\", error: %w", filePath, err)
		}
		filesData[filePath] = data
	}

	buf := new(bytes.Buffer)

	zipWriter := zip.NewWriter(buf)

	for filename, data := range filesData {
		name, err := filepath.Rel(folderPath, filename)
		writer, err := zipWriter.Create(name)
		if err != nil {
			return nil, fmt.Errorf("failed to create file inside zip: %w", err)
		}

		_, err = writer.Write(data)
		if err != nil {
			return nil, fmt.Errorf("failed to write content to file: %w", err)
		}
	}
	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zip writer: %w", err)
	}

	return buf.Bytes(), nil
}
