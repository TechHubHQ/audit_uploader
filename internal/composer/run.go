package composer

import (
	"audituploader/internal/cli"
	"audituploader/internal/data"
	"audituploader/internal/helpers"
	"audituploader/internal/log"
)

func Run() {
	log.Info("===============================================")
	log.Info("		Starting Audit Uploader")
	log.Info("===============================================")
	args := cli.ReadArgs()
	log.Info("Input file", "file", args.InputFile)
	absPath, err := helpers.FindFile(args.InputFile)
	if err != nil {
		log.Error("Error finding file", "file", args.InputFile, "error", err)
		return
	}
	log.Info("Resolved input file", "path", absPath)

	/*
		reduce month by -1 for user readability,
		since the user will provide month as 1-12,
		but we need 0-11 for processing
	*/
	if args.Month != 0 {
		args.Month--
	}

	/*
		reduce month range by -1 for user readability,
		since the user will provide month as 1-12,
		but we need 0-11 for processing
	*/
	var months []int
	if args.MonthRange != "" {
		months = helpers.ParseMonthRange(args.MonthRange)
		for i := range months {
			months[i]--
		}
	}

	rawRecords, err := data.ReadExcel(absPath, months, args.Month)
	if err != nil {
		log.Error("Error reading Excel file", "path", absPath, "error", err)
		return
	}
	log.Info("Records built", "count", len(rawRecords))

	auditRecords := make([]data.AuditRecord, 0, len(rawRecords))
	for _, rawRecord := range rawRecords {
		auditRecord := data.MapToAuditRecord(rawRecord)
		auditRecords = append(auditRecords, auditRecord)
	}
	log.Debug("Mapped Audit Records", "auditRecords", auditRecords)

	log.Info("===============================================")
	log.Info("		Audit Uploader Completed")
	log.Info("===============================================")
}
