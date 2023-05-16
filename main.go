package main

import (
	"log"
	"os"
	"strings"
	"time"
	"wget/packages/downloader"
	fp "wget/packages/flag-parser"
	"wget/packages/utils"
)

const outputFormat = "Content size: %d bytes [~ %.2f Mb]\nSaving file to: %s\n"

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Println("Usage: go run main.go URL")
		return
	}

	storage, err := fp.CreateParser().
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

	log.Printf("Start at %s\n", formatTime(time.Now()))

	fileName := ""
	dirPath := ""
	rateLimit := 0

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
		err := utils.CreateFolder(dirPath)
		if err != nil {
			log.Fatalf("Error creating folder: %v\n", err)
		}
	}

	if flag, err := storage.GetFlag("rate-limit"); err == nil {
		rateLimit, err = utils.StringSizeToBytes(flag.GetValue())
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}

	down := downloader.CreateDownloader(dirPath, rateLimit, len(urls) == 1)

	go func() {
		defer close(down.Result)

		for _, url := range urls {
			down.WG.Add(1)

			go down.DownloadFile(url, fileName)
		}

		down.WG.Wait()
	}()

	for result := range down.Result {
		if result.Err != nil {
			log.Fatalln(result.Err)
		}

		log.Printf(outputFormat, result.Size, bytesToMb(result.Size), result.File)
	}

	log.Printf("Finished at %s\n", formatTime(time.Now()))
}

func formatTime(t time.Time) string {
	return t.Format(time.DateTime)
}

func bytesToMb(bytes int64) float64 {
	return float64(bytes) / (1024 * 1024)
}
