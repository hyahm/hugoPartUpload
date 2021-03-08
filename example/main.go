package main

import (
	"log"

	"github.com/hyahm/hugoPartUpload"
)

func main() {
	pc := hugoPartUpload.PartClient{
		Filename:    "C:\\Users\\Admin\\Desktop\\a.mp4",
		Token:       "d26fcd4a11538c54071ad0c803034f0dea737a82",
		Identifier:  "aaaa",
		User:        "test",
		Title:       "test",
		Rule:        "test",
		Cat:         "mm_手机下载",
		Subcat:      []string{"经典"},
		Actor:       "test",
		Domain:      "http://192.168.50.72",
		NewFilename: "bbb.mp4",
	}
	err := pc.PartUpload()
	if err != nil {
		log.Fatal(err)
	}
}
