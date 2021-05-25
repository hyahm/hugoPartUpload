package main

import (
	"github.com/hyahm/golog"
	"github.com/hyahm/hugoPartUpload"
)

func main() {
	defer golog.Sync()
	pc := hugoPartUpload.PartClient{
		Filename:    "C:\\Users\\Admin\\Desktop\\dongman\\1.mp4",
		Token:       "xxxxxxxx",
		Identifier:  "aaaa",
		User:        "ceshi2",
		Title:       "test",
		Rule:        "test",
		Cat:         "mm_手机下载",
		Subcat:      []string{"动漫"},
		Actor:       "test",
		NewFilename: "bbb.mp4",
	}
	err := pc.PartUpload()
	if err != nil {
		golog.Fatal(err)
	}
}
