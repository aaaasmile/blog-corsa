package util

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kardianos/osext"
)

var rootPath string
var UseRelativeRoot = true

func GetEncryptedPasswordFromPrompt() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter the secret to encode")
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	buf := []byte(response)
	return base64.StdEncoding.EncodeToString(buf)
}

func GetFullPath(relPath string) string {
	if UseRelativeRoot {
		return relPath
	}

	if rootPath == "" {
		var err error
		rootPath, err = osext.ExecutableFolder()
		if err != nil {
			log.Fatalf("ExecutableFolder failed: %v", err)
		}
		log.Println("Executable folder (rootdir) is ", rootPath)
		//rootPath, _ = filepath.Split(os.Args[0]) // os.Args[0] can be "faked". (https://github.com/kardianos/osext)
	}
	r := filepath.Join(rootPath, relPath)
	return r
}
