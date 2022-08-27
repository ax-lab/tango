package kanji

import (
	"encoding/xml"
	"fmt"
	"io"
)

type Decoder struct {
	Info Info
	xml  *xml.Decoder
	init bool
}

const (
	xmlRootElement      = "kanjidic2"
	xmlHeaderElement    = "header"
	xmlCharacterElement = "character"
)

func NewDecoder(input io.Reader) *Decoder {
	out := &Decoder{
		xml: xml.NewDecoder(input),
	}
	return out
}

func (decoder *Decoder) ReadCharacter() (*Character, error) {
	input := decoder.xml
	if !decoder.init {
		decoder.init = true
		for {
			if token, err := input.Token(); err != nil {
				if err == io.EOF {
					return nil, fmt.Errorf("invalid schema: root `%s` element not found", xmlRootElement)
				}
				return nil, err
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
		} else if start, ok := token.(xml.StartElement); ok {
			var character Character
			switch start.Name.Local {
			case xmlHeaderElement:
				if headerErr := input.DecodeElement(&decoder.Info, &start); headerErr != nil {
					return nil, fmt.Errorf("decoding header: %v", headerErr)
				}
			case xmlCharacterElement:
				if characterErr := input.DecodeElement(&character, &start); characterErr != nil {
					return nil, fmt.Errorf("decoding character: %v", characterErr)
				}
				if character.Literal == "" {
					return nil, fmt.Errorf("invalid character: missing literal")
				}
				return &character, nil
			}
		}
	}
}
