package jmdict

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/ax-lab/tango/import/files"
)

const (
	xmlRootElement   = "JMdict"
	xmlEntryElement  = "entry"
	xmlDoctypePrefix = "DOCTYPE JMdict ["
)

type Decoder struct {
	Tags map[string]string
	xml  *xml.Decoder
	init bool
}

func NewDecoder(input io.Reader) *Decoder {
	out := &Decoder{
		xml: xml.NewDecoder(input),
	}
	return out
}

func (decoder *Decoder) ReadEntry() (*Entry, error) {
	input := decoder.xml
	if !decoder.init {
		decoder.init = true
		for {
			if token, err := input.Token(); err != nil {
				if err == io.EOF {
					return nil, fmt.Errorf("invalid schema: root `%s` element not found", xmlRootElement)
				}
				return nil, err
			} else if dir, ok := token.(xml.Directive); ok && bytes.HasPrefix(dir, []byte(xmlDoctypePrefix)) {

				// provides JMdict custom entities to the result
				decoder.Tags = files.ParseDoctypeTags(string(dir[len(xmlDoctypePrefix):]))

				// add the custom entities to the XML decoder as well, so they
				// can be parsed
				input.Entity = make(map[string]string)
				for key := range decoder.Tags {
					// we use the short entity codes as their own value since
					// descriptions are available in `decoder.Tags`
					input.Entity[key] = key
				}

			} else if start, ok := token.(xml.StartElement); ok && start.Name.Local == xmlRootElement {
				break
			}
		}
	}

	for {
		if token, err := input.Token(); err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		} else if start, ok := token.(xml.StartElement); ok && start.Name.Local == xmlEntryElement {
			var entry Entry
			if entryErr := input.DecodeElement(&entry, &start); entryErr != nil {
				return nil, fmt.Errorf("decoding entry: %v", entryErr)
			}
			if entry.Sequence == 0 {
				return nil, fmt.Errorf("invalid entry: missing sequence")
			}

			for i, it := range entry.Reading {
				entry.Reading[i].NoKanji = it.NoKanjiRaw != nil
				entry.Reading[i].NoKanjiRaw = nil
			}
			return &entry, nil
		}
	}
}
