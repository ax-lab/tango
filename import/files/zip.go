package files

import (
	"archive/zip"
	"io/fs"
	"os"
)

type ZipCloser struct {
	*zip.Reader
	file *os.File
}

func FindZip(fileName string) (*ZipCloser, error) {
	var size int64
	file, err := Find(fileName)
	if err == nil {
		stat, statErr := file.Stat()
		err = statErr
		size = stat.Size()
		if err == nil {
			zipFile, zipErr := zip.NewReader(file, size)
			err = zipErr
			if err == nil {
				return &ZipCloser{
					Reader: zipFile,
					file:   file,
				}, nil
			}
		}
	}
	if err != nil && file != nil {
		file.Close()
	}
	return nil, err
}

func (input *ZipCloser) OpenFileByName(name string) (fs.File, error) {
	for _, it := range input.File {
		if it.FileInfo().Name() == name {
			return input.Open(it.Name)
		}
	}
	return nil, ErrNotFound{name}
}

func (input *ZipCloser) Close() {
	input.file.Close()
}
