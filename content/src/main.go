package main

import (
	"corsa-blog/conf"
	"corsa-blog/content/src/watch"
	"corsa-blog/idl"
	"flag"
	"fmt"
	"log"
	"os"
)

// Example edit a post:
//
//	go run .\main.go -config ..\..\config.toml  -watch -target ..\posts-src\2024\11\08\
//
// Another example of edit
// go run .\main.go -config ..\..\config.toml  -editpost -date "2023-01-04"
//
// Example new post:
//
//	go run .\main.go -config ..\..\config.toml  -newpost "Quo Vadis" -date "2023-01-04"
func main() {
	var ver = flag.Bool("ver", false, "Prints the current version")
	var configfile = flag.String("config", "config.toml", "Configuration file path")
	var watchdir = flag.Bool("watch", false, "Watch the mdhtml file and generate the html")
	var target = flag.String("target", "", "file to watch")
	var newpost = flag.String("newpost", "", "title of the new post")
	var date = flag.String("date", "", "Date of the post, e.g. 2025-09-30")
	var editpost = flag.Bool("editpost", false, "edit post at date")
	flag.Parse()

	if *ver {
		fmt.Printf("%s, version: %s", idl.Appname, idl.Buildnr)
		os.Exit(0)
	}
	if _, err := conf.ReadConfig(*configfile); err != nil {
		log.Fatal("ERROR: ", err)
	}
	if *editpost {
		if err := watch.EditPost(*date); err != nil {
			log.Fatal("ERROR: ", err)
		}
	} else if *newpost != "" {
		if err := watch.NewPost(*newpost, *date, *watchdir); err != nil {
			log.Fatal("ERROR: ", err)
		}
	} else if *watchdir {
		if err := watch.RunWatcher(*target); err != nil {
			log.Fatal("ERROR: ", err)
		}
	}
	log.Println("That' all folks!")
}
