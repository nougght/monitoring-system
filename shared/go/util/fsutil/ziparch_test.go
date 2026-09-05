package fsutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZipFolder(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	dirName := "test_dir"
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	nestedDir := dirName + "/nested"
	nestedFiles := []string{"nested_file1.txt", "nested_file2.txt"}

	err := CreateDirWithEmptyFiles(dirName, files)
	require.NoError(t, err, "failed to create directory with empty files")
	err = CreateDirWithEmptyFiles(nestedDir, nestedFiles)
	require.NoError(t, err, "failed to create nested directory with empty files")

	defer func() {
		err := os.RemoveAll(dirName)
		require.NoError(t, err, "failed to remove directory")
	}()

	zipFilePath := tempDir + "/test.zip"
	err = ZipFolder(zipFilePath, dirName)
	require.NoError(t, err, "failed to zip folder")

	_, err = os.Stat(zipFilePath)
	require.NoError(t, err, "zip file does not exist")

}
