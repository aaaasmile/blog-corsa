package db

import (
	"corsa-blog/idl"
	"corsa-blog/util"
	"database/sql"
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

func (ld *LiteDB) GeCommentsForPostId(postid string) (*idl.CmtNode, error) {
	q := `SELECT id from comment where post_id = ? and parent_id = '';`
	rows, err := ld.connDb.Query(q, postid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	root_node := &idl.CmtNode{
		PostId:   postid,
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
		node, err := ld.getCommentNodeChild(level, item_id, postid)
		if err != nil {
			return nil, err
		}
		root_node.Children = append(root_node.Children, node)
	}

	return root_node, nil
}

func (ld *LiteDB) getCommentNodeChild(level int, parent_id int, post_id string) (*idl.CmtNode, error) {
	log.Println("[getCommentNodeChild] level ", level)
	child_node := &idl.CmtNode{
		PostId:   post_id,
		Children: []*idl.CmtNode{},
	}
	q := `SELECT id from comment where post_id = ? and parent_id = '';`
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
	}
	log.Printf("[getCommentNodeChild] on level %d found %d children with parent id %d", level, len(child_node.Children), parent_id)
	return child_node, nil
}
