package db

import (
	"database/sql"
	"sort"
	"strings"

	"github.com/ax-lab/tango/import/frequency"
)

type FrequencyWriter struct {
	db *sql.DB

	jparserPairs []*frequency.Pair
	mecabPairs   []*frequency.Pair
	charPairs    []*frequency.Pair

	wordInfo []*frequency.Info
	charInfo []*frequency.Info
}

func NewFrequencyWriter(outputFile string) (*FrequencyWriter, error) {
	db, err := Open(outputFile, `
		DROP TABLE IF EXISTS freq_char;
		CREATE TABLE freq_char (
			entry     TEXT PRIMARY KEY,
			count     INTEGER,
			blog_freq INTEGER,
			twit_freq INTEGER,
			news_freq INTEGER,
			blog_pm   TEXT,
			twit_pm   TEXT,
			news_pm   TEXT,
			blog_cd   INTEGER,
			twit_cd   INTEGER,
			news_cd   INTEGER,
			blog_cdp  TEXT,
			twit_cdp  TEXT,
			news_cdp  TEXT
		);

		DROP TABLE IF EXISTS freq_word;
		CREATE TABLE freq_word (
			entry     TEXT PRIMARY KEY,
			count     INTEGER,
			count_m   INTEGER,
			blog_freq INTEGER,
			twit_freq INTEGER,
			news_freq INTEGER,
			blog_pm   TEXT,
			twit_pm   TEXT,
			news_pm   TEXT,
			blog_cd   INTEGER,
			twit_cd   INTEGER,
			news_cd   INTEGER,
			blog_cdp  TEXT,
			twit_cdp  TEXT,
			news_cdp  TEXT
		);
	`)
	if err != nil {
		return nil, err
	}

	return &FrequencyWriter{
		db: db,
	}, nil
}

func (writer *FrequencyWriter) Close() {
	writer.db.Close()
}

func (writer *FrequencyWriter) AddWordPairs(jparser, mecab []*frequency.Pair) {
	writer.jparserPairs = append(writer.jparserPairs, jparser...)
	writer.mecabPairs = append(writer.mecabPairs, mecab...)
}

func (writer *FrequencyWriter) AddCharPairs(items []*frequency.Pair) {
	writer.charPairs = append(writer.charPairs, items...)
}

func (writer *FrequencyWriter) AddWordInfo(items []*frequency.Info) {
	writer.wordInfo = append(writer.wordInfo, items...)
}

func (writer *FrequencyWriter) AddCharInfo(items []*frequency.Info) {
	writer.charInfo = append(writer.charInfo, items...)
}

func (writer *FrequencyWriter) Write() error {
	tx := BeginTransaction(writer.db)

	writer.writeTable(tx, "word", writer.jparserPairs, writer.mecabPairs, writer.wordInfo)
	writer.writeTable(tx, "char", writer.charPairs, nil, writer.charInfo)

	writer.jparserPairs = nil
	writer.mecabPairs = nil
	writer.charPairs = nil
	writer.wordInfo = nil
	writer.charInfo = nil

	return tx.Finish()
}

func (writer *FrequencyWriter) writeTable(
	tx *WriterTransaction, name string,
	pairs []*frequency.Pair, extraPairs []*frequency.Pair,
	info []*frequency.Info,
) {
	entries := make(map[string]map[string]interface{})

	update := func(entry string, keys map[string]interface{}) {
		row := entries[entry]
		if row == nil {
			row = make(map[string]interface{})
		}
		for k, v := range keys {
			row[k] = v
		}
		entries[entry] = row
	}

	setData := func(prefix string, entry map[string]interface{}, data frequency.InfoData) {
		entry[prefix+"_freq"] = data.Freq
		entry[prefix+"_pm"] = data.FreqPM
		entry[prefix+"_cd"] = data.CD
		entry[prefix+"_cdp"] = data.CDPc
	}

	for _, it := range pairs {
		update(it.Entry, map[string]interface{}{"count": it.Count})
	}

	for _, it := range extraPairs {
		update(it.Entry, map[string]interface{}{"count_m": it.Count})
	}

	for _, it := range info {
		entry := make(map[string]interface{})
		setData("blog", entry, it.Blog)
		setData("twit", entry, it.Twitter)
		setData("news", entry, it.News)
		update(it.Entry, entry)
	}

	rows := make([]map[string]interface{}, 0, len(entries))
	for key, val := range entries {
		val["entry"] = key
		rows = append(rows, val)
	}

	sort.Slice(rows, func(a, b int) bool {
		ca, cb := rows[a]["count"], rows[b]["count"]
		ca2, cb2 := rows[a]["count_m"], rows[b]["count_m"]
		if ca == nil {
			ca = ca2
		}
		if cb == nil {
			cb = cb2
		}
		if ca != nil {
			if cb != nil {
				return ca.(int64) > cb.(int64)
			} else {
				return true
			}
		} else if cb != nil {
			return false
		} else {
			fa, fb := rows[a]["blog_freq"], rows[b]["blog_freq"]
			return fa.(int64) > fb.(int64)
		}
	})

	cols := []string{
		"entry",
		"count",
		"blog_freq",
		"twit_freq",
		"news_freq",
		"blog_pm",
		"twit_pm",
		"news_pm",
		"blog_cd",
		"twit_cd",
		"news_cd",
		"blog_cdp",
		"twit_cdp",
		"news_cdp",
	}
	if extraPairs != nil {
		cols = append(cols, "count_m")
	}

	colNames := strings.Join(cols, ", ")
	insertRow := tx.Prepare(
		"INSERT INTO freq_" + name + "(" + colNames + ") VALUES (?" +
			strings.Repeat(", ?", len(cols)-1) + ")")

	args := make([]any, len(cols))
	for _, row := range rows {
		for index, col := range cols {
			args[index] = row[col]
		}
		insertRow.Exec(args...)
	}
}
