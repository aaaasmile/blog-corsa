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
//	go run .\main.go -config ..\..\config.toml  -newpost "Quo Vadis" -date "2023-01-04" -watch
//
// # Example edit page
//
// go run .\main.go -config ..\..\config.toml  -editpage -name "autore"
// Example new page:
//
// go run .\main.go -config ..\..\config.toml  -newpage "statistiche" -date "2025-01-18"
// Example rebuild all
// go run .\main.go -config ..\..\config.toml  -rebuildall
// Example build only changed posts
// go run .\main.go -config ..\..\config.toml  -buildposts
// Scan and update post info in db
// go run .\main.go -config ..\..\config.toml  -scancontent
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
	var rebuildall = flag.Bool("rebuildall", false, "force to create all htmls (main, post and pages)")
	var scancontent = flag.Bool("scancontent", false, "fill the db table with souce content")
	var buildposts = flag.Bool("buildposts", false, "create posts (only changed)")
	flag.Parse()

	if *ver {
		fmt.Printf("%s, version: %s", idl.Appname, idl.Buildnr)
		os.Exit(0)
	}
	if _, err := conf.ReadConfig(*configfile, `../../cert`); err != nil {
		log.Fatal("ERROR: ", err)
	}
	if *scancontent {
		if err := watch.ScanContent(); err != nil {
			log.Fatal("ERROR: ", err)
		}
		return
	}
	if *buildposts {
		if err := watch.BuildPosts(); err != nil {
			log.Fatal("ERROR: ", err)
		}
	}
	if *rebuildall {
		if err := watch.RebuildAll(); err != nil {
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
