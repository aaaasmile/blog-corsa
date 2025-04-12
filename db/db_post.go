package db

import (
	"corsa-blog/idl"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func (ld *LiteDB) DeleteAllPostItem(tx *sql.Tx) error {
	q := `DELETE FROM post;`
	if ld.debugSQL {
		log.Println("SQL is:", q)
	}

	stm, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}
	res, err := tx.Stmt(stm).Exec()
	if ld.debugSQL {
		ra, err := res.RowsAffected()
		if err != nil {
			return err
		}
		log.Println("Row affected: ", ra)
	}
	return err
}

func (ld *LiteDB) InsertNewPost(tx *sql.Tx, postItem *idl.PostItem) error {
	log.Println("[LiteDB - InsertNewPost] insert new Post on post id ", postItem.PostId)

	q := `INSERT INTO post(title,post_id,timestamp,abstract,uri,title_img_uri) VALUES(?,?,?,?,?,?);`
	if ld.debugSQL {
		log.Println("Query is", q)
	}

	stmt, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}

	_, err = tx.Stmt(stmt).Exec(postItem.Title,
		postItem.PostId,
		postItem.DateTime.Local().Unix(),
		postItem.Abstract,
		postItem.Uri,
		postItem.TitleImgUri)
	if err != nil {
		return err
	}
	q = `SELECT last_insert_rowid()`
	if ld.debugSQL {
		log.Println("Query is", q)
	}
	var id int
	err = ld.connDb.QueryRow(q).Scan(&id)
	if err != nil {
		return err
	}
	postItem.Id = id
	log.Println("Post added into the db OK: ", postItem.Id)
	return nil
}
