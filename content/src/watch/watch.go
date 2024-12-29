package watch

import (
	"corsa-blog/conf"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/fsnotify/fsnotify"
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
		case err, ok := <-watcher.Errors:
			if !ok {
				return err
			}
			log.Println("error:", err)
		}
	}
}
