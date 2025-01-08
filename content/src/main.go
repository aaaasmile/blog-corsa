package main

import (
	"corsa-blog/content/src/watch"
	"corsa-blog/idl"
	"flag"
	"fmt"
	"log"
	"os"
)

// Example command line and the target should be a directory with one post and multiple images:
//
//	go run .\main.go -config ..\..\config_custom.toml  -watch -target ..\posts-src\2024\11\08\
func main() {
	var ver = flag.Bool("ver", false, "Prints the current version")
	var configfile = flag.String("config", "config.toml", "Configuration file path")
	var ww = flag.Bool("watch", false, "Watch the mdhtml file and generate the html")
	var target = flag.String("target", "", "file to watch")
	flag.Parse()

	if *ver {
		fmt.Printf("%s, version: %s", idl.Appname, idl.Buildnr)
		os.Exit(0)
	}
	if *ww {
		if err := watch.RunWatcher(*configfile, *target); err != nil {
			log.Fatal("ERROR: ", err)
		}
	}
	log.Println("That' all folks!")
}
