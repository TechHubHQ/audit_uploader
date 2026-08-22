package data

import (
	"audituploader/internal/log"
	"fmt"
	"sync"

	"github.com/xuri/excelize/v2"
)

func ReadExcel(file string, months []int, month int) ([]RawAuditRecord, error) {
	f, err := excelize.OpenFile(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	log.Info("Opened Excel file", "file", file, "sheets", len(f.GetSheetList()))

	var rawRecords []RawAuditRecord
	if month != 0 {
		log.Info("Processing month", "month", month+1)
		processedRecords, err := prepareRows(f, month, rawRecords)
		if err != nil {
			return nil, err
		}
		rawRecords = append(rawRecords, processedRecords...)
		log.Info("Records built for month", "month", month+1, "count", len(processedRecords))
	} else if len(months) != 0 {
		// read data for the specified month range
		var wg sync.WaitGroup

		results := make([][]RawAuditRecord, len(months))
		for i, m := range months {
			wg.Go(func() {
				processedRecords, err := prepareRows(f, m, nil)
				if err != nil {
					log.Error("Error processing month", "month", m+1, "error", err)
					return
				}
				results[i] = processedRecords
				log.Info("Records built for month", "month", m+1, "count", len(processedRecords))
			})
		}

		wg.Wait()

		for _, r := range results {
			rawRecords = append(rawRecords, r...)
		}
	} else {
		// process latest sheet
		sheetList := f.GetSheetList()
		if len(sheetList) == 0 {
			return nil, fmt.Errorf("no sheets found in the Excel file")
		}
		log.Info("Processing latest sheet", "sheet", sheetList[len(sheetList)-1])
		processedRecords, err := prepareRows(f, len(sheetList)-1, rawRecords)
		if err != nil {
			return nil, err
		}
		rawRecords = append(rawRecords, processedRecords...)
		log.Info("Records built for latest sheet", "sheet", sheetList[len(sheetList)-1], "count", len(processedRecords))
	}

	log.Info("Finished reading Excel file", "total_records", len(rawRecords))
	return rawRecords, nil
}
