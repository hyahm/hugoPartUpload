package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

func main() {
	postData := make(map[string]string)
	postData["anchorId"] = "361155076095561728"
	postData["searchBegin"] = "2019-03-01 00:00:00"
	postData["searchEnd"] = "2020-03-10 00:00:00"
	url := "http://127.0.0.1:8888/upload"
	PostWithFormData("POST", url, &postData)
}

func PostWithFormData(method, url string, postData *map[string]string) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k, v := range *postData {
		w.WriteField(k, v)
	}
	w.Close()
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, _ := http.DefaultClient.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(resp.StatusCode)
	fmt.Printf("%s", data)
}
