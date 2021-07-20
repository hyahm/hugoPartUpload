package main

import (
	"fmt"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/hugoPartUpload"
)

func main() {
	start := time.Now()
	defer golog.Sync()
	pc := hugoPartUpload.PartClient{
		Filename:    "C:\\Users\\Admin\\Desktop\\dongman\\1.mp4",
		Token:       "9252b3ee5f7f9a07f652d4e75faa9a268a6f27a3",
		Identifier:  "9252b3ee5f7f9a07f652d4e75faa9a268a6f27a3",
		User:        "admin",
		Title:       "0601测试2333",
		Rule:        "testadmin",
		Cat:         "admin",
		Subcat:      []string{"test"},
		Actor:       "test",
		NewFilename: "7d505b2ddf5tt39cc410847b7ab5018b8e1.mp4",
		InitUrl:     "/init",
		UploadUrl:   "/upload",
		CompleteUrl: "/complete",
		Domain:      "http://127.0.0.1:8888",
	}
	err := pc.PartUpload()
	if err != nil {
		golog.Fatal(err)
	}
	fmt.Println(time.Since(start).Seconds())
}
