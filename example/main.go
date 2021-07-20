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
		Token:       "9252b3ee5f7f9a078a6f27a3",
		Identifier:  "9252b3ee5f68a6f27a3",
		User:        "uu_upload",
		Title:       "0601测试2333",
		Rule:        "2021_API",
		Cat:         "uu手动上传",
		Subcat:      []string{"手动"},
		Actor:       "test",
		NewFilename: "7d505b2ddf5tt39cc410847b7ab5018b8e1.mp4",
		InitUrl:     "/init",
		UploadUrl:   "/upload",
		Domain:      "http://127.0.0.1:8888",
	}
	err := pc.PartUpload()
	if err != nil {
		golog.Fatal(err)
	}
	fmt.Println(time.Since(start).Seconds())
}
