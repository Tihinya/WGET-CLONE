package downloader

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"wget/packages/utils"
)

type FileInfo struct {
	path     []string
	fileType string
	// baseUrl  string
	fileName string
}

func (d *downloader) DownloadSite(url string, types, dirs []string) error {
	basePath := strings.SplitN(url, "://", 2)[1]
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	err := utils.CreateFolder(basePath)
	if err != nil {
		log.Fatalln(err)
	}

	html, err := DownloadHtml(url, basePath)
	if err != nil {
		log.Fatalln(err)
	}

	regs := []*regexp.Regexp{regexp.MustCompile(`(src|href)=['"]([0-9a-zA-Z./_-]+)['"]`), regexp.MustCompile(`url\(['"]([0-9a-zA-Z./_-]+)['"]\)`)}

	files := []*FileInfo{}
	for _, r := range regs {
		matches := r.FindAllStringSubmatch(html, -1)
		for _, match := range matches {
			if len(match) > 0 {
				url := match[len(match)-1]
				dir, fileName := path.Split(url)
				fileType := strings.Split(fileName, ".")
				file := &FileInfo{
					path:     strings.Split(dir, "/"),
					fileName: fileName,
					fileType: fileType[len(fileType)-1],
				}
				files = append(files, file)
			}
		}
	}

	allowed := make([]*FileInfo, 0)

	for _, file := range files {
		if utils.IsContainsArr(file.path, dirs) || utils.IsContains(types, file.fileType) {
			continue
		}

		allowed = append(allowed, file)
	}

	// download each file which allowed
	for _, file := range allowed {
		err = DownloadFile(url, basePath, *file)
		if err != nil {
			return err
		}
	}
	return nil
}

func DownloadHtml(url, trimmedUrl string) (string, error) {
	r, err := http.Get(url)
	if err != nil {
		return "", err
	}

	file, err := os.Create(path.Join(trimmedUrl, "index.html"))
	if err != nil {
		return "", err
	}

	str, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	_, err = file.Write(str)
	if err != nil {
		return "", err
	}

	return string(str), nil
}

func DownloadFile(baseUrl, basePath string, fileInfo FileInfo) error {
	dirPath := path.Join(fileInfo.path...)
	fileURL := baseUrl + path.Join(dirPath, fileInfo.fileName)

	response, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	err = utils.CreateFolder(path.Join(basePath, dirPath))
	if err != nil {
		return err
	}

	file, err := os.Create(path.Join(basePath, dirPath, fileInfo.fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
