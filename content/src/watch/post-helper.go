package watch

import (
	"corsa-blog/db"
	"corsa-blog/idl"
	"log"
)

func CreateMapLinks(liteDB *db.LiteDB) (*idl.MapPostsLinks, error) {
	mapLinks := &idl.MapPostsLinks{
		MapPost:  map[string]idl.PostLinks{},
		ListPost: []idl.PostItem{},
	}
	var err error
	mapLinks.ListPost, err = liteDB.GetPostList()
	if err != nil {
		return nil, err
	}
	//fmt.Println("*** Posts ", mapLinks.ListPost)
	last_ix := len(mapLinks.ListPost) - 1
	prev_item := &idl.PostItem{}
	next_item := &idl.PostItem{}

	for ix, item := range mapLinks.ListPost {
		postLinks := idl.PostLinks{
			Item: &item,
		}
		if last_ix > 0 {
			// at least 2 or more elements
			if ix == 0 {
				next_item = &mapLinks.ListPost[ix+1]
				postLinks.NextLink = next_item.Uri
				postLinks.NextPostID = next_item.PostId
			} else if ix == last_ix {
				postLinks.PrevLink = prev_item.Uri
				postLinks.PrevPostID = prev_item.PostId
			} else {
				next_item = &mapLinks.ListPost[ix+1]
				postLinks.NextLink = next_item.Uri
				postLinks.NextPostID = next_item.PostId
				postLinks.PrevLink = prev_item.Uri
				postLinks.PrevPostID = prev_item.PostId
			}
			prev_item = &mapLinks.ListPost[ix]
		}
		mapLinks.MapPost[item.PostId] = postLinks
	}
	//fmt.Println("*** map ", mapLinks.MapPost)
	log.Printf("Built map with %d items", len(mapLinks.MapPost))
	return mapLinks, nil
}
