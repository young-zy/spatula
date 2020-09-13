package downloader

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestDownloader(t *testing.T) {

}

func TestDownload(t *testing.T) {
	task := NewTask(
		"https://oxygenos.oneplus.net/OnePlus6Oxygen_22_OTA_047_all_2007191515_bd6f7476887846cb.zip",
		"",
		"",
		"",
	)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChannel
		cancel()
		fmt.Println(sig)
	}()
	task.Download(ctx, 102400000, 1)
}
