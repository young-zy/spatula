package downloader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

var (
	wg     sync.WaitGroup
	client *http.Client
)

type panShare struct {
	Code string `json:"code"`
	Data *data  `json:"datas"`
}

type data struct {
	Downlink []*downlink `json:"downlink"`
}

type downlink struct {
	Link string `json:"link"`
	Name string `json:"name"`
	Size string `json:"size"`
	Time string `json:"time"`
}

type task struct {
	url          string
	path         string
	filename     string
	tempFilename string
	cookies      []*http.Cookie
	size         int64
	useragent    string
}

func init() {
	client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        2000,
			MaxIdleConnsPerHost: 2000,
		},
		Timeout: 30 * time.Second,
	}
}

func NewTask(url string, path string, filename string, userAgent string) *task {
	return &task{
		url:       url,
		path:      path,
		filename:  filename,
		useragent: userAgent,
	}
}

func Resolve(link string) []*downlink {
	resp, err := http.Get("https://pan.naifei.cc/new/")
	if err != nil || resp.StatusCode != 200 {
		log.Println(err)
		panic("failed to connect to resolver")
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		panic("failed to retrieve sign from resolver")
	}
	body := *(*string)(unsafe.Pointer(&bodyBytes))
	//body := string(bodyBytes)
	reg := regexp.MustCompile("articleFrom\\['sign'] = \"(?P<sign>[\\S]*)\"")
	sign := reg.FindStringSubmatch(body)[1]
	resp, err = http.PostForm(
		"https://pan.naifei.cc/new/panshare.php",
		url.Values{"sign": {sign}, "link": {link}},
	)
	if err != nil || resp.StatusCode != 200 {
		log.Println(err)
		panic("resolver error")
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		panic("failed to resolve result from resolver")
	}
	var p panShare
	err = json.Unmarshal(bodyBytes, &p)
	if err != nil || p.Code != "200" {
		panic("failed to retrieve direct link")
	}
	reg = regexp.MustCompile("href=\"(?P<link>[\\S]*)\">")
	for _, d := range p.Data.Downlink {
		d.Link = reg.FindStringSubmatch(d.Link)[1]
	}
	return p.Data.Downlink
}

func (t *task) Download(blockSize int64, maxGoroutines int) {
	err := t.handle302()
	if err != nil {
		panic(err)
	}
	wait := make(chan struct{}, maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		wait <- struct{}{}
	}
	t.getSizeAndName()
	t.tempFilename = t.filename + ".tmp"
	t.createTempFile()
	for i := int64(0); i <= t.size; i += blockSize {
		offset := i
		wg.Add(1)
		<-wait
		go func() {
			count := 0
			for count < 20 {
				if t.size-offset < blockSize {
					blockSize = t.size - offset
				}
				err := t.downloadBlock(offset, blockSize)
				if err != nil {
					fmt.Println(err)
					count++
				} else {
					break
				}
			}
			if count >= 20 {
				panic(fmt.Sprintf("error while downloading %v-%v", offset, offset+blockSize))
			}
			wg.Done()
			wait <- struct{}{}
		}()
	}
	wg.Wait()
	err = os.Rename(t.path+t.tempFilename, t.path+t.filename)
	if err != nil {
		log.Println(err)
		panic("failed to rename file after completion")
	}
}

func (t *task) downloadBlock(offset int64, size int64) error {
	req, err := http.NewRequest("GET", t.url, nil)
	if err != nil {
		return err
	}
	downloadRange := fmt.Sprintf("bytes=%v-%v", offset, offset+size)
	req.Header.Add("Range", downloadRange)
	if t.useragent != "" {
		req.Header.Set("User-Agent", t.useragent)
	}
	for _, c := range t.cookies {
		if c.Name == "BAIDUID" {
			req.AddCookie(c)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(resp.Body)
	fp, err := os.OpenFile(t.path+t.tempFilename, os.O_RDWR, 0100644)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.Seek(offset, 0)
	_, err = reader.WriteTo(fp)
	if err != nil {
		return err
	}
	fmt.Printf("finished %v-%v\n", offset, offset+size)
	return nil
}

func (t *task) handle302() error {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(t.url)
	if err != nil {
		return err
	}
	if resp.StatusCode == 302 {
		t.url = resp.Header.Get("Location")
		t.cookies = resp.Cookies()
	}
	return nil
}

func (t *task) createTempFile() {
	_, err := os.Stat(t.path + t.filename) //check file exists
	if err == nil || os.IsExist(err) {
		panic("file already exists")
	}
	_, err = os.Stat(t.path + t.tempFilename) //check file exists
	if err == nil || os.IsExist(err) {
		panic("temp file already exists")
	}
	fp, err := os.OpenFile(t.path+t.tempFilename, os.O_RDWR|os.O_CREATE, 0100644)
	if err != nil {
		fmt.Println(err)
		panic("failed to open file")
	}
	defer fp.Close()
	_, err = fp.Seek(t.size, 0)
	if err != nil {
		fmt.Println(err)
		panic("failed to move to file pointer to file end")
	}
	_, err = fp.Write([]byte{0})
	if err != nil {
		fmt.Println(err)
		panic("failed to write zero to file end")
	}
}

func (t *task) getSizeAndName() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", t.url, nil)
	if err != nil {
		panic("error while getting size")
	}
	req.Header.Add("Range", "bytes=0-0")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		panic("error while trying to get file size")
	}
	size, err := strconv.ParseInt(resp.Header.Get("x-bs-file-size"), 10, 64)
	if err != nil {
		log.Println(err)
		panic("failed to transform result into int64")
	}
	reg := regexp.MustCompile("attachment;filename=\"(?P<filename>[\\S, ]*)\"")
	filename := reg.FindStringSubmatch(resp.Header.Get("Content-Disposition"))[1]
	if t.filename == "" {
		t.filename = filename
	}
	t.size = size
}
