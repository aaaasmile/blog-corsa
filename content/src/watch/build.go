package watch

import (
	"bytes"
	"corsa-blog/conf"
	"corsa-blog/content/src/mhproc"
	"corsa-blog/db"
	"corsa-blog/idl"
	"corsa-blog/util"
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type Builder struct {
	mdsFn    []string
	pages    []string
	liteDB   *db.LiteDB
	tx       *sql.Tx
	mapLinks *idl.MapPostsLinks
	force    bool
	debug    bool
}

func RebuildAll() error {
	start := time.Now()
	log.Println("[RebuildAll] the full site")

	bb := Builder{force: true}
	if err := bb.InitDBData(); err != nil {
		return err
	}
	if err := bb.scanMdHtml("../posts-src"); err != nil {
		return err
	}
	var err error
	if bb.mapLinks, err = CreateMapLinks(bb.liteDB); err != nil {
		return err
	}
	if err := bb.rebuildPosts("../posts-src"); err != nil {
		return err
	}
	if err := bb.rebuildFeed(); err != nil {
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

func BuildFeed() error {
	start := time.Now()
	log.Println("[BuildFeed] start")

	bb := Builder{}
	if err := bb.InitDBData(); err != nil {
		return err
	}

	if err := bb.rebuildFeed(); err != nil {
		return err
	}
	log.Println("[BuildFeed] completed, elapsed time ", time.Since(start))
	return nil
}

func BuildPosts() error {
	start := time.Now()
	log.Println("[BuildPosts] start")

	bb := Builder{
		debug: conf.Current.Debug,
	}
	if err := bb.InitDBData(); err != nil {
		return err
	}

	if err := bb.rebuildPosts("../posts-src"); err != nil {
		return err
	}
	if err := bb.rebuildFeed(); err != nil {
		return err
	}
	log.Println("[BuildPosts] completed, elapsed time ", time.Since(start))
	return nil
}

func BuildPages() error {
	start := time.Now()
	log.Println("[BuildPages] start")

	bb := Builder{}
	if err := bb.InitDBData(); err != nil {
		return err
	}
	if err := bb.rebuildPages("../page-src"); err != nil {
		return err
	}

	log.Println("[BuildPages] completed, elapsed time ", time.Since(start))
	return nil
}

func BuildMain() error {
	// TODO integrate this in build page
	start := time.Now()
	log.Println("[BuildMain] started")

	bb := Builder{}
	if err := bb.InitDBData(); err != nil {
		return err
	}
	if err := bb.rebuildMainPage(); err != nil {
		return err
	}

	log.Println("[BuildMain] completed, elapsed time ", time.Since(start))
	return nil
}

type PostWithData struct {
	DateFormatted string
	DateTime      string
	Title         string
	Link          string
}

func (bb *Builder) InitDBData() error {
	var err error
	if bb.liteDB, err = db.OpenSqliteDatabase(fmt.Sprintf("..\\..\\%s", conf.Current.Database.DbFileName),
		conf.Current.Database.SQLDebug); err != nil {
		return err
	}
	if bb.mapLinks, err = CreateMapLinks(bb.liteDB); err != nil {
		return err
	}
	return nil
}

func (bb *Builder) rebuildMainPage() error {
	log.Println("[rebuildMainPage] start")
	templDir := "templates/htmlgen"
	templName := path.Join(templDir, "mainpage.html")
	var partFirst bytes.Buffer
	tmplPage := template.Must(template.New("Page").ParseFiles(templName))
	latestPosts := []*PostWithData{}
	for ix, item := range bb.mapLinks.ListPost {
		pwd := PostWithData{
			DateFormatted: util.FormatDateIt(item.DateTime),
			DateTime:      item.DateTime.Format("2006-01-02 15:00"),
			Title:         item.Title,
			Link:          item.Uri,
		}
		latestPosts = append(latestPosts, &pwd)
		if ix >= 7 {
			break
		}
	}
	CtxFirst := struct {
		Title       string
		LatestPosts []*PostWithData
	}{
		Title:       "IgorRun Blog",
		LatestPosts: latestPosts,
	}

	if err := tmplPage.ExecuteTemplate(&partFirst, "mainpage", CtxFirst); err != nil {
		return err
	}
	prc := mhproc.NewMdHtmlProcess(false, nil)
	prc.RootStaticDir = fmt.Sprintf("..\\..\\static\\%s", conf.Current.StaticBlogDir)
	prc.HtmlGen = partFirst.String()
	prc.TargetDir = prc.RootStaticDir
	err := prc.CreateOnlyIndexStaticHtml()

	return err
}

func (bb *Builder) rebuildFeed() error {
	log.Println("[rebuildFeed] start ")
	templDir := "templates/xml"
	templName := path.Join(templDir, "feed.xml")
	var partFirst bytes.Buffer
	tmplPage := template.Must(template.New("FeedSrc").ParseFiles(templName))

	if err := tmplPage.ExecuteTemplate(&partFirst, "feedbeg", bb.mapLinks.ListPost); err != nil {
		return err
	}
	rootStaticDir := fmt.Sprintf("..\\..\\static\\%s\\", conf.Current.StaticBlogDir)
	fname := path.Join(rootStaticDir, "feed")
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(partFirst.Bytes()); err != nil {
		return err
	}
	log.Println("file created ", fname)

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
			if bb.debug {
				log.Println("[buildItem] ignore because unchanged", mdHtmlFname)
			}
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
