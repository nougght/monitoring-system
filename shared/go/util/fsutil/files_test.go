package fsutil

import (
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyFile(t *testing.T) {
	fileName := "test.txt"

	srcDir := t.TempDir()
	dstDir := t.TempDir()
	defer func() {
		err := os.RemoveAll(srcDir)
		require.NoError(t, err, "failed to remove all src dir temp dir")
		err = os.RemoveAll(dstDir)
		require.NoError(t, err, "failed to remove all dst dir temp dir")
	}()
	srcFilePath := fmt.Sprintf("%s/%s", srcDir, fileName)

	f, err := os.Create(srcFilePath)
	require.NoError(t, err, "failed to create test file")
	err = f.Close()
	require.NoError(t, err, "failed to close test file")

	dstFilePath := fmt.Sprintf("%s/%s", dstDir, fileName)

	err = CopyFile(srcFilePath, dstFilePath)
	require.NoError(t, err, "failed to copy file")

	_, err = os.Stat(dstFilePath)
	require.NoError(t, err, "not found copied file")
}

func TestCreateDirWithEmptyFiles(t *testing.T) {
	dirName := "test_dir"
	files := []string{"file1.txt", "file2.txt", "file3.txt"}

	err := CreateDirWithEmptyFiles(dirName, files)
	require.NoError(t, err, "failed to create directory with empty files")
	defer func() {
		err := os.RemoveAll(dirName)
		require.NoError(t, err, "failed to remove directory")
	}()

	for _, fileName := range files {
		filePath := dirName + "/" + fileName
		_, err := os.Stat(filePath)
		require.NoError(t, err, "file does not exist: %s", filePath)
	}
}
func TestAllFilesInDir(t *testing.T) {
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

	expectedNonRecursive := make([]string, 0, len(files)+len(nestedFiles))
	for _, file := range files {
		expectedNonRecursive = append(expectedNonRecursive, dirName+"/"+file)
	}
	expectedRecursive := slices.Clone(expectedNonRecursive)
	for _, file := range nestedFiles {
		expectedRecursive = append(expectedRecursive, nestedDir+"/"+file)
	}

	allFiles, err := AllFilesInDir(dirName, false)
	require.NoError(t, err, "failed to get all files in directory (non-recursive)")

	require.ElementsMatch(t, expectedNonRecursive, allFiles, "unexpected files in directory (non-recursive)")

	allFiles, err = AllFilesInDir(dirName, true)
	require.NoError(t, err, "failed to get all files in directory (recursive)")
	require.ElementsMatch(t, expectedRecursive, allFiles, "unexpected files in directory (recursive)")

}
