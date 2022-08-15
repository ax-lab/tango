package db

import (
	"database/sql"
)

func Open(databaseFile string, schemaSQL string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databaseFile)
	if err == nil {
		if schemaSQL != "" {
			_, err = db.Exec(schemaSQL)
		}
	}
	if err != nil && db != nil {
		db.Close()
		db = nil
	}
	return db, err
}

func BeginTransaction(db *sql.DB) *WriterTransaction {
	tx, err := db.Begin()
	return &WriterTransaction{tx, err}
}

type WriterTransaction struct {
	tx  *sql.Tx
	err error
}

func (trans *WriterTransaction) Finish() error {
	if trans.err == nil {
		trans.err = trans.tx.Commit()
	} else if trans.tx != nil {
		trans.tx.Rollback()
	}
	return trans.err
}

func (trans *WriterTransaction) Prepare(sql string) *WriterCommand {
	out := &WriterCommand{trans: trans}
	if trans.err == nil {
		out.stmt, trans.err = trans.tx.Prepare(sql)
	}
	return out
}

type WriterCommand struct {
	trans *WriterTransaction
	stmt  *sql.Stmt
}

func (cmd *WriterCommand) Exec(args ...any) bool {
	if cmd.trans.err == nil {
		_, cmd.trans.err = cmd.stmt.Exec(args...)
	}
	return cmd.trans.err == nil
}
