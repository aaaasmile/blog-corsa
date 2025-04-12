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

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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
	bb.liteDB.DeleteAllPostItem(tx)

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
	prc := mhproc.NewMdHtmlProcess(false, nil)
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
	subDir := conf.Current.PostSubDir
	arr, err := mhproc.GetDirNameArray(mdHtmlFname)
	if err != nil {
		return err
	}
	last_ix := len(arr) - 1
	dir_stack := []string{arr[last_ix-3], arr[last_ix-2], arr[last_ix-1], arr[last_ix]}
	remain := strings.Join(dir_stack, "/")
	postItem.Uri = fmt.Sprintf("/%s/%s/#", subDir, remain)
	//fmt.Println("*** uri is ", postItem.Uri)
	//fmt.Println("*** title is ", postItem.Title)
	bufRead := strings.NewReader(prc.HtmlGen)
	doc, err := html.Parse(bufRead)
	if err != nil {
		return err
	}
	traverse(doc, &postItem)

	err = bb.liteDB.InsertNewPost(tx, &postItem)
	if err != nil {
		return err
	}

	return nil
}
func traverse(doc *html.Node, postItem *idl.PostItem) {
	// We need here the title, abstract and header image
	// Information are from parsing the mdhtml file
	section_first := false
	title_first := false
	has_title_img := false
	for n := range doc.Descendants() {
		if !title_first && n.Type == html.ElementNode && n.DataAtom == atom.Header {
			for _, a := range n.Attr {
				if a.Key == "class" {
					if a.Val == "withimg" {
						//fmt.Println("** has an image in title ")
						has_title_img = true
					}
					break
				}
			}
		}
		if !title_first && n.Type == html.ElementNode && n.DataAtom == atom.H1 {
			if n.FirstChild != nil {
				title := n.FirstChild.Data
				//fmt.Println("** title ", title)
				postItem.Title = title
			}
			title_first = true
		}
		if has_title_img && n.Type == html.ElementNode && n.DataAtom == atom.Img {
			has_title_img = false
			for _, a := range n.Attr {
				if a.Key == "src" {
					img_src := a.Val
					//fmt.Println("*** image in title ", img_src)
					postItem.TitleImgUri = strings.TrimRight(postItem.Uri, "/#")
					postItem.TitleImgUri = fmt.Sprintf("%s/%s", postItem.TitleImgUri, img_src)
					//fmt.Println("*** TitleImgUri ", postItem.TitleImgUri)
					break
				}
			}
		}
		if n.Type == html.ElementNode && n.DataAtom == atom.Section {
			section_first = true
			has_title_img = false
		}
		if section_first && n.Type == html.ElementNode && n.DataAtom == atom.P {
			if n.FirstChild != nil {
				abstract := n.FirstChild.Data
				abstract = strings.Trim(abstract, " ")
				abstract = strings.Trim(abstract, "\n")
				abstract = strings.Trim(abstract, " ")
				maxlen := 40
				if len(abstract) > maxlen-3 {
					abstract = fmt.Sprintf("%s...", abstract[0:maxlen])
				}
				//fmt.Println("** abstract ", abstract)
				if len(abstract) > 4 {
					postItem.Abstract = abstract
				}
			}
			return
		}
	}
}
