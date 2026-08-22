package data

import (
	"audituploader/internal/log"
	"strings"
)

func validateColName(colName string) bool {
	if strings.ToLower(colName) == "na" {
		return false
	}
	return true
}

func trimCols(cols []string) []string {
	trimmedCols := make([]string, 0, len(cols))
	for _, col := range cols {
		trimmedCol := strings.TrimSpace(col)
		if validateColName(trimmedCol) {
			trimmedCols = append(trimmedCols, trimmedCol)
		}
	}
	log.Debug("Trimmed columns", "original_count", len(cols), "trimmed_count", len(trimmedCols))
	return trimmedCols
}
