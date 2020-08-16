package downloader

import (
	"fmt"
	"testing"
)

func TestDownloader(t *testing.T)  {
	url := ""
	offset := int64(10240)
	size := int64(1024)
	err := downloadBlock(url, offset, size, "", "test.tmp")
	if err != nil {
		t.Error(err)
	}
}

func TestDownload(t *testing.T) {
	url := ""
	blockSize := int64(10240)
	maxGoroutines := 1000
	Download(url, blockSize, maxGoroutines, "", "test")
}

func TestResolve(t *testing.T) {
	link := "https://pan.baidu.com/s/1mgpAh76"
	fmt.Printf("%v", Resolve(link))
}