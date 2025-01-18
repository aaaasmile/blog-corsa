package watch

import (
	"corsa-blog/conf"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Page struct {
	Datetime     time.Time
	DatetimeOrig string
	Name         string
	NameCompress string
	mdhtmlName   string
	contentDir   string
	templDir     string
	Id           string
}

func EditPage(name string) error {
	if name == "" {
		return fmt.Errorf("[EditPage] page name could not be empty")
	}
	page := Page{
		Name: name,
	}
	if err := page.editPage("../page-src"); err != nil {
		return err
	}

	return nil
}

func (pg *Page) editPage(targetRootDir string) error {
	log.Printf("[editPage] on '%s'", pg.Name)
	contentDir := filepath.Join(targetRootDir, pg.Name)
	log.Println("source page content dir ", contentDir)
	log.Println("destination page is ", conf.Current.PageSubDir)
	info, err := os.Stat(contentDir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("[editPage] expected dir on %s", contentDir)
	}
	if err := RunWatcher(contentDir, conf.Current.PageSubDir, true); err != nil {
		log.Println("[editPage] error on watch")
		return err
	}
	return nil
}
