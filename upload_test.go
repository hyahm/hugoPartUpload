package hugoPartUpload

import (
	"os"
	"testing"

	"github.com/hyahm/golog"
)

func TestSeek(t *testing.T) {
	defer golog.Sync()
	f, err := os.Open("aa.txt")
	if err != nil {
		golog.Error(err)
		return
	}
	ret, err := f.Seek(13, 0)
	if err != nil {
		return
	}
	golog.Info(ret)
	b := make([]byte, 10)
	n, err := f.Read(b)
	if err != nil {
		golog.Error(err)
	}
	golog.Info(n)
	golog.Info(string(b))
}
