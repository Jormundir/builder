package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func MakeDirectories(path string, mode os.FileMode) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, mode)
	}
	return err
}

func CopyFile(dest, src string) (int64, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return 0, err
	}
	defer destFile.Close()

	return io.Copy(destFile, srcFile)
}

func ReferencePath(path, prefix string) string {
	return filepath.ToSlash(TrimPath(path, prefix))
}

func TrimPath(path, prefix string) string {
	newPath := strings.TrimPrefix(path, prefix)
	return strings.TrimSuffix(newPath, filepath.Ext(newPath))
}
