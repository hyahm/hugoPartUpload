package hugoPartUpload

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type PartClient struct {
	Token       string // 必填
	Identifier  string // 必填
	User        string // 必填
	Title       string
	Rule        string   // 必填
	Cat         string   // 必填
	Subcat      []string // 必填
	Actor       string
	Domain      string
	Filename    string // 必填
	UploadId    int
	NewFilename string
}

func (pc *PartClient) checkFiled() error {
	if pc.Domain == "" {
		pc.Domain = "http://admin.hugocut.com/"
	}

	if pc.Filename == "" {
		return errors.New("Filename not be empty")
	}
	if pc.User == "" {
		return errors.New("User not be empty")
	}

	if pc.Identifier == "" {
		return errors.New("Identifier not be empty")
	}

	if pc.Token == "" {
		return errors.New("Token not be empty")
	}

	if pc.Title == "" {
		pc.Title = pc.Rule
	}

	if pc.Rule == "" {
		return errors.New("Rule not be empty")
	}

	if pc.Cat == "" {
		return errors.New("Cat not be empty")
	}
	if pc.Domain[len(pc.Domain)-1:] == "/" {
		pc.Domain = pc.Domain[:len(pc.Domain)-1]
	}
	if pc.NewFilename == "" {
		i := strings.LastIndex(pc.Filename, ".")
		pc.NewFilename = pc.Identifier + pc.Filename[i:]
	}
	return nil
}

func (pc *PartClient) PartUpload() error {
	err := pc.checkFiled()
	if err != nil {
		return err
	}
	err = pc.initfunc()
	if err != nil {
		return err
	}
	return pc.dataForm()
}

type Data struct {
	UploadId int `json:"uploadId"`
}

type InitData struct {
	Code    int    `json:"code"`
	Data    Data   `json:"data"`
	Message string `json:"message"`
}

var PARTSIZE int64 = 1024 * 1024 * 10

func (pc *PartClient) initfunc() error {

	x := `
	{
		"fileName": "%s",
		"totalParts": %d,
		"totalSize": %d,
		"user": "%s"
	}
	`
	f, err := os.Open(pc.Filename)
	if err != nil {
		return err
	}
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	tp := fi.Size() / PARTSIZE
	if fi.Size()%PARTSIZE != 0 {
		tp++
	}
	x = fmt.Sprintf(x, pc.NewFilename, tp, fi.Size(), pc.User)
	cli := &http.Client{}

	r, err := http.NewRequest("POST", pc.Domain+"/audio.php/VideoUpload/initiateMultipartUpload", strings.NewReader(x))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Token", pc.Token)
	resp, err := cli.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	init := &InitData{}

	err = json.NewDecoder(resp.Body).Decode(init)
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", *init)
	pc.UploadId = init.Data.UploadId
	return nil
}

func (pc *PartClient) dataForm() error {
	b, err := ioutil.ReadFile(pc.Filename)
	if err != nil {
		return err
	}

	i := 0
	l := len(b)
	cli := http.Client{}
	for int64(i)*PARTSIZE < int64(l) {
		buf := new(bytes.Buffer)
		w := multipart.NewWriter(buf)
		w.WriteField("uploadId", fmt.Sprintf("%d", pc.UploadId))
		w.WriteField("user", pc.User)
		part, err := w.CreateFormFile("file", fmt.Sprintf("%d.mp4", i))
		if err != nil {
			return err
		}
		if int64(i+1)*PARTSIZE > int64(l) {
			_, err = io.Copy(part, bytes.NewReader(b[int64(i)*PARTSIZE:]))
		} else {
			_, err = io.Copy(part, bytes.NewReader(b[int64(i)*PARTSIZE:(int64(i)+1)*PARTSIZE]))
		}
		w.WriteField("partNumber", fmt.Sprintf("%d", i+1))
		req, err := http.NewRequest("POST", pc.Domain+"/audio.php/VideoUpload/uploadPart", buf)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
		req.Header.Set("Token", pc.Token)
		resp, err := cli.Do(req)
		if err != nil {
			return err

		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err

		}
		fmt.Println(string(b))
		i++
		// return
	}
	return pc.complate()
}

func (pc *PartClient) complate() error {
	cli := http.Client{}
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	w.WriteField("uploadId", fmt.Sprintf("%d", pc.UploadId))
	w.WriteField("user", pc.User)
	w.WriteField("identifier", pc.Identifier)
	w.WriteField("title", pc.Title)
	w.WriteField("rule", pc.Rule)
	w.WriteField("cat", pc.Cat)
	w.WriteField("subcat", strings.Join(pc.Subcat, ","))
	w.WriteField("actor", pc.Actor)

	req, err := http.NewRequest("POST", pc.Domain+"/audio.php/VideoUpload/completeMultipartUpload", buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Token", pc.Token)
	resp, err := cli.Do(req)
	if err != nil {
		return err

	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err

	}
	fmt.Println(string(b))
	return nil
}
