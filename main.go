package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
	flag_parser "wget/packages/flag-parser"
)

const progressBarWidth = 40
const outputFormat = "Content size: %d bytes [~ %.2f Mb]\nSaving file to: %s\nFinished at %s\n"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go URL")
		return
	}

	flags := os.Args[1 : len(os.Args)-1]
	url := os.Args[len(os.Args)-1]

	storage, err := flag_parser.CreateParser().
		Add("B", "backgound download. When the program containing this flag is executed it should output : Output will be written to `wget-log`", true).
		Add("O", "specifies file name", false).
		Add("P", "specifies file location", false).
		Add("rate-limit", "specifies limit of speed rate", false).
		Add("i", "asynchronously download multiple files from given URLs", false).
		Add("mirror", "download entire website", true).
		Add("reject", "specifies which file suffixes will be avoided", false).
		Add("exclude", "specifies which paths will be avoided", false).
		Parse(flags)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	storage.GetFlag("B")

	fmt.Printf("Start at %s\n", formatTime(time.Now()))

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error fetching URL: recieved status code", resp.StatusCode)
	} else {
		fmt.Println("Sending request, awaiting response... status 200 OK")
	}

	fileName := ""
	dirPath := ""

	if flag, err := storage.GetFlag("O"); err == nil {
		fileName = flag.GetValue()
	} else {
		fileName = path.Base(url)
	}

	if flag, err := storage.GetFlag("P"); err == nil {
		dirPath = flag.GetValue()
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			err := os.Mkdir(dirPath, fs.ModeDir|0755)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	file, err := os.Create(dirPath + fileName)
	if err != nil {
		fmt.Printf("Error creating file %v", err)
		return
	}

	defer file.Close()

	currentSize := int64(0)

	if resp.ContentLength > 0 {
		go oneSecondTick(&currentSize, resp)
	}

	_, err = io.Copy(file, io.TeeReader(resp.Body, &progressReader{&currentSize}))
	if err != nil {
		fmt.Printf("Cannot download file %v", err)
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
	return strings.Repeat("=", progress) + strings.Repeat("-", progressBarWidth-progress)
}

func formatTime(t time.Time) string {
	return t.Format(time.DateTime)
}

func bytesToMb(bytes int64) float64 {
	return float64(bytes) / (1024 * 1024)
}
