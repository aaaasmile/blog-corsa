package watch

import (
	"corsa-blog/conf"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Post struct {
	Datetime      time.Time
	DatetimeOrig  string
	Title         string
	TitleCompress string
	mdhtmlName    string
	contentDir    string
	templDir      string
	postId        string
}

func EditPost(datepost string) error {
	post := Post{
		DatetimeOrig: datepost,
	}
	if err := post.setDateTimeFromString(datepost); err != nil {
		return err
	}
	if err := post.editPost("../posts-src"); err != nil {
		return err
	}
	return nil
}

func (pp *Post) editPost(targetRootDir string) error {
	log.Printf("[editPost] on '%s'", pp.Datetime)
	yy := fmt.Sprintf("%d", pp.Datetime.Year())
	mm := fmt.Sprintf("%02d", pp.Datetime.Month())
	dd := fmt.Sprintf("%02d", pp.Datetime.Day())
	contentDir := filepath.Join(targetRootDir, yy, mm, dd)
	log.Println("source post content dir ", contentDir)
	info, err := os.Stat(contentDir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("[editPost] expected dir on %s", contentDir)
	}
	if err := RunWatcher(contentDir, conf.Current.PostSubDir, false); err != nil {
		log.Println("[editPost] error on watch")
		return err
	}
	return nil
}
