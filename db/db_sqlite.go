package db

import (
	"corsa-blog/idl"
	"corsa-blog/util"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

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

func (ld *LiteDB) GetCommentForId(id string) (*idl.CmtNode, error) {
	log.Println("[GetCommentForId] get comment id ", id)
	q := `SELECT id,parent_id,post_id,name,email,comment,timestamp,status from comment where id = ?;`
	if ld.debugSQL {
		log.Println("Query is", q)
	}
	rows, err := ld.connDb.Query(q, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var rowid int
	var parent_id int
	post_id := ""
	level0_ids := []int{}
	arrCmtItem := []*idl.CmtItem{}
	for rows.Next() {
		var ts int64
		statustxt := ""
		cmtItem := idl.CmtItem{}
		if err := rows.Scan(&rowid, &parent_id, &cmtItem.PostId, &cmtItem.Name, &cmtItem.Email, &cmtItem.Comment, &ts, &statustxt); err != nil {
			return nil, err
		}
		cmtItem.Id = rowid
		cmtItem.ParentId = parent_id
		cmtItem.DateTime = time.Unix(ts, 0)
		status, err := strconv.Atoi(statustxt)
		if err != nil {
			return nil, err
		}
		cmtItem.Status = idl.StatusType(status)
		arrCmtItem = append(arrCmtItem, &cmtItem)
		level0_ids = append(level0_ids, rowid)
		if post_id == "" {
			post_id = cmtItem.PostId
		}
	}
	if len(arrCmtItem) == 0 {
		return nil, fmt.Errorf("comment id %s not found", id)
	}
	if len(arrCmtItem) > 1 {
		return nil, fmt.Errorf("comment id %s multiple instance?", id)
	}
	root_node := &idl.CmtNode{
		Children: []*idl.CmtNode{},
		CmtItem:  arrCmtItem[0],
		PostId:   post_id,
	}
	level := 0
	for ix, item_id := range level0_ids {
		node := &idl.CmtNode{
			PostId:    post_id,
			Children:  []*idl.CmtNode{},
			CmtItem:   arrCmtItem[ix],
			NodeCount: 1,
		}
		node.CmtItem = arrCmtItem[ix]
		children, err := ld.getCommentNodeChildren(level, item_id, post_id)
		if err != nil {
			return nil, err
		}
		if len(children) > 0 {
			node.Children = append(node.Children, children...)
		}

		root_node.Children = append(root_node.Children, node)
		root_node.NodeCount += node.NodeCount
	}
	log.Println("[GetCommentForId] found level 0 items: ", len(level0_ids))

	return root_node, nil
}

func (ld *LiteDB) GeCommentsForPostId(post_id string) (*idl.CmtNode, error) {
	log.Println("[LiteDB-SELECT] get comments for post id ", post_id)
	q := `SELECT id,name,email,comment,timestamp,status from comment where post_id = ? and parent_id = 0;`
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
		CmtItem:  &idl.CmtItem{},
	}
	var rowid int
	level0_ids := []int{}
	arrCmtItem := []*idl.CmtItem{}
	for rows.Next() {
		var ts int64
		statustxt := ""
		cmtItem := idl.CmtItem{}
		if err := rows.Scan(&rowid, &cmtItem.Name, &cmtItem.Email, &cmtItem.Comment, &ts, &statustxt); err != nil {
			return nil, err
		}
		cmtItem.Id = rowid
		cmtItem.PostId = post_id
		cmtItem.DateTime = time.Unix(ts, 0)
		status, err := strconv.Atoi(statustxt)
		if err != nil {
			return nil, err
		}
		cmtItem.Status = idl.StatusType(status)
		arrCmtItem = append(arrCmtItem, &cmtItem)
		level0_ids = append(level0_ids, rowid)
	}
	level := 0
	for ix, item_id := range level0_ids {
		node := &idl.CmtNode{
			PostId:    post_id,
			Children:  []*idl.CmtNode{},
			CmtItem:   arrCmtItem[ix],
			NodeCount: 1,
		}
		children, err := ld.getCommentNodeChildren(level, item_id, post_id)
		if err != nil {
			return nil, err
		}
		if len(children) > 0 {
			node.Children = append(node.Children, children...)
			for _, item := range children {
				node.NodeCount += item.NodeCount
			}
		}
		root_node.Children = append(root_node.Children, node)
		root_node.NodeCount += node.NodeCount
	}
	log.Println("[LiteDB-SELECT] found level 0 items: ", len(level0_ids))

	return root_node, nil
}

func (ld *LiteDB) getCommentNodeChildren(level int, parent_id int, post_id string) ([]*idl.CmtNode, error) {
	log.Println("[getCommentNodeChildren] level ", level)
	q := `SELECT id,name,email,comment,timestamp,status from comment where post_id = ? and parent_id = ?;`
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
	arrCmtItem := []*idl.CmtItem{}

	for rows.Next() {
		var ts int64
		cmtItem := idl.CmtItem{}
		statustxt := ""
		if err := rows.Scan(&rowid, &cmtItem.Name, &cmtItem.Email, &cmtItem.Comment, &ts, &statustxt); err != nil {
			return nil, err
		}
		cmtItem.Id = rowid
		cmtItem.ParentId = parent_id
		cmtItem.PostId = post_id
		cmtItem.DateTime = time.Unix(ts, 0)
		cmtItem.Indent = level
		status, err := strconv.Atoi(statustxt)
		if err != nil {
			return nil, err
		}
		cmtItem.Status = idl.StatusType(status)
		arrCmtItem = append(arrCmtItem, &cmtItem)
		level_ids = append(level_ids, rowid)
	}
	nex_level := level + 1
	nodes := []*idl.CmtNode{}
	subNodeCount := 0
	for ix, item_id := range level_ids {
		node := &idl.CmtNode{
			PostId:    post_id,
			Children:  []*idl.CmtNode{},
			NodeCount: 0,
		}
		node.CmtItem = arrCmtItem[ix]
		node.NodeCount += 1

		children, err := ld.getCommentNodeChildren(nex_level, item_id, post_id)
		if err != nil {
			return nil, err
		}
		if len(children) > 0 {
			node.Children = append(node.Children, children...)
			for _, item := range children {
				node.NodeCount += item.NodeCount
				subNodeCount += item.NodeCount
			}
		}
		nodes = append(nodes, node)
	}
	log.Printf("[getCommentNodeChild] on level %d found %d nodes with parent id %d, sub-count %d", level, len(nodes), parent_id, subNodeCount)
	return nodes, nil
}
