package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	pb "wget/packages/progress-bar"
)

func Download(url, fileName, dirPath string, rateLimit int) (int64, error) {
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

	limit := pb.NewLimitedReader(resp.Body, rateLimit, &currentSize)
	buffer := make([]byte, 32768)
	if rateLimit > 0 {
		buffer = make([]byte, rateLimit)
	}

	for {
		n, err := limit.Read(buffer)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if n == 0 {
			break
		}

		// Write the chunk to disk
		if _, err := file.Write(buffer[:n]); err != nil {
			return 0, err
		}
	}

	if resp.ContentLength > 0 {
		fmt.Printf("\r[%s] %.2f%% of %d bytes\n", pb.ProgressBar(currentSize, resp.ContentLength), float64(currentSize)/float64(resp.ContentLength)*100, resp.ContentLength)
	}

	return currentSize, nil
}

func oneSecondTick(size *int64, r *http.Response) {
	ticker := time.NewTicker(time.Second)

	for range ticker.C {
		fmt.Printf("\r[%s] %.2f%% of %d bytes", pb.ProgressBar(*size, r.ContentLength), float64(*size)/float64(r.ContentLength)*100, r.ContentLength)
	}
}
