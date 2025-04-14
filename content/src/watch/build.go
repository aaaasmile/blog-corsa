package watch

import (
	"corsa-blog/conf"
	"corsa-blog/content/src/mhproc"
	"corsa-blog/db"
	"corsa-blog/idl"
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Builder struct {
	mdsFn    []string
	pages    []string
	liteDB   *db.LiteDB
	tx       *sql.Tx
	mapLinks *idl.MapPostsLinks
	force    bool
}

func RebuildAll() error {
	start := time.Now()
	log.Println("[RebuildAll] the full site")

	bb := Builder{force: true}
	var err error
	if bb.liteDB, err = db.OpenSqliteDatabase(fmt.Sprintf("..\\..\\%s", conf.Current.Database.DbFileName),
		conf.Current.Database.SQLDebug); err != nil {
		return err
	}
	if bb.mapLinks, err = CreateMapLinks(bb.liteDB); err != nil {
		return err
	}
	if err := bb.rebuildPosts("../posts-src"); err != nil {
		return err
	}
	if err := bb.rebuildPages("../page-src"); err != nil {
		return err
	}
	if err := bb.rebuildMainPage(); err != nil {
		return err
	}
	log.Println("[RebuildAll] completed, elapsed time ", time.Since(start))
	return nil
}

func BuildPosts() error {
	start := time.Now()
	log.Println("[BuildPosts] changed posts")

	bb := Builder{}
	var err error
	if bb.liteDB, err = db.OpenSqliteDatabase(fmt.Sprintf("..\\..\\%s", conf.Current.Database.DbFileName),
		conf.Current.Database.SQLDebug); err != nil {
		return err
	}
	if bb.mapLinks, err = CreateMapLinks(bb.liteDB); err != nil {
		return err
	}
	if err := bb.rebuildPosts("../posts-src"); err != nil {
		return err
	}
	log.Println("[BuildPosts] completed, elapsed time ", time.Since(start))
	return nil
}

func (bb *Builder) rebuildMainPage() error {
	//TODO
	return nil
}

func (bb *Builder) rebuildPosts(srcDir string) error {
	bb.mdsFn = make([]string, 0)
	var err error
	bb.mdsFn, err = getFilesinDir(srcDir, bb.mdsFn)
	if err != nil {
		return err
	}
	bb.tx, err = bb.liteDB.GetTransaction()
	if err != nil {
		return err
	}
	log.Printf("%d mdhtml posts  found ", len(bb.mdsFn))
	for _, item := range bb.mdsFn {
		if err := bb.buildItem(item, false); err != nil {
			return err
		}
	}
	bb.tx.Commit()
	log.Printf("%d posts processed ", len(bb.mdsFn))
	return nil
}

func (bb *Builder) rebuildPages(srcDir string) error {
	bb.pages = make([]string, 0)
	var err error
	bb.pages, err = getFilesinDir(srcDir, bb.pages)
	if err != nil {
		return err
	}
	log.Printf("%d mdhtml pages found ", len(bb.pages))
	for _, item := range bb.pages {
		if err := bb.buildItem(item, true); err != nil {
			return err
		}
	}
	log.Printf("%d pages processed ", len(bb.pages))
	return nil
}

func getFilesinDir(dirAbs string, ini []string) ([]string, error) {
	r := ini
	//log.Println("Scan dir ", dirAbs)
	files, err := os.ReadDir(dirAbs)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		itemAbs := path.Join(dirAbs, f.Name())
		if info, err := os.Stat(itemAbs); err == nil && info.IsDir() {
			//fmt.Println("** Sub dir found ", f.Name())
			r, err = getFilesinDir(itemAbs, r)
			if err != nil {
				return nil, err
			}
		} else {
			//fmt.Println("** file is ", f.Name())
			ext := filepath.Ext(itemAbs)
			if strings.HasPrefix(ext, ".mdhtml") {
				r = append(r, path.Join(dirAbs, f.Name()))
			}
		}
	}
	return r, nil
}

func (bb *Builder) buildItem(mdHtmlFname string, is_page bool) error {
	var err error
	wmh := WatcherMdHtml{
		debug:         conf.Current.Debug,
		staticBlogDir: conf.Current.StaticBlogDir,
		is_page:       is_page,
		mapLinks:      bb.mapLinks,
	}
	is_same := true
	postItem := &idl.PostItem{}
	if is_page {
		wmh.staticSubDir = conf.Current.PageSubDir
	} else {
		postItem, is_same, err = bb.hasSameMd5(mdHtmlFname)
		if err != nil {
			return err
		}
		wmh.staticSubDir = conf.Current.PostSubDir
		if !bb.force && is_same {
			log.Println("[buildItem] ignore because unchanged", mdHtmlFname)
			return nil
		}
	}
	if err := wmh.BuildFromMdHtml(mdHtmlFname); err != nil {
		return err
	}
	if (postItem.PostId != "") && !is_same {
		if err := bb.liteDB.UpdateMd5Post(bb.tx, postItem); err != nil {
			return err
		}
	}
	log.Println("created HTML: ", wmh.CreatedHtmlFile)
	return nil
}

func (bb *Builder) hasSameMd5(mdHtmlFname string) (*idl.PostItem, bool, error) {
	f, err := os.Open(mdHtmlFname)
	if err != nil {
		return nil, false, err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, false, err
	}
	mMd5 := string(h.Sum(nil))

	mdhtml, err := os.ReadFile(mdHtmlFname)
	if err != nil {
		return nil, false, err
	}
	prc := mhproc.NewMdHtmlProcess(false, nil)
	if err := prc.ProcessToHtml(string(mdhtml)); err != nil {
		log.Println("[hasSameMd5] ProcessToHtml error: ", err)
		return nil, false, err
	}
	gr := prc.GetScriptGrammar()
	mMd5Db, ok := bb.mapLinks.MapPost[gr.PostId]
	if !ok {
		return nil, false, fmt.Errorf("[hasSameMd5] post id %s not found in MapLinks. Is the post table in db syncronized?", gr.PostId)
	}
	same := mMd5 == mMd5Db.Item.Md5
	postItem := idl.PostItem{PostId: gr.PostId, Md5: mMd5}
	return &postItem, same, nil

}
