# hugo 分片上传包


```go
func main() {
	pc := PartClient{
		Filename:    "C:\\Users\\Admin\\Desktop\\a.mp4",
		Token:       "xxxxxxxxxxxx",
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
```