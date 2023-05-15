package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const progressBarWidth = 20

type DownloadResult struct {
	Size int64
	File string
	Err  error
}

func Download(url, fileName, dirPath string, ch chan DownloadResult, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		ch <- DownloadResult{
			Err: err,
		}
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Sending request, awaiting response... status 200 OK\n")
	} else {
		fmt.Printf("Error downloading %s: recieved status code %d\n", url, resp.StatusCode)
	}

	file, err := os.Create(dirPath + fileName)
	if err != nil {
		ch <- DownloadResult{
			Err: err,
		}
	}

	defer file.Close()

	currentSize := int64(0)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	if resp.ContentLength > 0 {
		go oneSecondTick(&currentSize, resp, ticker)
	}

	_, err = io.Copy(file, io.TeeReader(resp.Body, &progressReader{&currentSize}))
	if err != nil {
		ch <- DownloadResult{
			Err: err,
		}
	}

	if resp.ContentLength > 0 {
		fmt.Printf("\r[%s] %.2f%% of %d bytes\n", progressBar(currentSize, resp.ContentLength), float64(currentSize)/float64(resp.ContentLength)*100, resp.ContentLength)
	}

	ch <- DownloadResult{
		Size: currentSize,
		File: dirPath + fileName,
		Err:  nil,
	}
}

type progressReader struct {
	currentSize *int64
}

func (pr *progressReader) Write(p []byte) (n int, err error) {
	n = len(p)
	*pr.currentSize += int64(n)
	return
}

func progressBar(currentSize, totalSize int64) string {
	progress := int(float64(currentSize) / float64(totalSize) * float64(progressBarWidth))
	return strings.Repeat("=", progress) + strings.Repeat("-", progressBarWidth-progress)
}

func oneSecondTick(size *int64, r *http.Response, ticker *time.Ticker) {
	for range ticker.C {
		fmt.Printf("\r[%s] %.2f%% of %d bytes", progressBar(*size, r.ContentLength), float64(*size)/float64(r.ContentLength)*100, r.ContentLength)
	}
}
