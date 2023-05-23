# wget

## Description

This project aims to recreate the functionality of the `wget` command-line tool using a compiled language. It allows downloading files from URLs, with options to save them under different names or in specific directories. Additional features include setting download speed limits, background downloading, asynchronous downloading of multiple files, and mirroring entire websites.

Audit can be accessed at https://github.com/01-edu/public/blob/master/subjects/wget/audit/README.md

## Usage

To run the program, use the following command:

```console
go run main.go [FLAGS] [URL]
```

Alternatively, you can build a binary file and run it directly:

```console
go build -o wget main.go
./wget [FLAGS] [URL]
```

## Flags

The program supports the following flags:

1. `-B`: Enables background download. Output will be written to "wget-log".
2. `-O=<name.jpg>`: Specifies the file name.
3. `-P=</Download/>`: Specifies the file location.
4. `--rate-limit=<200k, 2M>`: Sets the speed limit for downloads.
5. `-i=<file with links>`: Downloads multiple files asynchronously from the given URLs.
6. `--mirror`: Downloads an entire website.
   1. `--reject=<jpg,gif> -R=<jpg,gif>`: Specifies file suffixes to be avoided during the download.
   2. `--exclude=</assets,/css> -X=</assets,/css>`: Specifies paths to be excluded from the download.

## How Mirror Works

The mirror functionality downloads an HTML file and searches for the following tags:

1. `<a>` with `href` attribute
2. `<link>` with `src` attribute
3. `<script>` with `src` attribute
4. `<img>` with `src` attribute

The mirror only downloads files that are not links to other websites. It also checks CSS files for the `url()` function and attempts to find linked resources.

## Features

- Download files from URLs with various options
  - Save files under different names or in specific directories
  - Limit the download speed
  - Perform background downloads
- Download multiple files asynchronously using a file with links
- Mirror entire website pages
- Display information such as:
  - Start time of the download (in the format yyyy-mm-dd hh:mm:ss)
  - Status of the request (OK or error)
  - Size of the downloaded content (in Mb or Gb)
  - Name and path of the saved file
  - Progress bar showing the amount downloaded, percentage, and remaining time
  - Finish time of the download (in the format yyyy-mm-dd hh:mm:ss)

Example output:

```console
Start at 2017-10-14 03:46:06
Sending request, awaiting response... Status: 200 OK
Content size: 56370 [~0.06MB]
Saving file to: ./EMtmPFLWkAA8CIS.jpg
 55.05 KiB / 55.05 KiB [==================================] 100.00% 1.24 MiB/s 0s

Downloaded [https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg]
Finished at 2017-10-14 03:46:07
```
