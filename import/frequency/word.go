package frequency

import (
	"fmt"
	"strconv"
	"strings"
)

type Word struct {
	Entry string
	Count int64
}

func ParseWord(line string) (*Word, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}

	wrapErr := func(msg string) error {
		return fmt.Errorf("parsing word frequency: %s", msg)
	}

	if fields := strings.Split(line, "\t"); len(fields) == 2 {
		if count, err := strconv.ParseInt(fields[0], 10, 64); err != nil {
			return nil, wrapErr(err.Error())
		} else {
			return &Word{Entry: fields[1], Count: count}, nil
		}
	}

	return nil, wrapErr("invalid line")
}
