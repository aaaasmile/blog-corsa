package watch

import (
	"bytes"
	"corsa-blog/conf"
	"fmt"
	"image"
	"image/jpeg"
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
	debug      bool
	dirContent string
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
		wwa := WatcherMdHtml{dirContent: targetDir,
			debug: conf.Current.Debug,
		}
		if err := wwa.doWatch(); err != nil {
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

func (wwa *WatcherMdHtml) doWatch() error {
	log.Println("setup watch on ", wwa.dirContent)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	err = watcher.Add(wwa.dirContent)
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
				}
			}
			if event.Has(fsnotify.Create) {
				if time.Since(lastWriteEv) > time.Duration(500)*time.Millisecond {
					log.Println("Create file:", event.Name)
					lastWriteEv = time.Now()
					if err := wwa.processNewImage(event.Name); err != nil {
						return err
					}
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

func (wwa *WatcherMdHtml) processNewImage(newFname string) error {
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
	if isJpeg {
		original_image, err := jpeg.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return err
		}
		ff := "your_image_resized.jpg_trf"
		ff_full := filepath.Join(wwa.dirContent, ff)
		output, _ := os.Create(ff_full)
		defer output.Close()
		dst := image.NewRGBA(image.Rect(0, 0, original_image.Bounds().Max.X/2, original_image.Bounds().Max.Y/2))
		draw.CatmullRom.Scale(dst, dst.Rect, original_image, original_image.Bounds(), draw.Over, nil)
		jpOpt := jpeg.Options{Quality: 100}
		jpeg.Encode(output, dst, &jpOpt)
		log.Println("image created: ", ff_full)
	}

	return nil
}
