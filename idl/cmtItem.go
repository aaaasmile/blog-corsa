package idl

import (
	"corsa-blog/conf"
	"fmt"
	"log"
	"strings"
)

func (ci *CmtItem) GetLocationFromPostId() string {
	if ll, ok := GetStaticLocation(ci.PostId); ok {
		return ll
	}
	return "/"
}

func GetStaticLocation(id string) (res string, ok bool) {
	ok = false
	arr := strings.Split(id, "-")
	if len(arr) < 2 {
		log.Println("[GetStaticLocation] WARN id not recognized", id)
		return
	}
	last := arr[len(arr)-1]
	switch last {
	case "PS":
		// Example of result "/posts/2024/11/08/24-11-08-ProssimaGara/" from PostId 24-11-08-ProssimaGara-PS
		if len(arr) < 5 {
			log.Println("[GetStaticLocation] WARN id for Post not recognized", id)
			return
		}
		yy := arr[0]
		mm := arr[1]
		dd := arr[2]
		title_ix := len(arr) - 1
		title := strings.Join(arr[0:title_ix], "-")
		res = fmt.Sprintf("/%s/20%s/%s/%s/%s/",
			conf.Current.PostSubDir,
			yy,
			mm,
			dd,
			title,
		)
		ok = true
	case "PG":
		if len(arr) != 2 {
			log.Println("[GetStaticLocation] WARN id for Page not recognized", id)
			return
		}
		title := arr[0]
		res = fmt.Sprintf("/%s/%s", conf.Current.PageSubDir, title)
		ok = true
	}
	return
}
