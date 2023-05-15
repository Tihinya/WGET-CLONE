package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

func SwitchCases(flag string) (int, error) {
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
	case "Kb":
		return amount * 1024, nil
	case "Mb":
		return amount * 1024 * 1024, nil
	case "Gb":
		return amount * 1024 * 1024 * 1024, nil
	default:
		return 0, fmt.Errorf("unrecognized unit in flag")
	}
}
