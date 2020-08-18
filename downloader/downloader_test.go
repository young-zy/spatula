package downloader

import (
	"fmt"
	"testing"
)

func TestDownloader(t *testing.T) {

}

func TestDownload(t *testing.T) {

}

func TestResolve(t *testing.T) {
	link := "https://pan.baidu.com/s/1mgpAh76"
	fmt.Printf("%v", Resolve(link))
}
