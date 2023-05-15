package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
	"time"
	"wget/packages/downloader"
	flag_parser "wget/packages/flag-parser"
	"wget/packages/utils"
)

const outputFormat = "Content size: %d bytes [~ %.2f Mb]\nSaving file to: %s%s\nFinished at %s\n"

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

	fmt.Printf("Start at %s\n", formatTime(time.Now()))

	fileName := ""
	dirPath := ""
	limit := 0

	if flag, err := storage.GetFlag("O"); err == nil {
		fileName = flag.GetValue()
	} else {
		fileName = path.Base(url)
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

	currentSize, err := downloader.Download(url, fileName, dirPath, limit)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf(outputFormat, currentSize, bytesToMb(currentSize), dirPath, fileName, formatTime(time.Now()))
}

func formatTime(t time.Time) string {
	return t.Format(time.DateTime)
}

func bytesToMb(bytes int64) float64 {
	return float64(bytes) / (1024 * 1024)
}
