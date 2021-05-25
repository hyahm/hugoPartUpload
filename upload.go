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

	"github.com/hyahm/golog"
)

type PartClient struct {
	Token       string // 必填
	Identifier  string // 必填
	User        string // 必填
	Title       string
	Audio       string
	Rule        string   // 必填
	Cat         string   // 必填
	Subcat      []string // 必填
	Actor       string
	Domain      string
	Filename    string // 必填
	Cover       string // 封面图
	UploadId    int
	NewFilename string
}

func (pc *PartClient) checkFiled() error {
	if pc.Domain == "" {
		pc.Domain = "http://admin.hugocut.com"
	}

	if pc.Filename == "" {
		return errors.New("filename not be empty")
	}
	if pc.User == "" {
		return errors.New("user not be empty")
	}

	if pc.Identifier == "" {
		return errors.New("identifier not be empty")
	}

	if pc.Token == "" {
		return errors.New("token not be empty")
	}

	if pc.Title == "" {
		pc.Title = pc.Rule
	}

	if pc.Rule == "" {
		return errors.New("rule not be empty")
	}

	if pc.Cat == "" {
		return errors.New("cat not be empty")
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

func (pc *PartClient) Upload() error {
	err := pc.checkFiled()
	if err != nil {
		return err
	}
	return pc.upload()
}

type Data struct {
	UploadId int `json:"uploadId"`
}

type InitData struct {
	Code    int    `json:"code"`
	Data    Data   `json:"data"`
	Message string `json:"message"`
}

var PARTSIZE int64 = 1024 * 1024 * 10 // 10M

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
		golog.Error(err)
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
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Token", pc.Token)
	resp, err := cli.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	init := &InitData{}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	golog.Info(string(b))
	err = json.Unmarshal(b, init)
	if err != nil {
		return err
	}
	pc.UploadId = init.Data.UploadId
	return nil
}

func (pc *PartClient) dataForm() error {
	f, err := os.Open(pc.Filename)
	// b, err := ioutil.ReadFile(pc.Filename)
	if err != nil {
		return err
	}
	var i int64 = 0
	// l := len(b)

	for {

		buf := new(bytes.Buffer)
		w := multipart.NewWriter(buf)
		w.WriteField("uploadId", fmt.Sprintf("%d", pc.UploadId))
		w.WriteField("user", pc.User)
		part, err := w.CreateFormFile("file", fmt.Sprintf("%d.mp4", i))
		if err != nil {
			return err
		}
		_, err = f.Seek(i*PARTSIZE, 0)
		if err != nil {
			return err
		}
		b := make([]byte, PARTSIZE)
		n, err := f.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		golog.Info(n)
		_, err = io.Copy(part, bytes.NewReader(b[:n]))
		if err != nil {
			golog.Info(err)
			return err
		}
		// if int64(i+1)*PARTSIZE > int64(l) {
		// 	_, err = io.Copy(part, bytes.NewReader(b[int64(i)*PARTSIZE:]))
		// 	if err != nil {
		// 		return err
		// 	}
		// } else {
		// 	_, err = io.Copy(part, bytes.NewReader(b[int64(i)*PARTSIZE:(int64(i)+1)*PARTSIZE]))
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		w.WriteField("partNumber", fmt.Sprintf("%d", i+1))
		req, err := http.NewRequest("POST", pc.Domain+"/audio.php/VideoUpload/uploadPart", buf)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
		req.Header.Set("Token", pc.Token)
		cli := http.Client{}
		resp, err := cli.Do(req)
		if err != nil {
			return err

		}
		rb, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		golog.Info(string(rb))
		i++
		// return
	}
	return pc.complate()
}

func (pc *PartClient) complate() error {
	cli := http.Client{}
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	if pc.Cover != "" {
		imageb, err := ioutil.ReadFile(pc.Cover)
		if err != nil {
			golog.Error(err)
			return err
		}
		image, err := w.CreateFormFile("image", fmt.Sprintf("%s.jpg", pc.Identifier))
		if err != nil {
			golog.Error(err)
			return err
		}
		_, err = io.Copy(image, bytes.NewReader(imageb))
		if err != nil {
			return err
		}
	}

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
	golog.Info(string(b))
	return nil
}

func (pc *PartClient) upload() error {
	if pc.Token == "" {
		return errors.New("token not be empty")
	}
	if pc.Audio == "" {
		return errors.New("audio not be empty")
	}
	cli := http.Client{}
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	videob, err := ioutil.ReadFile(pc.Audio)
	if err != nil {
		return err
	}
	audio, err := w.CreateFormFile("audio", pc.Filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(audio, bytes.NewReader(videob))
	if err != nil {
		return err
	}
	imageb, err := ioutil.ReadFile(pc.Cover)
	if err != nil {
		return err
	}
	image, err := w.CreateFormFile("image", fmt.Sprintf("%s.jpg", pc.Identifier))
	if err != nil {
		return err
	}

	_, err = io.Copy(image, bytes.NewReader(imageb))
	if err != nil {
		return err
	}
	w.WriteField("uploadId", fmt.Sprintf("%d", pc.UploadId))
	w.WriteField("user", pc.User)
	w.WriteField("identifier", pc.Identifier)
	w.WriteField("title", pc.Title)
	w.WriteField("rule", pc.Rule)
	w.WriteField("cat", pc.Cat)
	w.WriteField("subcat", strings.Join(pc.Subcat, ","))
	w.WriteField("actor", pc.Actor)
	w.WriteField("filename", pc.Filename)

	req, err := http.NewRequest("POST", pc.Domain+"/audio.php/VideoUpload/index", buf)
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
