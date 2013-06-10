package site

import (
	"io"
	"os"
	"path/filepath"
)

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

func MakeDirectoriesTo(path string, mode os.FileMode) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(path), mode)
	}
	return err
}

func mapToFuncs(vars map[string]string) map[string]interface{} {
	funcs := make(map[string]interface{})
	for name, val := range vars {
		funcs[name] = func() string { return val }
	}
	return funcs
}
