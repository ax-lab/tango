package jmnedict

import (
	"compress/gzip"
	"io"

	"github.com/ax-lab/tango/import/files"
)

const DefaultFilePath = "vendor/data/entries/JMnedict.xml.gz"

func Load(filePath string) (io.Reader, error) {
	input, err := files.Find(filePath)
	if err != nil {
		return nil, err
	}
	return gzip.NewReader(input)
}
