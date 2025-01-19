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
// go run .\main.go -config ..\..\config.toml  -editpost -date "2023-01-04"
//
// Example new post:
//
//	go run .\main.go -config ..\..\config.toml  -newpost "Quo Vadis" -date "2023-01-04"
//
// # Example edit page
//
// go run .\main.go -config ..\..\config.toml  -editpage -name "autore"
// Example new page:
//
// go run .\main.go -config ..\..\config.toml  -newpage "statistiche" -date "2025-01-18"
// Example build all
// go run .\main.go -config ..\..\config.toml  -build
func main() {
	var ver = flag.Bool("ver", false, "Prints the current version")
	var configfile = flag.String("config", "config.toml", "Configuration file path")
	var watchdir = flag.Bool("watch", false, "Watch the mdhtml file and generate the html")
	var newpost = flag.String("newpost", "", "title of the new post")
	var date = flag.String("date", "", "Date of the post, e.g. 2025-09-30")
	var editpost = flag.Bool("editpost", false, "edit post at date")
	var editpage = flag.Bool("editpage", false, "edit page at name")
	var newpage = flag.String("newpage", "", "name of the new page")
	var name = flag.String("name", "", "name of the page")
	var build = flag.Bool("build", false, "create all htmls (post and pages)")
	flag.Parse()

	if *ver {
		fmt.Printf("%s, version: %s", idl.Appname, idl.Buildnr)
		os.Exit(0)
	}
	if _, err := conf.ReadConfig(*configfile); err != nil {
		log.Fatal("ERROR: ", err)
	}
	if *build {
		if err := watch.Build(); err != nil {
			log.Fatal("ERROR: ", err)
		}
	} else if *editpost {
		if err := watch.EditPost(*date); err != nil {
			log.Fatal("ERROR: ", err)
		}
	} else if *newpost != "" {
		if err := watch.NewPost(*newpost, *date, *watchdir); err != nil {
			log.Fatal("ERROR: ", err)
		}
	} else if *newpage != "" {
		if err := watch.NewPage(*newpage, *date); err != nil {
			log.Fatal("ERROR: ", err)
		}
	} else if *editpage {
		if err := watch.EditPage(*name); err != nil {
			log.Fatal("ERROR: ", err)
		}
	}
	log.Println("That' all folks!")
}
