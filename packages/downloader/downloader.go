package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
	pb "wget/packages/progress-bar"
)

type DownloadResult struct {
	Size int64
	File string
	Err  error
}

func Download(url, fileName, dirPath string, rateLimit int, ch chan DownloadResult, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		ch <- DownloadResult{
			Err: err,
		}
		return
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Downloading file: %s. Status 200 OK\n", fileName)
	} else {
		fmt.Printf("Error downloading %s: recieved status code %d\n", url, resp.StatusCode)
	}

	file, err := os.Create(dirPath + fileName)
	if err != nil {
		ch <- DownloadResult{
			Err: err,
		}
		return
	}

	defer file.Close()

	currentSize := int64(0)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	if resp.ContentLength > 0 {
		go oneSecondTick(&currentSize, resp, ticker)
	}

	limit := pb.NewLimitedReader(resp.Body, rateLimit, &currentSize)
	buffer := make([]byte, 32768)
	if rateLimit > 0 {
		buffer = make([]byte, rateLimit)
	}

	for {
		n, err := limit.Read(buffer)
		if err != nil && err != io.EOF {
			ch <- DownloadResult{
				Err: err,
			}
			return
		}
		if n == 0 {
			break
		}

		// Write the chunk to disk
		if _, err := file.Write(buffer[:n]); err != nil {
			ch <- DownloadResult{
				Err: err,
			}
			return
		}
	}

	if resp.ContentLength > 0 {
		fmt.Printf("\r[%s] %.2f%% of %d bytes\n", pb.ProgressBar(currentSize, resp.ContentLength), float64(currentSize)/float64(resp.ContentLength)*100, resp.ContentLength)
	}

	ch <- DownloadResult{
		Size: currentSize,
		File: dirPath + fileName,
		Err:  nil,
	}
}

func oneSecondTick(size *int64, r *http.Response, ticker *time.Ticker) {
	for range ticker.C {
		fmt.Printf("\r[%s] %.2f%% of %d bytes", pb.ProgressBar(*size, r.ContentLength), float64(*size)/float64(r.ContentLength)*100, r.ContentLength)
	}
}
