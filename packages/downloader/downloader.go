package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const progressBarWidth = 40

func Download(url, fileName, dirPath string, limit int64) (int64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Sending request, awaiting response... status 200 OK\n")
	} else {
		fmt.Printf("Error downloading %s: recieved status code %d\n", url, resp.StatusCode)
	}

	file, err := os.Create(dirPath + fileName)
	if err != nil {
		return 0, err
	}

	defer file.Close()

	currentSize := int64(0)

	if resp.ContentLength > 0 {
		go oneSecondTick(&currentSize, resp)
	}

	if limit > 0 {
		buffer := make([]byte, limit)
		for {
			n, err := io.TeeReader(resp.Body, &progressReader{&currentSize}).Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				return 0, err
			}
			_, err = file.Write(buffer[:n])
			if err != nil {
				return 0, err
			}
			time.Sleep(time.Second)
		}
	} else {
		_, err = io.Copy(file, io.TeeReader(resp.Body, &progressReader{&currentSize}))
		if err != nil {
			return 0, err
		}
	}

	if resp.ContentLength > 0 {
		fmt.Printf("\r[%s] %.2f%% of %d bytes\n", progressBar(currentSize, resp.ContentLength), float64(currentSize)/float64(resp.ContentLength)*100, resp.ContentLength)
	}

	return currentSize, nil
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

func oneSecondTick(size *int64, r *http.Response) {
	ticker := time.NewTicker(time.Second)

	for range ticker.C {
		fmt.Printf("\r[%s] %.2f%% of %d bytes", progressBar(*size, r.ContentLength), float64(*size)/float64(r.ContentLength)*100, r.ContentLength)
	}
}
