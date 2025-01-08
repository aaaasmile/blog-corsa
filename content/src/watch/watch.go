package watch

import (
	"bytes"
	"corsa-blog/conf"
	"corsa-blog/content/src/mhproc"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/image/draw"
)

type WatcherMdHtml struct {
	debug         bool
	dirContent    string
	staticBlogDir string
	postSubDir    string
}

func RunWatcher(configfile string, targetDir string) error {
	if targetDir == "" {
		return fmt.Errorf("target dir is empty")
	}
	log.Println("watching ", targetDir)
	fs, err := os.Stat(targetDir)
	if err != nil {
		return err
	}
	if !fs.IsDir() {
		return fmt.Errorf("watch make sense only on a directory with content and images")
	}
	if _, err := conf.ReadConfig(configfile); err != nil {
		return err
	}

	chShutdown := make(chan struct{}, 1)
	go func(chs chan struct{}) {
		wmh := WatcherMdHtml{dirContent: targetDir,
			debug:         conf.Current.Debug,
			staticBlogDir: conf.Current.StaticBlogDir,
			postSubDir:    conf.Current.PostSubDir,
		}
		if err := wmh.doWatch(); err != nil {
			log.Println("Server is not watching anymore because: ", err)
		}
		log.Println("watch end")
		chs <- struct{}{}
	}(chShutdown)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	log.Println("Enter in server blocking loop")

loop:
	for {
		select {
		case <-sig:
			log.Println("stop because interrupt")
			break loop
		case <-chShutdown:
			log.Println("stop because service shutdown on watch")
			break loop
		}
	}

	log.Println("Bye, service")
	return nil
}

func (wmh *WatcherMdHtml) doWatch() error {
	log.Println("setup watch on ", wmh.dirContent)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	err = watcher.Add(wmh.dirContent)
	if err != nil {
		return err
	}

	lastWriteEv := time.Now()
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return fmt.Errorf("watch event failed")
			}
			//log.Println("event:", event)
			if event.Has(fsnotify.Write) {
				if time.Since(lastWriteEv) > time.Duration(500)*time.Millisecond {
					log.Println("WRITE modified file:", event.Name)
					lastWriteEv = time.Now()
					if err := wmh.processMdHtmlChange(event.Name); err != nil {
						return err
					}
				}
			}
			if event.Has(fsnotify.Create) {
				if time.Since(lastWriteEv) > time.Duration(500)*time.Millisecond {
					log.Println("Create file:", event.Name)
					if err := wmh.processNewImage(event.Name); err != nil {
						return err
					}
					lastWriteEv = time.Now() // important because a new jpg image is created
				}
			}
			if event.Has(fsnotify.Rename) {
				log.Println("Rename file:", event.Name) // remember that is followed by a create event
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return err
			}
			log.Println("error:", err)
		}
	}
}

func (wmh *WatcherMdHtml) processMdHtmlChange(newFname string) error {
	if wmh.staticBlogDir == "" {
		return fmt.Errorf("static blog dir config is empty")
	}
	if wmh.postSubDir == "" {
		return fmt.Errorf("post sub dir config is empty")
	}
	_, err := os.Stat(newFname)
	if err != nil {
		return err
	}
	ext := filepath.Ext(newFname)
	if !strings.HasPrefix(ext, ".mdhtml") {
		log.Println("file ignored", newFname)
		return nil
	}
	mdhtml, err := os.ReadFile(newFname)
	if err != nil {
		return err
	}
	//log.Println("read: ", mdhtml)
	prc := mhproc.NewMdHtmlProcess(false)
	if err := prc.ProcessToHtml(string(mdhtml)); err != nil {
		log.Println("HTML error: ", err)
		return nil
	}
	log.Println("html created with size: ", len(prc.HtmlGen))
	prc.RootStaticDir = fmt.Sprintf("..\\..\\static\\%s\\%s", wmh.staticBlogDir, wmh.postSubDir)
	if err = prc.CreateOrUpdateStaticHtml(newFname); err != nil {
		return err
	}
	return nil
}

func (wmh *WatcherMdHtml) processNewImage(newFname string) error {
	_, err := os.Stat(newFname)
	if err != nil {
		return err
	}

	ext := filepath.Ext(newFname)
	log.Println("extension new file ", ext)
	isPng := strings.HasPrefix(ext, ".png")
	isJpeg := strings.HasPrefix(ext, ".jpg")
	if !(isJpeg || isPng) {
		log.Println("file ignored", newFname)
		return nil
	}

	imageBytes, err := os.ReadFile(newFname)
	if err != nil {
		return err
	}
	newWidth := 320
	base_ff := filepath.Base(newFname)
	ff := strings.Replace(base_ff, ext, "", 1)
	if isJpeg {
		ff = fmt.Sprintf("%s_%d.jpg", ff, newWidth)
	} else if isPng {
		ff = fmt.Sprintf("%s_%d.png", ff, newWidth)
	} else {
		return fmt.Errorf("image format %s not supported", ext)
	}
	ff_full := filepath.Join(wmh.dirContent, ff)

	var original_image image.Image
	if isJpeg {
		if original_image, err = jpeg.Decode(bytes.NewReader(imageBytes)); err != nil {
			return err
		}
	} else if isPng {
		if original_image, err = png.Decode(bytes.NewReader(imageBytes)); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("image format %s not supported", ext)
	}
	if original_image.Bounds().Max.X <= newWidth {
		log.Println("image is already on resize width or smaller", newWidth)
		return nil
	}

	output, _ := os.Create(ff_full)
	defer output.Close()
	log.Println("current image size ", original_image.Bounds().Max)
	ratiof := float32(original_image.Bounds().Max.X) / float32(newWidth)
	if ratiof == 0.0 {
		return fmt.Errorf("invalid source image, attempt division by zero")
	}
	newHeightf := float32(original_image.Bounds().Max.Y) / ratiof
	newHeight := int(newHeightf)
	log.Printf("new rect width %d height %d ratio %f ", newWidth, newHeight, ratiof)
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(dst, dst.Rect, original_image, original_image.Bounds(), draw.Over, nil)
	if isJpeg {
		jpOpt := jpeg.Options{Quality: 100}
		if err = jpeg.Encode(output, dst, &jpOpt); err != nil {
			return err
		}
	} else if isPng {
		if err = png.Encode(output, dst); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("image format %s not supported", ext)
	}
	log.Println("image created: ", ff_full)

	return nil
}
