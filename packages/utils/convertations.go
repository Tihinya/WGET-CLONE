package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

func StringSizeToBytes(flag string) (int, error) {
	re := `^(\d+)((?:G|M|K)?b)$`

	match := regexp.MustCompile(re).FindStringSubmatch(flag)
	if len(match) < 3 {
		return 0, fmt.Errorf("you used unallowed format, input example: --rate-limit=200Kb")
	}

	amountStr := match[1]
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return 0, fmt.Errorf("invalid number in flag: %v", err)
	}

	unit := match[2]
	switch unit {
	case "b":
		return amount, nil
	case "Kb":
		return amount * 1024, nil
	case "Mb":
		return amount * 1024 * 1024, nil
	default:
		return 0, fmt.Errorf("unrecognized unit in flag")
	}
}

type size struct {
	Amount float64
	Unit   string
}

var units = [...]string{"b", "Kb", "Mb", "Gb"}

func FromBytesToBiggest(bytes int64) size {
	i := 0
	temp := float64(bytes)
	for temp >= 1024.0 && i < 3 {
		temp /= 1024.0
		i++
	}

	return size{
		Amount: temp,
		Unit:   units[i],
	}
}
