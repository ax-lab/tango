package files

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type ErrNotFound struct {
	File string
}

func (e ErrNotFound) Unwrap() error {
	return fs.ErrNotExist
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("import file `%s` not found", e.File)
}

func Find(filePath string) (*os.File, error) {
	currentPath, err := os.Getwd()
	for err == nil && currentPath != "" {
		fullPath := filepath.Join(currentPath, filePath)
		file, fileErr := os.Open(fullPath)
		if fileErr == nil {
			return file, nil
		}

		if errors.Is(fileErr, fs.ErrNotExist) {
			parentDir := filepath.Dir(currentPath)
			if parentDir != currentPath {
				currentPath = parentDir
			} else {
				err = ErrNotFound{filePath}
			}
		} else {
			return nil, fileErr
		}
	}

	return nil, err
}
