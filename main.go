package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
	"wget/packages/downloader"
	flag_parser "wget/packages/flag-parser"
	"wget/packages/utils"
)

const outputFormat = "Content size: %d bytes [~ %.2f Mb]\nSaving file to: %s\n"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go URL")
		return
	}

	storage, err := flag_parser.CreateParser().
		Add("B", "backgound download. When the program containing this flag is executed it should output : Output will be written to `wget-log`", true).
		Add("O", "specifies file name", false).
		Add("P", "specifies file location", false).
		Add("rate-limit", "specifies limit of speed rate", false).
		Add("i", "asynchronously download multiple files from given URLs", false).
		Add("mirror", "download entire website", true).
		Add("reject", "specifies which file suffixes will be avoided", false).
		Add("exclude", "specifies which paths will be avoided", false).
		Parse(os.Args[1:])

	if err != nil {
		log.Fatalf("Error parsing flag: %v\n", err)
	}

	urls := make([]string, 0)

	if urlArg := storage.GetTags(); len(urlArg) > 0 {
		urls = append(urls, urlArg[0])
	}

	fmt.Printf("Start at %s\n", formatTime(time.Now()))

	fileName := ""
	dirPath := ""
	limit := 0

	if flag, err := storage.GetFlag("O"); err == nil && !storage.HasFlag("i") {
		fileName = flag.GetValue()
	}

	if flag, err := storage.GetFlag("i"); err == nil {
		urls, err = utils.ReadLines(flag.GetValue())
		if err != nil {
			log.Fatalf("Error reading file: %v\n", err)
		}
	}

	if flag, err := storage.GetFlag("P"); err == nil {
		dirPath = flag.GetValue()
		if !strings.HasSuffix(dirPath, "/") {
			dirPath += "/"
		}

		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			err := os.Mkdir(dirPath, fs.ModeDir|0755)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	if flag, err := storage.GetFlag("rate-limit"); err == nil {
		limit, err = utils.SwitchCases(flag.GetValue())
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}

	ch := make(chan downloader.DownloadResult)

	var wg sync.WaitGroup

	go func() {
		defer close(ch)

		for _, url := range urls {
			wg.Add(1)

			temp := fileName
			if fileName == "" {
				temp = path.Base(url)
			}

			go downloader.Download(url, temp, dirPath, limit, ch, &wg)
		}

		wg.Wait()
	}()

	for result := range ch {
		if result.Err != nil {
			log.Fatalln(result.Err)
		}

		fmt.Printf(outputFormat, result.Size, bytesToMb(result.Size), result.File)
	}

	fmt.Printf("Finished at %s\n", formatTime(time.Now()))
}

func formatTime(t time.Time) string {
	return t.Format(time.DateTime)
}

func bytesToMb(bytes int64) float64 {
	return float64(bytes) / (1024 * 1024)
}
