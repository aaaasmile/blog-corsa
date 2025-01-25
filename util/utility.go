package util

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

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

func FormatDateIt(tt time.Time) string {
	res := fmt.Sprintf("%d %s %d", tt.Day(), MonthToStringIT(tt.Month()), tt.Year())
	return res
}

func MonthToStringIT(month time.Month) string {
	switch month {
	case time.January:
		return "Gennaio"
	case time.February:
		return "Febbraio"
	case time.March:
		return "Marzo"
	case time.April:
		return "Aprile"
	case time.May:
		return "Maggio"
	case time.June:
		return "Giugno"
	case time.July:
		return "Luglio"
	case time.August:
		return "Agosto"
	case time.September:
		return "Settembre"
	case time.October:
		return "Ottobre"
	case time.November:
		return "Novembre"
	case time.December:
		return "Dicembre"
	default:
		return ""
	}
}

func PseudoUuid() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, nil
}
