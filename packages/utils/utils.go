package utils

import (
	"bufio"
	"io/fs"
	"os"
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
		err := os.Mkdir(path, fs.ModeDir|0755)
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
