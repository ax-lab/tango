package frequency

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Entry struct {
	Text    string
	Blog    EntryInfo
	Twitter EntryInfo
	News    EntryInfo
}

type EntryInfo struct {
	// Raw frequency count.
	Freq int64

	// Frequency per-million as a decimal floating point.
	FreqPM string

	// Contextual diversity.
	CD int64

	// Percentage of contextual diversity as a decimal floating point.
	CDPc string
}

var reEntryDecimal = regexp.MustCompile(`^\d+(\.\d+)?$`)

func ParseEntry(input string) (*Entry, error) {
	line := strings.TrimRightFunc(input, unicode.IsSpace)
	if line == "" {
		return nil, nil
	}

	wrapErr := func(msg string) error {
		return fmt.Errorf("parsing frequency entry: %s", msg)
	}

	if fields := strings.Split(line, "\t"); len(fields) == 13 {
		var err error

		parseInt := func(field string) (out int64) {
			if err == nil {
				out, err = strconv.ParseInt(field, 10, 64)
			}
			return out
		}

		parseDec := func(field string) string {
			if err == nil && !reEntryDecimal.MatchString(field) {
				err = fmt.Errorf("invalid decimal: %s", field)
			}
			return field
		}

		parseInfo := func(index int) EntryInfo {
			return EntryInfo{
				Freq:   parseInt(fields[index+0]),
				FreqPM: parseDec(fields[index+1]),
				CD:     parseInt(fields[index+2]),
				CDPc:   parseDec(fields[index+3]),
			}
		}

		if strings.TrimSpace(fields[0]) == "" {
			return nil, nil
		}

		entry := &Entry{
			Text:    fields[0],
			Blog:    parseInfo(1),
			Twitter: parseInfo(5),
			News:    parseInfo(9),
		}

		if err != nil {
			return nil, wrapErr(err.Error())
		}

		return entry, nil
	}

	return nil, wrapErr("invalid line")
}
