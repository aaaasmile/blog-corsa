package db

import (
	"corsa-blog/idl"
	"corsa-blog/util"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type LiteDB struct {
	connDb   *sql.DB
	dBPath   string
	debugSQL bool
}

func OpenSqliteDatabase(dbPath string, debugSql bool) (*LiteDB, error) {
	ld := &LiteDB{
		dBPath:   dbPath,
		debugSQL: debugSql,
	}
	if err := ld.openSqliteDatabase(); err != nil {
		log.Println("[OpenSqliteDatabase] error")
		return nil, err
	}
	return ld, nil
}

func (ld *LiteDB) openSqliteDatabase() error {
	var err error
	dbname := util.GetFullPath(ld.dBPath)
	log.Println("Using the sqlite file: ", dbname)
	ld.connDb, err = sql.Open("sqlite3", dbname)
	if err != nil {
		return err
	}
	return nil
}
func (ld *LiteDB) DeleteComment(cmtItem *idl.CmtItem) error {
	log.Println("[LiteDB - DELETE] delete comment on id ", cmtItem.Id)
	q := fmt.Sprintf(`DELETE FROM comment WHERE id=%d AND req_id='%s';`, cmtItem.Id, cmtItem.ReqId)
	if ld.debugSQL {
		log.Println("Query is", q)
	}
	stmt, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	log.Println("[LiteDB - DELETE] ok")
	return nil
}

func (ld *LiteDB) InsertNewComment(cmtItem *idl.CmtItem) error {
	log.Println("[LiteDB - INSERT] insert new comment on post id ", cmtItem.PostId)

	q := `INSERT INTO comment(parent_id,name,email,comment,timestamp,post_id,status,req_id) VALUES(?,?,?,?,?,?,?,?);`
	if ld.debugSQL {
		log.Println("Query is", q)
	}

	stmt, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(cmtItem.ParentId,
		cmtItem.Name,
		cmtItem.Email,
		cmtItem.Comment,
		cmtItem.DateTime.Local().Unix(),
		cmtItem.PostId,
		cmtItem.Status,
		cmtItem.ReqId)
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
	cmtItem.Id = id
	log.Println("Comment added into the db OK: ", cmtItem.Id)
	return nil
}

func (ld *LiteDB) GeCommentsForPostId(post_id string) (*idl.CmtNode, error) {
	log.Println("[LiteDB-SELECT] get comments for post id ", post_id)
	q := `SELECT id from comment where post_id = ? and parent_id = 0;`
	if ld.debugSQL {
		log.Println("Query is", q)
	}
	rows, err := ld.connDb.Query(q, post_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	root_node := &idl.CmtNode{
		PostId:   post_id,
		Children: []*idl.CmtNode{},
	}
	var rowid int
	level0_ids := []int{}
	for rows.Next() {
		if err := rows.Scan(&rowid); err != nil {
			return nil, err
		}
		level0_ids = append(level0_ids, rowid)
	}
	level := 0
	for _, item_id := range level0_ids {
		node, err := ld.getCommentNodeChild(level, item_id, post_id)
		if err != nil {
			return nil, err
		}
		root_node.Children = append(root_node.Children, node)
		root_node.NodeCount += node.NodeCount
	}
	log.Println("[LiteDB-SELECT] found level 0 items: ", len(level0_ids))

	return root_node, nil
}

func (ld *LiteDB) getCommentNodeChild(level int, parent_id int, post_id string) (*idl.CmtNode, error) {
	log.Println("[getCommentNodeChild] level ", level)
	child_node := &idl.CmtNode{
		PostId:    post_id,
		Children:  []*idl.CmtNode{},
		NodeCount: 1,
	}
	q := `SELECT id from comment where post_id = ? and parent_id = ?;`
	if ld.debugSQL {
		log.Println("Query is", q)
	}
	rows, err := ld.connDb.Query(q, post_id, parent_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rowid int
	level_ids := []int{}
	for rows.Next() {
		if err := rows.Scan(&rowid); err != nil {
			return nil, err
		}
		level_ids = append(level_ids, rowid)
	}
	nex_level := level + 1
	for _, item_id := range level_ids {
		node, err := ld.getCommentNodeChild(nex_level, item_id, post_id)
		if err != nil {
			return nil, err
		}
		child_node.Children = append(child_node.Children, node)
		child_node.NodeCount += node.NodeCount
	}
	log.Printf("[getCommentNodeChild] on level %d found %d children with parent id %d, sub-count %d", level, len(child_node.Children), parent_id, child_node.NodeCount)
	return child_node, nil
}
