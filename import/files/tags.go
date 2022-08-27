package files

import (
	"strings"
	"unicode"
)

// Parse entity tags in a `<!DOCTYPE...>` directive and returns a map of
// entities to their respective values.
//
// Entities are given as `<!ENTITY uK "word usually written using kanji alone">`
func ParseDoctypeTags(text string) (resultTags map[string]string) {
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
