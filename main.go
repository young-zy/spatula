package main

import (
	"fmt"
	"log"
	"os"
	"spatula/downloader"
)

func main() {
	link := os.Args[len(os.Args)-1]
	links := downloader.Resolve(link)
	index := 0
	if len(links) > 1 {
		fmt.Println("Multiple files detected, pls select(folders are not supported): ")
		for i, o := range links {
			fmt.Printf("%4d\t%10v\n", i+1, o.Name)
		}
		fmt.Print("Please Select: ")
		_, err := fmt.Scanln(&index)
		if err != nil {
			log.Println(err)
			panic("error when input")
		}
		index--
	}
	if index > len(links) {
		panic("index out of range")
	}
	link = links[index].Link
	println("If being limited, value 10240 is suggested")
	print("Please input file block size[default: 4194304]: ")
	var blockSize int64
	var maxGoRoutines int
	_, err := fmt.Scanln(&blockSize)
	if err != nil {
		blockSize = 4194304
	}
	println("Value beyond 1000 is not suggested")
	print("Please input file block size[default: 20]: ")
	_, err = fmt.Scanln(&maxGoRoutines)
	if err != nil {
		maxGoRoutines = 20
	}
	downloader.Download(link, blockSize, maxGoRoutines, "", "")
}
