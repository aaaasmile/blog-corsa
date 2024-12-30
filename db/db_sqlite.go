package db

import (
	"corsa-blog/util"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type LiteDB struct {
	connDb   *sql.DB
	DBPath   string
	DebugSQL bool
}

func (ld *LiteDB) OpenSqliteDatabase() error {
	var err error
	dbname := util.GetFullPath(ld.DBPath)
	log.Println("Using the sqlite file: ", dbname)
	ld.connDb, err = sql.Open("sqlite3", dbname)
	if err != nil {
		return err
	}

	return nil
}
