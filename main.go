package main

import (
	"flag"
	"spatula/downloader"
)

var (
	link          = flag.String("l", "", "the download link, must not be empty")
	useragent     = flag.String("u", "", "customized useragent")
	path          = flag.String("p", "", "output path")
	filename      = flag.String("o", "", "out put filename")
	blockSize     = flag.Int64("c", 4194304, "size of each file block, if download speed is limited, set a value below the limit")
	maxGoRoutines = flag.Int("g", 20, "max goroutines opened, value over 1000 is not recommended")
)

func main() {
	flag.Parse()
	//reg := regexp.MustCompile(`(链接)?([:, ：]?)( *)(https://pan\.baidu\.com)?/?s?/?(?P<hash>[\S]*)( *)(提取码)?([:, ：]?)( *)(?P<code>[\S]{4})( *)(复制这段内容后打开百度网盘手机App，操作更方便哦)?`)
	//matches := reg.FindStringSubmatch(*link)
	//hash := matches[5]
	//code := matches[10]
	//*link = fmt.Sprintf("链接: https://pan.baidu.com/s/%v 提取码: %v", hash, code)
	//links := downloader.Resolve(*link)
	//index := 0
	//if len(links) < 1 {
	//	fmt.Println("resolve failed, no file found (Note that folders are not supported) ")
	//}
	//if len(links) > 1 {
	//	fmt.Println("Multiple files detected, pls select(folders are not supported): ")
	//	for i, o := range links {
	//		fmt.Printf("%4d\t%10v\n", i+1, o.Name)
	//	}
	//	fmt.Print("Please Select: ")
	//	_, err := fmt.Scanln(&index)
	//	if err != nil {
	//		log.Println(err)
	//		panic("error when input")
	//	}
	//	index--
	//}
	//if index > len(links) {
	//	panic("index out of range")
	//}
	//*link = links[index].Link
	//println("If being limited, value 10240 is suggested")
	//print("Please input file block size[default: 4194304]: ")
	//var blockSize int64
	//var maxGoRoutines int
	//_, err := fmt.Scanln(&blockSize)
	//if err != nil {
	//	blockSize = 4194304
	//}
	//println("Value beyond 1000 is not suggested")
	//print("Please input file block size[default: 20]: ")
	//_, err = fmt.Scanln(&maxGoRoutines)
	//if err != nil {
	//	maxGoRoutines = 20
	//}
	t := downloader.NewTask(*link, *path, *filename, *useragent)
	t.Download(*blockSize, *maxGoRoutines)
}
