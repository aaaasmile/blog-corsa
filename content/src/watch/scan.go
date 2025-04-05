package watch

import (
	"corsa-blog/conf"
	"corsa-blog/content/src/mhproc"
	"corsa-blog/db"
	"corsa-blog/idl"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func ScanContent() error {
	start := time.Now()
	bb := Builder{}
	if err := bb.scanMdHtml("../posts-src"); err != nil {
		return err
	}
	log.Println("[ScanContent] completed, elapsed time ", time.Since(start))
	return nil
}

func (bb *Builder) scanMdHtml(srcDir string) error {
	var err error
	if bb.liteDB, err = db.OpenSqliteDatabase(fmt.Sprintf("..\\..\\%s", conf.Current.Database.DbFileName),
		conf.Current.Database.SQLDebug); err != nil {
		return err
	}
	bb.mdsFn = make([]string, 0)
	bb.mdsFn, err = getFilesinDir(srcDir, bb.mdsFn)
	if err != nil {
		return err
	}
	log.Printf("%d mdhtml posts  found ", len(bb.mdsFn))
	tx, err := bb.liteDB.GetTransaction()
	if err != nil {
		return err
	}

	for _, item := range bb.mdsFn {
		if err := bb.scanPostItem(item, tx); err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("%d posts processed ", len(bb.mdsFn))
	return nil
}

func (bb *Builder) scanPostItem(mdHtmlFname string, tx *sql.Tx) error {
	log.Println("[scanPostItem] file is ", mdHtmlFname)

	mdhtml, err := os.ReadFile(mdHtmlFname)
	if err != nil {
		return err
	}
	//log.Println("read: ", mdhtml)
	prc := mhproc.NewMdHtmlProcess(false)
	if err := prc.ProcessToHtml(string(mdhtml)); err != nil {
		log.Println("[scanPostItem] HTML error: ", err)
		return err
	}
	grm := prc.GetScriptGrammar()
	postItem := idl.PostItem{
		Title:    grm.Title,
		PostId:   grm.PostId,
		DateTime: grm.Datetime,
	}
	//staticBlogDir := conf.Current.StaticBlogDir
	subDir := conf.Current.PostSubDir
	arr, err := mhproc.GetDirNameArray(mdHtmlFname)
	if err != nil {
		return err
	}
	last_ix := len(arr) - 1
	dir_stack := []string{arr[last_ix-3], arr[last_ix-2], arr[last_ix-1], arr[last_ix]}
	remain := strings.Join(dir_stack, "/")
	postItem.Uri = fmt.Sprintf("/%s/%s/#", subDir, remain)
	fmt.Println("*** uri is ", postItem.Uri)

	return nil
}
