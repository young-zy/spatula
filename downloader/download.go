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
	"unsafe"
)

var (
	wg sync.WaitGroup
	cookies []*http.Cookie
	client *http.Client
)

type panShare struct {
	Code string `json:"code"`
	Data data   `json:"datas"`
}

type data struct{
	Downlink []downlink `json:"downlink"`
}

type downlink struct {
	Link string 	`json:"link"`
	Name string		`json:"name"`
	Size string		`json:"size"`
	Time string		`json:"time"`
}

func init(){
	client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns: 2000,
			MaxIdleConnsPerHost: 2000,
		},
	}
}

func Resolve(link string) []downlink {
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
	body = *(*string)(unsafe.Pointer(&bodyBytes))
	var p panShare
	err = json.Unmarshal(bodyBytes, &p)
	if err != nil || p.Code != "200" {
		panic("failed to retrieve direct link")
	}
	reg = regexp.MustCompile("href=\"(?P<link>[\\S]*)\">")
	for i, d := range p.Data.Downlink {
		p.Data.Downlink[i].Link = reg.FindStringSubmatch(d.Link)[1]
	}
	return p.Data.Downlink
}

func Download(url string, blockSize int64, maxGoroutines int, path string, userFilename string)  {
	url, err := handle302(url)
	if err != nil {
		println(err)
	}
	wait := make(chan struct{}, maxGoroutines)
	for i:=0; i < maxGoroutines; i++ {
		wait <- struct{}{}
	}
	size, filename := getSizeAndName(url)
	if userFilename != ""{
		filename = userFilename
	}
	filename += ".tmp"
	fp, err := os.OpenFile(path+filename, os.O_RDWR|os.O_CREATE, 0100644)
	if err != nil {
		fmt.Println(err)
	}
	_, err = fp.Seek(size, 0)
	if err != nil {
		fmt.Println(err)
	}
	_, err = fp.Write([]byte{0})
	if err != nil {
		fmt.Println(err)
	}
	err = fp.Close()
	if err != nil {
		fmt.Println(err)
	}
	for i := int64(0); i<=size; i+=blockSize {
		offset := i
		wg.Add(1)
		<- wait
		go func() {
			count := 0
			for count < 20{
				err := downloadBlock(url, offset, blockSize, path, filename)
				if err != nil{
					fmt.Println(err)
					count++
				}else{
					break
				}
			}
			if count >= 20{
				panic(fmt.Sprintf("error while downloading %v-%v", offset, offset+blockSize))
			}
			wg.Done()
			wait <- struct{}{}
		}()
	}
	wg.Wait()
	filename = filename[:len(filename)-4]
	err = os.Rename(path+filename+".tmp", path+filename)
	if err != nil {
		log.Println(err)
		panic("failed to rename file after completion")
	}
}

func downloadBlock(url string, offset int64, size int64, path string, filename string) error{
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	downloadRange := fmt.Sprintf("bytes=%v-%v", offset, offset+size)
	req.Header.Add("Range", downloadRange)
	for _, c := range cookies {
		if c.Name == "BAIDUID" {
			req.AddCookie(c)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(resp.Body)
	fp, err := os.OpenFile(path+filename, os.O_RDWR|os.O_CREATE, 0100644)
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

func handle302(url string) (location string,err error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 302{
		location = resp.Header.Get("Location")
		cookies = resp.Cookies()
	} else {
		location = url
	}
	return location, err
}

func getSizeAndName(url string) (int64, string){
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
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
	return size, filename
}
