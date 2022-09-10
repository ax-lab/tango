package data

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path"
)

type DataWriterJSON interface {
	WriteJSON(output io.Writer) error
}

type DataWriter struct {
	err   error
	dir   string
	files map[string]*dataWriterFile
}

func NewDataWriter(outputDirectory string) *DataWriter {
	err := os.MkdirAll(outputDirectory, os.ModePerm)
	return &DataWriter{
		err:   err,
		dir:   outputDirectory,
		files: make(map[string]*dataWriterFile),
	}
}

func (w *DataWriter) Error() error {
	return w.err
}

func (w *DataWriter) Close() {
	for _, file := range w.files {
		if !file.closed {
			file.WriteString("\n]\n")
		}
		file.Close()
	}
	w.files = nil
}

func (w *DataWriter) WriteToFile(fileName string, item any) {
	w.doWriteToFile(fileName, item, "", "")
}

func (w *DataWriter) WriteToFileIndent(fileName string, item any) {
	w.doWriteToFile(fileName, item, "", "\t")
}

func (w *DataWriter) doWriteToFile(fileName string, item any, prefix, indent string) {
	if file, isNew := w.getFile(fileName); file != nil {
		if !isNew {
			panic("writing to non-empty file: " + fileName)
		}
		file.WriteJSON(item, prefix, indent)
		file.closed = true
	}
}

func (w *DataWriter) AppendToFile(fileName string, item any) {
	w.doAppendToFile(fileName, item, "", "")
}

func (w *DataWriter) AppendToFileIndent(fileName string, item any) {
	w.doAppendToFile(fileName, item, "\t", "\t")
}

func (w *DataWriter) doAppendToFile(fileName string, item any, prefix, indent string) {
	if file, isNew := w.getFile(fileName); file != nil {
		if file.closed {
			panic("appending to finished file: " + fileName)
		}
		if isNew {
			file.WriteString("[\n")
		}
		if file.items {
			file.WriteString(",\n")
		}
		file.items = true
		if prefix != "" {
			file.WriteString(prefix)
		}
		file.WriteJSON(item, prefix, indent)
	}
}

func (w *DataWriter) getFile(fileName string) (out *dataWriterFile, isNew bool) {
	if w.err != nil {
		return nil, false
	}
	out = w.files[fileName]
	if out == nil {
		var fullPath = path.Join(w.dir, fileName)
		isNew = true
		if parentPath := path.Dir(fullPath); parentPath != "" {
			if err := os.MkdirAll(parentPath, os.ModePerm); err != nil {
				w.err = err
				return nil, false
			}
		}
		fp, err := os.Create(fullPath)
		if err != nil {
			w.err = err
			return nil, false
		}
		out = &dataWriterFile{src: w, file: fp, out: bufio.NewWriter(fp)}
		w.files[fileName] = out
	}
	return out, isNew
}

type dataWriterFile struct {
	src    *DataWriter
	file   io.WriteCloser
	out    *bufio.Writer
	items  bool
	closed bool
}

func (w *dataWriterFile) WriteString(data string) {
	if w.src.err != nil {
		return
	}
	_, w.src.err = w.Write(([]byte)(data))
}

func (w *dataWriterFile) WriteJSON(value any, prefix, indent string) {
	if w.src.err != nil {
		return
	}

	var bytes []byte
	var err error
	if indent != "" {
		bytes, err = json.MarshalIndent(value, prefix, indent)
	} else {
		bytes, err = json.Marshal(value)
	}

	if err == nil {
		_, w.src.err = w.Write(bytes)
	} else {
		w.src.err = err
	}
}

func (w *dataWriterFile) Close() error {
	w.out.Flush()
	return w.file.Close()
}

func (w *dataWriterFile) Write(data []byte) (int, error) {
	return w.out.Write(data)
}
