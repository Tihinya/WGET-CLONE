package utils

import (
	"bufio"
	"io/fs"
	"os"
	"strings"
)

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func CreateFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, fs.ModeDir|0755)
		return err
	}
	return nil
}

func IsContains[V int | string](arr []V, value V) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}

	return false
}

func IsContainsArr(arr []string, include []string) bool {
	for _, v := range arr {
		if IsContains(include, v) {
			return true
		}
	}

	return false
}
func ParsePath(path string) (string, error) {
	currentPath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(path, "~") {
		currentPath += path[1:]
	}
	return currentPath, nil
}
