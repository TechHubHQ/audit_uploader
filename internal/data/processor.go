package data

import (
	"audituploader/internal/helpers"
	"audituploader/internal/log"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

func prepareRows(f *excelize.File, month int, rawRecords []RawAuditRecord) ([]RawAuditRecord, error) {
	sheetName := f.GetSheetName(month)
	log.Debug("processing", "month", sheetName)
	rows, err := f.Rows(sheetName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	isHeader := true
	var headers []string
	rowNum := 0
	for rows.Next() {
		rowNum++
		cols, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		cleanedCols := trimCols(cols)

		if isHeader {
			// consider first row as header
			headers = cleanedCols
			log.Debug("Read header row", "row_number", rowNum, "headers", headers)
			isHeader = false
			continue
		}

		rowData := make(map[string]string, len(cleanedCols))
		for i, val := range cleanedCols {
			key := ""
			if i < len(headers) && headers[i] != "" {
				key = strings.ToLower(
					helpers.RemoveSpaces(headers[i]),
				)
			} else {
				key = "column_" + fmt.Sprint(i+1)
			}
			rowData[key] = val
		}
		rawRecords = append(rawRecords, buildRawAuditRecord(rowData))
	}

	if err := rows.Error(); err != nil {
		return nil, err
	}
	log.Debug("Processed raw records", "count", len(rawRecords))
	return rawRecords, nil
}
