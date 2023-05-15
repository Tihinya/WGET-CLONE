# wget

## Features

- download file with given URL
  - ✅ under a different name
  - ✅ in specific directory
  - ✅ with limiting the rate speed of a download
  - download file in backgorund
- download multiple files by reading links from specified file
- download entire website page
- display:
  - ✅ time when download started(format yyyy-mm-dd hh:mm:ss)
  - ✅ status request if OK, otherwise show error
  - size of downloaded content(Mb, Gb)
  - ✅ name and path of saved file
  - progress bar:
    - size which already downloaded in Kb or Mb
    - ✅ percantage
    - remaining time
  - ✅ time when download finished(format yyyy-mm-dd hh:mm:ss)

Example of output:

```console
start at 2017-10-14 03:46:06
sending request, awaiting response... status 200 OK
content size: 56370 [~0.06MB]
saving file to: ./EMtmPFLWkAA8CIS.jpg
 55.05 KiB / 55.05 KiB [================================================================================================================] 100.00% 1.24 MiB/s 0s

Downloaded [https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg]
finished at 2017-10-14 03:46:07
```

## Flags

1. `-B` backgound download. When the program containing this flag is executed it should output : Output will be written to "wget-log"
2. `-O=<name.jpg>` specifies file name
3. `-P=</Download/>` specifies file location
4. `--rate-limit=<200k, 2M>` specifies limit of speed rate
5. `-i=<file with links>` asynchronously download multiple files from given URLs
6. `--mirror` download entire website
   1. `--reject=<jpg,gif> -R=<jpg,gif>` specifies which file suffixes will be avoided
   2. `--exclude=</assets,/css> -X=</assets,/css>` specifies which paths will be avoided
