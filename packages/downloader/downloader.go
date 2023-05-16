package downloader

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
	pb "wget/packages/progress-bar"
)

type downloader struct {
	path        string
	rateLimit   int
	progressBar bool
	Result      chan downloadResult
	WG          *sync.WaitGroup
}

type downloadResult struct {
	Size int64
	File string
	Err  error
}

func CreateDownloader(path string, rateLimit int, progressBar bool) *downloader {
	ch := make(chan downloadResult)

	var wg sync.WaitGroup
	return &downloader{
		path:        path,
		rateLimit:   rateLimit,
		progressBar: progressBar,
		Result:      ch,
		WG:          &wg,
	}
}

func (d *downloader) DownloadFile(url, fileName string) {
	defer d.WG.Done()

	if fileName == "" {
		fileName = path.Base(url)
	}

	resp, err := http.Get(url)
	if err != nil {
		d.Result <- downloadResult{
			Err: err,
		}
		return
	}

	if resp.StatusCode == http.StatusOK {
		log.Printf("Downloading file: %s. Status 200 OK\n", fileName)
	} else {
		log.Printf("Error downloading %s: recieved status code %d\n", url, resp.StatusCode)
	}

	file, err := os.Create(d.path + fileName)
	if err != nil {
		d.Result <- downloadResult{
			Err: err,
		}
		return
	}

	defer file.Close()

	currentSize := int64(0)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	if resp.ContentLength > 0 && d.progressBar {
		go oneSecondTick(&currentSize, resp, ticker)
	}

	limit := pb.NewLimitedReader(resp.Body, d.rateLimit, &currentSize)
	buffer := make([]byte, 32768)
	if d.rateLimit > 0 {
		buffer = make([]byte, d.rateLimit)
	}

	for {
		n, err := limit.Read(buffer)
		if err != nil && err != io.EOF {
			d.Result <- downloadResult{
				Err: err,
			}
			return
		}
		if n == 0 {
			break
		}

		if _, err := file.Write(buffer[:n]); err != nil {
			d.Result <- downloadResult{
				Err: err,
			}
			return
		}
	}

	if resp.ContentLength > 0 && d.progressBar {
		log.Printf("\r[%s] %.2f%% of %d bytes\n", pb.ProgressBar(currentSize, resp.ContentLength), float64(currentSize)/float64(resp.ContentLength)*100, resp.ContentLength)
	}

	d.Result <- downloadResult{
		Size: currentSize,
		File: d.path + fileName,
		Err:  nil,
	}
}

func oneSecondTick(size *int64, r *http.Response, ticker *time.Ticker) {
	for range ticker.C {
		log.Printf("\r[%s] %.2f%% of %d bytes", pb.ProgressBar(*size, r.ContentLength), float64(*size)/float64(r.ContentLength)*100, r.ContentLength)
	}
}
