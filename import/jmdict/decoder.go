package jmdict

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"unicode"
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
				decoder.Tags = parseTags(string(dir[len(xmlDoctypePrefix):]))

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
			return &entry, nil
		}
	}
}

// Parse entity tags in a `<!DOCTYPE...>` directive and returns a map of
// entities to their respective values.
//
// Entities are given as `<!ENTITY uK "word usually written using kanji alone">`
func parseTags(text string) (resultTags map[string]string) {
	const entityTag = "<!ENTITY "

	resultTags = make(map[string]string)
	next := text
	for len(next) > 0 {
		pos := strings.Index(next, entityTag)
		if pos < 0 {
			break
		}

		next = next[pos+len(entityTag):]

		isSpace := func(chr byte) bool {
			return unicode.IsSpace(rune(chr))
		}

		state, valid, end := ' ', false, false
		for pos = 0; !end && pos < len(next); pos++ {
			switch chr := next[pos]; state {
			case ' ':
				if !isSpace(chr) {
					state = 'n'
				}
			case 'n':
				if isSpace(chr) {
					state = 'q'
				}
			case 'q':
				if !isSpace(chr) {
					if chr == '\'' || chr == '"' {
						state = rune(chr)
					} else {
						state = 'x'
					}
				}
			case '\'', '"':
				if rune(chr) == state {
					end, valid = true, true
				}
			default:
				end = chr == '>'
			}
		}

		if valid {
			entity := strings.SplitN(strings.TrimSpace(next[:pos-1]), " ", 2)
			if len(entity) == 2 {
				name := strings.TrimSpace(entity[0])
				value := strings.TrimSpace(entity[1][1:]) // skips opening quote
				if name != "" && value != "" {
					resultTags[name] = value
				}
			}
		}

		next = next[pos:]
	}

	return resultTags
}
