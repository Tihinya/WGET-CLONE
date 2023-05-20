package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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

	if storage.HasFlag("B") {
		file, err := os.OpenFile("wget-log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer file.Close()

		fmt.Println(`Output will be written to "wget-log".`)

		// cmd := exec.Command(os.Args[0], "https://golang.org/dl/go1.16.3.linux-amd64.tar.gz")
		cmd := exec.Command(os.Args[0], storage.ArgsExcluded("B")...)
		cmd.Stdout = file
		cmd.Stderr = file

		cmd.Start()
		os.Exit(0)
	}

	urls := make([]string, 0)

	if urlArg := storage.GetTags(); len(urlArg) > 0 {
		urls = append(urls, urlArg[0])
	}

	fmt.Printf("Start at %s\n", formatTime(time.Now()))

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

	if !storage.HasFlag("mirror") {
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
	} else {
		types := []string{}
		dirs := []string{}

		if flag, err := storage.GetFlag("reject"); err == nil {
			types = strings.Split(flag.GetValue(), ",")
		}
		if flag, err := storage.GetFlag("exclude"); err == nil {
			temp := strings.Split(flag.GetValue(), ",")
			for _, v := range temp {
				dirs = append(dirs, strings.TrimPrefix(v, "/"))
			}
		}

		err = down.DownloadSite(urls[0], types, dirs)
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("Finished at %s\n", formatTime(time.Now()))
}

func formatTime(t time.Time) string {
	return t.Format(time.DateTime)
}

func bytesToMb(bytes int64) float64 {
	return float64(bytes) / (1024 * 1024)
}
