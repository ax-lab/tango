package frequency

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/ax-lab/tango/import/files"
)

const (
	DefaultWordsFile   = "vendor/data/frequency/Innocent_Novel_Analysis_120526.zip"
	DefaultEntriesFile = "vendor/data/frequency/Jap.Freq.2.zip"

	EntriesFileChars = "Jap.Char.Freq.2.txt"
	EntriesFileWords = "Jap.Freq.2.txt"

	WordsFileJparser = "word_freq_report_jparser.txt"
	WordsFileMecab   = "word_freq_report_mecab.txt"
	WordsFileKanji   = "kanji_freq_report.txt"
)

func LoadEntries(fileName string) (words, chars []*Entry, err error) {
	zipFile, err := files.FindZip(fileName)
	if err != nil {
		return nil, nil, err
	}
	defer zipFile.Close()

	readEntries := func(path string) (list []*Entry, err error) {
		input, err := zipFile.Open(path)
		if err == nil {
			defer input.Close()
			scanner := bufio.NewScanner(input)
			for err == nil && scanner.Scan() {
				var entry *Entry
				if text := scanner.Text(); !strings.Contains(text, "BlogFreqPm") {
					entry, err = ParseEntry(text)
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

	words, err = readEntries(EntriesFileWords)
	if err != nil {
		return nil, nil, fmt.Errorf("loading word entries: %v", err)
	}

	chars, err = readEntries(EntriesFileChars)
	if err != nil {
		return nil, nil, fmt.Errorf("loading char entries: %v", err)
	}

	return words, chars, nil
}

func LoadWords(fileName string) (jparser, mecab, kanji []*Word, err error) {
	zipFile, err := files.FindZip(fileName)
	if err != nil {
		return nil, nil, nil, err
	}
	defer zipFile.Close()

	readWords := func(path string) (list []*Word, err error) {
		input, err := zipFile.OpenFileByName(path)
		if err == nil {
			defer input.Close()
			scanner := bufio.NewScanner(input)
			for err == nil && scanner.Scan() {
				// files have the byte order mark, so we need to strip it
				text := strings.TrimLeft(scanner.Text(), "\uFEFF")
				var word *Word
				word, err = ParseWord(text)
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

	kanji, err = readWords(WordsFileKanji)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("loading kanji entries: %v", err)
	}

	jparser, err = readWords(WordsFileJparser)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("loading jparser entries: %v", err)
	}

	mecab, err = readWords(WordsFileMecab)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("loading mecab entries: %v", err)
	}

	return jparser, mecab, kanji, nil
}
