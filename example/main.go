package main

import (
	"github.com/hyahm/golog"
	"github.com/hyahm/hugoPartUpload"
)

func main() {
	defer golog.Sync()
	pc := hugoPartUpload.PartClient{
		Filename:    "C:\\Users\\Admin\\Desktop\\dongman\\1.mp4",
		Token:       "xxxxxx",
		Identifier:  "aaaa7",
		User:        "ceshi2",
		Title:       "test",
		Rule:        "test",
		Cat:         "mm_在线电影",
		Subcat:      []string{"222"},
		Actor:       "test",
		Domain:      "",
		NewFilename: "bbb.mp4",
	}
	err := pc.PartUpload()
	if err != nil {
		golog.Fatal(err)
	}
}
