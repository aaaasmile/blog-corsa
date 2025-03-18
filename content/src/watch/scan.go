package watch

import (
	"log"
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
	bb.mdsFn = make([]string, 0)
	var err error
	bb.mdsFn, err = getFilesinDir(srcDir, bb.mdsFn)
	if err != nil {
		return err
	}
	log.Printf("%d mdhtml posts  found ", len(bb.mdsFn))
	for _, item := range bb.mdsFn {
		if err := bb.scanItem(item, false); err != nil {
			return err
		}
	}
	log.Printf("%d posts processed ", len(bb.mdsFn))
	return nil
}

func (bb *Builder) scanItem(mdHtmlFname string, is_page bool) error {
	//TODO: implement
	return nil
}
