package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

const progressBarWidth = 60
const outputFormat = "Content size: %d bytes [~ %.2f Mb]\nSaving file to: %s\nFinished at %s\n"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go URL")
		return
	}

	url := os.Args[len(os.Args)-1]

	fmt.Println("Start at", formatTime(time.Now()))

	resp, err := http.Get(url)
	if err !=nil {
		fmt.Println("Error fetching URL:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error fetching URL: recieved status code", resp.StatusCode)
	} else {
		fmt.Println("Sending request, awaiting response... status 200 OK")
	}

	fileName := path.Base(url)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file", err)
		return
	}

	defer file.Close()

	currentSize := int64(0)

	if resp.ContentLength > 0 {
		go oneSecondTick(&currentSize, resp)
	}

	_, err = io.Copy(file, io.TeeReader(resp.Body, &progressReader{&currentSize}))
	if err != nil {
		fmt.Println("Cannot download file", err)
		return
	}

	if resp.ContentLength > 0 {
		fmt.Printf("\r[%s] %.2f%% of %d bytes\n", progressBar(currentSize, resp.ContentLength), float64(currentSize)/float64(resp.ContentLength)*100, resp.ContentLength)
	}
	fmt.Printf(outputFormat, currentSize, bytesToMb(currentSize), fileName, formatTime(time.Now()))
}

func oneSecondTick(size *int64, r *http.Response) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Printf("\r[%s] %.2f%% of %d bytes", progressBar(*size, r.ContentLength), float64(*size)/float64(r.ContentLength)*100, r.ContentLength)
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
	return  strings.Repeat("=", progress) + strings.Repeat("-", progressBarWidth-progress)
}

func formatTime(t time.Time) string {
	return t.Format(time.DateTime)
}

func bytesToMb(bytes int64) float64 {
	return float64(bytes) / (1024*1024)
}
