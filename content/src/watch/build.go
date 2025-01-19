package watch

import (
	"corsa-blog/conf"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Builder struct {
	mdsFn []string
}

func Build() error {
	start := time.Now()
	bb := Builder{}
	if err := bb.rebuildPosts("../posts-src"); err != nil {
		return err
	}
	log.Println("[Build] completed, elapsed time ", time.Since(start))
	return nil
}

func (bb *Builder) rebuildPosts(srcDir string) error {
	bb.mdsFn = make([]string, 0)
	var err error
	bb.mdsFn, err = getFilesinDir(srcDir, bb.mdsFn)
	if err != nil {
		return err
	}
	log.Printf("%d mdhtml posts  found ", len(bb.mdsFn))
	for _, item := range bb.mdsFn {
		if err := bb.buildItem(item, false); err != nil {
			return err
		}
	}
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
	wmh := WatcherMdHtml{
		debug:         conf.Current.Debug,
		staticBlogDir: conf.Current.StaticBlogDir,
		staticSubDir:  conf.Current.PostSubDir,
		is_page:       is_page,
	}
	if err := wmh.BuildFromMdHtml(mdHtmlFname); err != nil {
		return err
	}
	log.Println("create HTML: ", wmh.CreatedHtmlFile)
	return nil
}
