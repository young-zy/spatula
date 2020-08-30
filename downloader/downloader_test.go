package downloader

import (
	"testing"
)

func TestDownloader(t *testing.T) {

}

func TestDownload(t *testing.T) {
	task := NewTask(
		"https://oxygenos.oneplus.net/OnePlus6Oxygen_22_OTA_047_all_2007191515_bd6f7476887846cb.zip",
		"",
		"OnePlus6Oxygen_22_OTA_047_all_2007191515_bd6f7476887846cb.zip",
		"",
	)
	task.Download(102400, 500)
}
