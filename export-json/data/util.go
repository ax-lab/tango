package data

import "strings"

func TSV(value interface{}) []string {
	input := value.(string)
	if input == "" {
		return nil
	}
	return strings.Split(input, "\t")
}
