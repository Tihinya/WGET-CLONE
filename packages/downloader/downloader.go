package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
	pb "wget/packages/progress-bar"
	"wget/packages/utils"
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

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		d.Result <- downloadResult{
			Err: err,
		}
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		d.Result <- downloadResult{
			Err: err,
		}
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Downloading file: %s. Status 200 OK\n", fileName)
	} else {
		d.Result <- downloadResult{
			Err: fmt.Errorf("error downloading %s: recieved status code %d", url, resp.StatusCode),
		}
		return
	}

	file, err := os.Create(path.Join(d.path, fileName))
	if err != nil {
		d.Result <- downloadResult{
			Err: err,
		}
		return
	}

	defer file.Close()
	currentSize := int64(0)
	lastTime := time.Now()
	lastSize := int64(0)

	if resp.ContentLength > 0 && d.progressBar {
		ticker := time.NewTicker(time.Second)

		defer ticker.Stop()
		go oneSecondTick(&currentSize, resp, ticker, &lastTime, &lastSize)
	}

	limit := pb.NewLimitedReader(resp.Body, d.rateLimit, &currentSize)
	buffer := make([]byte, 32768)
	if d.rateLimit > 0 {
		// rate limit > 4 Mb (max size for buffer)
		if d.rateLimit >= (4 * 1024 * 1024) {
			a := float64(d.rateLimit) / (4.0 * 1024.0 * 1024.0)
			c := float64(time.Second) / a
			limit.Interval = time.Duration(c)
		}
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

	ts := utils.FromBytesToBiggest(resp.ContentLength)
	cs := utils.FromBytesToBiggest(currentSize)

	if resp.ContentLength > 0 && d.progressBar {
		fmt.Printf("\r\033[K%.2f %s / %.2f %s [%s] %.2f%%\n", cs.Amount, cs.Unit, ts.Amount, ts.Unit, pb.ProgressBar(currentSize, resp.ContentLength), float64(currentSize)/float64(resp.ContentLength)*100)
	}

	d.Result <- downloadResult{
		Size: currentSize,
		File: d.path + fileName,
		Err:  nil,
	}
}

func oneSecondTick(totalSize *int64, r *http.Response, ticker *time.Ticker, lastTime *time.Time, lastSize *int64) {
	for range ticker.C {
		cs := utils.FromBytesToBiggest(*totalSize)
		ts := utils.FromBytesToBiggest(r.ContentLength)

		timeInterval := int64(time.Since(*lastTime))
		speed := float64(*totalSize-*lastSize) / (float64(timeInterval) / float64(time.Second))
		if speed == 0 {
			speed = 1
		}
		temp := time.Now()
		lastTime = &temp
		*lastSize = *totalSize

		downloadSpeed := utils.FromBytesToBiggest(int64(speed))
		timeRemaining := (r.ContentLength - *totalSize) / int64(speed)
		fmt.Printf("\r\033[K%.2f %s / %.2f %s [%s] %.2f%%  %.2f %s/s %vs",
			cs.Amount, cs.Unit,
			ts.Amount, ts.Unit,
			pb.ProgressBar(*totalSize, r.ContentLength),
			float64(*totalSize)/float64(r.ContentLength)*100,
			downloadSpeed.Amount, downloadSpeed.Unit,
			timeRemaining)
	}
}
