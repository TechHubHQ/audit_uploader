package helpers

import (
	"audituploader/internal/log"
	"strconv"
	"strings"
)

func ParseMonthRange(monthRange string) []int {
	bounds := strings.SplitN(monthRange, "-", 2)
	if len(bounds) != 2 {
		log.Error("Invalid month range format", "input", monthRange)
		return nil
	}

	start, err := strconv.Atoi(strings.TrimSpace(bounds[0]))
	if err != nil {
		log.Error("Invalid start month in range", "month", bounds[0], "error", err)
		return nil
	}
	end, err := strconv.Atoi(strings.TrimSpace(bounds[1]))
	if err != nil {
		log.Error("Invalid end month in range", "month", bounds[1], "error", err)
		return nil
	}
	if start > end {
		log.Error("Invalid month range: start greater than end", "start", start, "end", end)
		return nil
	}

	monthInts := make([]int, 0, end-start+1)
	for m := start; m <= end; m++ {
		monthInts = append(monthInts, m)
	}

	log.Debug("Parsed month range", "input", monthRange, "months", monthInts)
	return monthInts
}
