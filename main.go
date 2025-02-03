package main

import (
	"corsa-blog/crypto"
	"corsa-blog/idl"
	"corsa-blog/web"
	"flag"
	"fmt"
	"log"
	"mime"
	"os"
)

func main() {
	var ver = flag.Bool("ver", false, "Prints the current version")
	var configfile = flag.String("config", "config.toml", "Configuration file path")
	var simulate = flag.Bool("simulate", false, "Simulate sending alarm")
	var initAccount = flag.Bool("initaccount", false, "Initialize the Admin account credentials")
	flag.Parse()

	if *ver {
		fmt.Printf("%s, version: %s", idl.Appname, idl.Buildnr)
		os.Exit(0)
	}

	if *initAccount {
		ac := crypto.NewUserCred()
		if err := ac.CreateAdminCredentials(); err != nil {
			log.Fatal("Error: ", err)
		}
		log.Println("Credential for admin successfully created. Please restart.")
		os.Exit(0)
	}

	if err := web.RunService(*configfile, *simulate); err != nil {
		panic(err)
	}
}

func init() {
	_ = mime.AddExtensionType(".js", "text/javascript")
	_ = mime.AddExtensionType(".css", "text/css")
	_ = mime.AddExtensionType(".mjs", "text/javascript")
	log.Printf("Init App %s with Mime override", idl.Appname)
}
