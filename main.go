package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"spatula/downloader"
	"syscall"
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
	t := downloader.NewTask(*link, *path, *filename, *useragent)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		sig := <-sigChannel
		cancel()
		fmt.Println(sig)
	}()
	t.Download(ctx, *blockSize, *maxGoRoutines)
}
