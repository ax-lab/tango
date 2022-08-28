package frequency

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/ax-lab/tango/import/files"
)

const (
	DefaultPairsFile = "vendor/data/frequency/Innocent_Novel_Analysis_120526.zip"
	DefaultInfoFile  = "vendor/data/frequency/Jap.Freq.2.zip"

	InfoFileChars = "Jap.Char.Freq.2.txt"
	InfoFileWords = "Jap.Freq.2.txt"

	PairsFileJparser = "word_freq_report_jparser.txt"
	PairsFileMecab   = "word_freq_report_mecab.txt"
	PairsFileKanji   = "kanji_freq_report.txt"
)

func LoadInfo(fileName string) (words, chars []*Info, err error) {
	zipFile, err := files.FindZip(fileName)
	if err != nil {
		return nil, nil, err
	}
	defer zipFile.Close()

	readInfo := func(path string) (list []*Info, err error) {
		input, err := zipFile.Open(path)
		if err == nil {
			defer input.Close()
			scanner := bufio.NewScanner(input)
			for err == nil && scanner.Scan() {
				var entry *Info
				if text := scanner.Text(); !strings.Contains(text, "BlogFreqPm") {
					entry, err = ParseInfo(text)
					if entry != nil {
						list = append(list, entry)
					}
				}
			}
			if err == nil {
				err = scanner.Err()
			}
		}
		return list, err
	}

	words, err = readInfo(InfoFileWords)
	if err != nil {
		return nil, nil, fmt.Errorf("loading word info: %v", err)
	}

	chars, err = readInfo(InfoFileChars)
	if err != nil {
		return nil, nil, fmt.Errorf("loading char info: %v", err)
	}

	return words, chars, nil
}

func LoadPairs(fileName string) (jparser, mecab, kanji []*Pair, err error) {
	zipFile, err := files.FindZip(fileName)
	if err != nil {
		return nil, nil, nil, err
	}
	defer zipFile.Close()

	readPairs := func(path string) (list []*Pair, err error) {
		input, err := zipFile.OpenFileByName(path)
		if err == nil {
			defer input.Close()
			scanner := bufio.NewScanner(input)
			for err == nil && scanner.Scan() {
				// files have the byte order mark, so we need to strip it
				text := strings.TrimLeft(scanner.Text(), "\uFEFF")
				var word *Pair
				word, err = ParsePair(text)
				if word != nil {
					list = append(list, word)
				}
			}
			if err == nil {
				err = scanner.Err()
			}
		}
		return list, err
	}

	kanji, err = readPairs(PairsFileKanji)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("loading kanji pairs: %v", err)
	}

	jparser, err = readPairs(PairsFileJparser)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("loading jparser pairs: %v", err)
	}

	mecab, err = readPairs(PairsFileMecab)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("loading mecab pairs: %v", err)
	}

	return jparser, mecab, kanji, nil
}
