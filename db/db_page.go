package db

import (
	"corsa-blog/idl"
	"database/sql"
	"log"
	"time"
)

func (ld *LiteDB) UpdateMd5Page(tx *sql.Tx, pageItem *idl.PageItem) error {
	log.Println("[LiteDB - UpdateMd5Page] update md5 on page id ", pageItem.PageId)
	q := `UPDATE page SET md5=? WHERE page_id=?;`
	if ld.debugSQL {
		log.Println("Query is", q)
	}
	stm, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}

	res, err := tx.Stmt(stm).Exec(pageItem.Md5, pageItem.PageId)
	if ld.debugSQL {
		ra, err := res.RowsAffected()
		if err != nil {
			return err
		}
		log.Println("Row affected: ", ra)
	}

	return err
}

func (ld *LiteDB) GetPageList() ([]idl.PageItem, error) {
	log.Println("[LiteDB - GetPageList] select all pages")

	q := `SELECT id,title,page_id,timestamp,uri,md5 from page ORDER BY page_id DESC;`
	if ld.debugSQL {
		log.Println("Query is", q)
	}
	rows, err := ld.connDb.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []idl.PageItem{}
	for rows.Next() {
		var ts int64
		item := idl.PageItem{}
		if err := rows.Scan(&item.Id,
			&item.Title,
			&item.PageId,
			&ts,
			&item.Uri,
			&item.Md5); err != nil {
			return nil, err
		}
		item.DateTime = time.Unix(ts, 0)
		item.DateTimeRfC822 = item.DateTime.Format(time.RFC822Z)
		res = append(res, item)
	}
	log.Printf("[LiteDB - GetPageList] posts read %d", len(res))
	return res, nil
}

func (ld *LiteDB) DeleteAllPageItem() error {
	log.Println("[LiteDB - DeleteAllPageItem] start")
	q := `DELETE FROM page;`
	if ld.debugSQL {
		log.Println("SQL is:", q)
	}

	stm, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}
	res, err := stm.Exec()
	if ld.debugSQL {
		ra, err := res.RowsAffected()
		if err != nil {
			return err
		}
		log.Println("Row affected: ", ra)
	}
	return err
}

func (ld *LiteDB) InsertNewPage(tx *sql.Tx, postItem *idl.PageItem) error {
	if ld.debugSQL {
		log.Println("[LiteDB - InsertNewPage] insert new Post on post id ", postItem.PageId)
	}

	q := `INSERT INTO page(title,page_id,timestamp,uri,md5) VALUES(?,?,?,?,?);`
	if ld.debugSQL {
		log.Println("Query is", q)
	}

	stmt, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}

	_, err = tx.Stmt(stmt).Exec(postItem.Title,
		postItem.PageId,
		postItem.DateTime.Local().Unix(),
		postItem.Uri,
		postItem.Md5)
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
	if ld.debugSQL {
		log.Println("Post added into the db OK: ", postItem.Id)
	}
	return nil
}
