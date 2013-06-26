package site

import (
	"os"
	fp "path/filepath"
	"strings"
)

func makeDirsTo(p string, mode os.FileMode) error {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		err = os.MkdirAll(fp.Dir(p), mode)
	}
	return err
}

func superMatch(name string, patterns ...string) (matched bool, err error) {
	pathParts := strings.Split(name, string(fp.Separator))
	for _, part := range pathParts {
		for _, pattern := range patterns {
			matched, err = fp.Match(pattern, part)
			if matched {
				return
			}
		}
	}
	return
}
