package frequency

import (
	"archive/zip"
	"bufio"
	"fmt"
	"strings"

	"github.com/ax-lab/tango/import/files"
)

const (
	DefaultWordFrequency    = "vendor/data/frequency/Innocent_Novel_Analysis_120526.zip"
	DefaultFrequencyEntries = "vendor/data/frequency/Jap.Freq.2.zip"

	FrequencyCharFileName = "Jap.Char.Freq.2.txt"
	FrequencyWordFileName = "Jap.Freq.2.txt"
)

func LoadFrequencyEntries(fileName string) (words, chars []*Entry, err error) {
	var size int64
	file, err := files.Find(fileName)
	if err == nil {
		stat, statErr := file.Stat()
		err = statErr
		size = stat.Size()
	}

	if err != nil {
		return nil, nil, err
	}

	defer file.Close()
	inputZip, err := zip.NewReader(file, size)
	if err != nil {
		return nil, nil, err
	}

	readEntries := func(path string) (list []*Entry, err error) {
		input, err := inputZip.Open(path)
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

	words, err = readEntries(FrequencyWordFileName)
	if err != nil {
		return nil, nil, fmt.Errorf("loading word entries: %v", err)
	}

	chars, err = readEntries(FrequencyCharFileName)
	if err != nil {
		return nil, nil, fmt.Errorf("loading char entries: %v", err)
	}

	return words, chars, nil
}

func LoadWordFrequency(file string) (mecab, jparser []*Word, err error) {
	panic("nie")
}
