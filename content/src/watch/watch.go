package watch

import (
	"corsa-blog/conf"
	"fmt"
	"log"
	"os"
	"os/signal"

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
	if _, err := os.Stat(targetDir); err != nil {
		return err
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

	return nil
}
