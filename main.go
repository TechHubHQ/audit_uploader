package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// ------------------------------------------------------------
// TYPES
// ------------------------------------------------------------

type FormField struct {
	Question string   `json:"question"`
	EntryID  string   `json:"entry_id"`
	Type     int      `json:"type"`
	Options  []string `json:"options,omitempty"`
}

type AuditRecord struct {
	RowNumber int
	Values    map[string]string
}

type SubmissionResult struct {
	ExcelRow int
	UHID     string
	Date     string
	Status   string
	Error    string
}

// ------------------------------------------------------------
// GOOGLE FORM SUBMISSION
// ------------------------------------------------------------

func submitForm(payload url.Values) error {
	const formURL = "https://docs.google.com/forms/u/0/d/e/1FAIpQLScuaiAFmILMbBua-wxXuQqh-3_uJAg_bxmSbFNix8kiw5LiGQ/formResponse"

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.PostForm(formURL, payload)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return fmt.Errorf(
			"Google Form returned HTTP status %s",
			resp.Status,
		)
	}

	return nil
}

// ------------------------------------------------------------
// WRITE SUBMISSION LOG
// ------------------------------------------------------------

func writeSubmissionLog(result SubmissionResult) error {
	const filename = "submission_log.csv"

	_, err := os.Stat(filename)
	fileExists := err == nil

	file, err := os.OpenFile(
		filename,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if !fileExists {
		err := writer.Write([]string{
			"ExcelRow",
			"UHID",
			"AuditDate",
			"Status",
			"Error",
		})
		if err != nil {
			return err
		}
	}

	return writer.Write([]string{
		strconv.Itoa(result.ExcelRow),
		result.UHID,
		result.Date,
		result.Status,
		result.Error,
	})
}

// ------------------------------------------------------------
// NORMALIZE
// ------------------------------------------------------------

func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, "\u00a0", " ")

	var b strings.Builder

	for _, r := range s {
		switch r {
		case '/', '-', '_', ',', '.', '(', ')', '&':
			b.WriteRune(' ')
		default:
			b.WriteRune(r)
		}
	}

	return strings.Join(strings.Fields(b.String()), " ")
}

// ------------------------------------------------------------
// LOAD FORM FIELDS
// ------------------------------------------------------------

func loadFormFields() ([]FormField, error) {
	data, err := os.ReadFile("form_fields.json")
	if err != nil {
		return nil, err
	}

	var fields []FormField

	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}

	return fields, nil
}

// ------------------------------------------------------------
// LOAD MAPPING
// ------------------------------------------------------------

func loadMapping() (map[string]string, error) {
	data, err := os.ReadFile("mapping.json")
	if err != nil {
		return nil, err
	}

	var mapping map[string]string

	if err := json.Unmarshal(data, &mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}

// ------------------------------------------------------------
// FIND FORM FIELD
// ------------------------------------------------------------

func findFormField(
	excelHeader string,
	fields []FormField,
	mapping map[string]string,
) (FormField, bool) {

	formQuestion := excelHeader

	if mapped, ok := mapping[excelHeader]; ok {
		formQuestion = mapped
	}

	target := normalize(formQuestion)

	for _, field := range fields {
		if normalize(field.Question) == target {
			return field, true
		}
	}

	return FormField{}, false
}

// ------------------------------------------------------------
// PARSE AUDIT DATE
// ------------------------------------------------------------

func parseAuditDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)

	layouts := []string{
		"01-02-06",
		"01-02-2006",
		"01/02/06",
		"01/02/2006",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf(
		"unsupported date format %q",
		value,
	)
}

// ------------------------------------------------------------
// NORMALIZE CHOICE
// ------------------------------------------------------------

func normalizeChoice(
	value string,
	options []string,
) (string, error) {

	value = strings.TrimSpace(value)

	if value == "" {
		return "", nil
	}

	for _, option := range options {
		if strings.EqualFold(
			strings.TrimSpace(option),
			value,
		) {
			return option, nil
		}
	}

	return "", fmt.Errorf(
		"value %q is not one of form options %v",
		value,
		options,
	)
}

// ------------------------------------------------------------
// VALIDATE CHOICE
// ------------------------------------------------------------

func validateChoice(
	field FormField,
	value string,
) error {

	if value == "" {
		return nil
	}

	for _, option := range field.Options {
		if strings.EqualFold(
			strings.TrimSpace(option),
			strings.TrimSpace(value),
		) {
			return nil
		}
	}

	return fmt.Errorf(
		"%q is not a valid option %v",
		value,
		field.Options,
	)
}

// ------------------------------------------------------------
// INTEGER CHECK
// ------------------------------------------------------------

func isInteger(value string) bool {
	if value == "" {
		return true
	}

	_, err := strconv.Atoi(
		strings.TrimSpace(value),
	)

	return err == nil
}

// ------------------------------------------------------------
// CONDITIONAL COUNT VALIDATION
//
// Blank count fields are allowed.
// NA is allowed.
// Populated values must be non-negative integers.
// ------------------------------------------------------------

func validateConditionalCounts(
	record AuditRecord,
	fields []FormField,
	mapping map[string]string,
) []string {

	var errors []string

	excelByFormQuestion := make(map[string]string)

	for excelHeader := range record.Values {

		if strings.EqualFold(
			strings.TrimSpace(excelHeader),
			"NA",
		) {
			continue
		}

		field, ok := findFormField(
			excelHeader,
			fields,
			mapping,
		)

		if !ok {
			continue
		}

		excelByFormQuestion[normalize(field.Question)] = excelHeader
	}

	for i := 0; i < len(fields); i++ {

		if fields[i].Type != 2 {
			continue
		}

		var countField *FormField

		for j := i + 1; j < len(fields); j++ {

			if fields[j].Type == 2 {
				break
			}

			if fields[j].Type == 0 {
				countField = &fields[j]
				break
			}
		}

		if countField == nil {
			continue
		}

		excelCountHeader, ok :=
			excelByFormQuestion[normalize(countField.Question)]

		if !ok {
			continue
		}

		count := strings.TrimSpace(
			record.Values[excelCountHeader],
		)

		// Blank is valid.
		if count == "" {
			continue
		}

		// NA is valid.
		if strings.EqualFold(count, "NA") {
			continue
		}

		n, err := strconv.Atoi(count)

		if err != nil {
			errors = append(
				errors,
				fmt.Sprintf(
					"%s must be a non-negative integer; got %q",
					excelCountHeader,
					count,
				),
			)

			continue
		}

		if n < 0 {
			errors = append(
				errors,
				fmt.Sprintf(
					"%s cannot be negative; got %d",
					excelCountHeader,
					n,
				),
			)
		}
	}

	return errors
}

// ------------------------------------------------------------
// VALIDATE ROW
// ------------------------------------------------------------

func validateRow(
	record AuditRecord,
	fields []FormField,
	mapping map[string]string,
) []string {

	var errors []string

	get := func(header string) string {
		return strings.TrimSpace(
			record.Values[header],
		)
	}

	// --------------------------------------------------------
	// REQUIRED FIELDS
	// --------------------------------------------------------

	required := []string{
		"UHID/IP Number",
		"Doctor Name",
		"Department",
		"Audit Date",
		"Location",
	}

	for _, header := range required {

		if get(header) == "" {
			errors = append(
				errors,
				fmt.Sprintf(
					"%s is empty",
					header,
				),
			)
		}
	}

	// --------------------------------------------------------
	// AUDIT DATE
	// --------------------------------------------------------

	if dateValue := get("Audit Date"); dateValue != "" {

		t, err := parseAuditDate(
			dateValue,
		)

		if err != nil {
			errors = append(
				errors,
				err.Error(),
			)
		} else if t.Month() != time.August {

			errors = append(
				errors,
				fmt.Sprintf(
					"audit date %s is not in AUGUST",
					t.Format("2006-01-02"),
				),
			)
		}
	}

	// --------------------------------------------------------
	// VALIDATE POPULATED EXCEL FIELDS
	// --------------------------------------------------------

	for excelHeader, value := range record.Values {

		excelHeader = strings.TrimSpace(excelHeader)
		value = strings.TrimSpace(value)

		if excelHeader == "" {
			continue
		}

		if strings.EqualFold(
			excelHeader,
			"NA",
		) {
			continue
		}

		if value == "" {
			continue
		}

		if strings.EqualFold(
			value,
			"NA",
		) {
			continue
		}

		field, ok := findFormField(
			excelHeader,
			fields,
			mapping,
		)

		if !ok {

			errors = append(
				errors,
				fmt.Sprintf(
					"no Form mapping for %q",
					excelHeader,
				),
			)

			continue
		}

		// Choice fields.
		if field.Type == 2 {

			if err := validateChoice(
				field,
				value,
			); err != nil {

				errors = append(
					errors,
					fmt.Sprintf(
						"%s: %v",
						excelHeader,
						err,
					),
				)
			}
		}
	}

	// --------------------------------------------------------
	// NUMERIC COUNT FIELDS
	// --------------------------------------------------------

	countFields := []string{
		"Total number of drugs in the prescription",

		"How many drugs did not have doses stated appropriately",

		"How many drugs did not have Frequency stated appropriately",

		"How many drugs did no have Route stated approriately",

		"Howmany drugs not mentioned Units",

		"Howmany drugs did not mentioned concentration",

		"Howmany drugs didnot have rate of administration not mentioned",

		"Howmany drugs did have incorrect drug selection",

		"Howmany drugs are illegible",

		"How many Non approved abbreviations used",

		"How many drug names were not written in capital letters",

		"How many drugs had non modification of drug dose keeping in mind drug drug interactions",

		"How many drugs had non modification of time of administration or dose keeping in mind drug food interactions",

		"Howmany Wrong Formulation Transcribed/Indented",

		"Howmany wrong drugs Transcribed/Indented",

		"How many Wrong Strengths Transcribed/Indented",

		"How many wrong drugs dispensed",

		"Howmany wrong dose dispensed",

		"How many wrong drug formulations dispensed",

		"How many Expired /Near Expiry drugs dispensed",

		"How many drugs dispensed in wrong /No drug labelling",

		"How many drugs were dispensed lately",

		"How many drugs generic substitute done without consultation",

		"How many drugs were administered to wrong patient",

		"How many drugs were omitted to the patient",

		"How many drug doses were adminstered improperly",

		"How many wrong drugs were adminstered",

		"How many wrong dosage forms were adminstered",

		"How many wrong routes of drugs administered",

		"How many drugs were adminstered in wrong Rate",

		"Howmany drugs were not administered at correct duration",

		"How many drugs were administered in wrong time",

		"How many drugs documentation was not done properly",

		"Howmany drugs documentation completly and properly not done by nursing staff",

		"How many drugs documented without administration",
	}

	for _, header := range countFields {

		value := get(header)

		if value == "" {
			continue
		}

		if strings.EqualFold(value, "NA") {
			continue
		}

		if !isInteger(value) {

			errors = append(
				errors,
				fmt.Sprintf(
					"%s must be an integer, got %q",
					header,
					value,
				),
			)
		}
	}

	// --------------------------------------------------------
	// CONDITIONAL COUNTS
	// --------------------------------------------------------

	errors = append(
		errors,
		validateConditionalCounts(
			record,
			fields,
			mapping,
		)...,
	)

	return errors
}

// ------------------------------------------------------------
// BUILD GOOGLE FORM PAYLOAD
//
// This function ONLY builds the payload.
// It NEVER sends an HTTP request.
// ------------------------------------------------------------

func BuildPayload(
	record AuditRecord,
	fields []FormField,
	mapping map[string]string,
) (url.Values, error) {

	values := url.Values{}

	for excelHeader, rawValue := range record.Values {

		excelHeader = strings.TrimSpace(
			excelHeader,
		)

		rawValue = strings.TrimSpace(
			rawValue,
		)

		// Empty header.
		if excelHeader == "" {
			continue
		}

		// Ignore Excel column named NA.
		if strings.EqualFold(
			excelHeader,
			"NA",
		) {
			continue
		}

		// NA means Not Applicable.
		if strings.EqualFold(
			rawValue,
			"NA",
		) {
			continue
		}

		// Blank Excel cells are not submitted.
		if rawValue == "" {
			continue
		}

		field, ok := findFormField(
			excelHeader,
			fields,
			mapping,
		)

		if !ok {

			return nil, fmt.Errorf(
				"no Form field found for Excel column %q",
				excelHeader,
			)
		}

		value := rawValue

		// ----------------------------------------------------
		// DATE
		// ----------------------------------------------------

		if field.Type == 9 {

			t, err := parseAuditDate(
				rawValue,
			)

			if err != nil {

				return nil, fmt.Errorf(
					"%s: %w",
					excelHeader,
					err,
				)
			}

			value = t.Format(
				"2006-01-02",
			)
		}

		// ----------------------------------------------------
		// CHOICE
		// ----------------------------------------------------

		if field.Type == 2 {

			normalized, err := normalizeChoice(
				value,
				field.Options,
			)

			if err != nil {

				return nil, fmt.Errorf(
					"%s: %w",
					excelHeader,
					err,
				)
			}

			value = normalized
		}

		values.Set(
			"entry."+field.EntryID,
			value,
		)
	}

	return values, nil
}

// ------------------------------------------------------------
// WRITE DRY-RUN CSV
// ------------------------------------------------------------

func writeDryRunReport(
	results [][]string,
) error {

	file, err := os.Create(
		"dry_run.csv",
	)

	if err != nil {
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)

	header := []string{
		"ExcelRow",
		"UHID",
		"AuditDate",
		"FieldCount",
		"Status",
		"Error",
	}

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, result := range results {

		if err := writer.Write(result); err != nil {
			return err
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}

// ------------------------------------------------------------
// MAIN
// ------------------------------------------------------------

func main() {

	valid := 0
	invalid := 0
	submitted := 0

	submitMode := flag.Bool(
		"submit",
		false,
		"submit validated rows to Google Forms",
	)

	limit := flag.Int(
		"limit",
		0,
		"maximum number of READY rows to submit; 0 means unlimited",
	)

	flag.Parse()

	// --------------------------------------------------------
	// MODE
	// --------------------------------------------------------

	if *submitMode {

		fmt.Println("SUBMISSION MODE")

		if *limit > 0 {
			fmt.Printf(
				"Submission limit: %d\n",
				*limit,
			)
		} else {
			fmt.Println(
				"Submission limit: unlimited",
			)
		}

	} else {
		fmt.Println("DRY RUN")
	}

	fmt.Println(
		"========================================",
	)

	fmt.Println(
		"AUGUST VALIDATION + DRY RUN",
	)

	fmt.Println(
		"========================================",
	)

	// --------------------------------------------------------
	// OPEN EXCEL
	// --------------------------------------------------------

	file, err := excelize.OpenFile(
		"IP Rx AUDIT.xlsx",
	)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	// --------------------------------------------------------
	// READ AUGUST SHEET
	// --------------------------------------------------------

	rows, err := file.GetRows(
		"AUGUST",
	)

	if err != nil {
		panic(err)
	}

	if len(rows) < 2 {
		panic(
			"AUGUST contains no audit records",
		)
	}

	// --------------------------------------------------------
	// LOAD FORM SCHEMA
	// --------------------------------------------------------

	fields, err := loadFormFields()

	if err != nil {
		panic(err)
	}

	// --------------------------------------------------------
	// LOAD EXCEL → FORM MAPPING
	// --------------------------------------------------------

	mapping, err := loadMapping()

	if err != nil {
		panic(err)
	}

	// --------------------------------------------------------
	// HEADERS
	// --------------------------------------------------------

	headers := rows[0]

	fmt.Printf(
		"Data rows: %d\n\n",
		len(rows)-1,
	)

	// Store one CSV row for every Excel row.
	dryRunResults := make(
		[][]string,
		0,
		len(rows)-1,
	)

	// --------------------------------------------------------
	// PROCESS EACH EXCEL ROW
	// --------------------------------------------------------

	for rowIndex := 1; rowIndex < len(rows); rowIndex++ {

		row := rows[rowIndex]

		record := AuditRecord{
			RowNumber: rowIndex + 1,
			Values:    make(map[string]string),
		}

		// ----------------------------------------------------
		// CONVERT EXCEL ROW → AuditRecord
		// ----------------------------------------------------

		for colIndex, header := range headers {

			header = strings.TrimSpace(
				header,
			)

			if header == "" {
				continue
			}

			value := ""

			if colIndex < len(row) {
				value = strings.TrimSpace(
					row[colIndex],
				)
			}

			record.Values[header] = value
		}

		uhid := record.Values["UHID/IP Number"]
		auditDate := record.Values["Audit Date"]

		// ----------------------------------------------------
		// VALIDATE ROW
		// ----------------------------------------------------

		errors := validateRow(
			record,
			fields,
			mapping,
		)

		// ----------------------------------------------------
		// INVALID ROW
		// ----------------------------------------------------

		if len(errors) > 0 {

			invalid++

			errorText := strings.Join(
				errors,
				"; ",
			)

			fmt.Printf(
				"\nRow %-3d UHID %-15s INVALID\n",
				record.RowNumber,
				uhid,
			)

			for _, validationError := range errors {

				fmt.Printf(
					"    - %s\n",
					validationError,
				)
			}

			dryRunResults = append(
				dryRunResults,
				[]string{
					strconv.Itoa(
						record.RowNumber,
					),
					uhid,
					auditDate,
					"0",
					"BLOCKED",
					errorText,
				},
			)

			continue
		}

		// ----------------------------------------------------
		// BUILD PAYLOAD
		// ----------------------------------------------------

		payload, err := BuildPayload(
			record,
			fields,
			mapping,
		)

		if err != nil {

			invalid++

			fmt.Printf(
				"\nRow %-3d UHID %-15s INVALID PAYLOAD\n",
				record.RowNumber,
				uhid,
			)

			fmt.Println(
				"    -",
				err,
			)

			dryRunResults = append(
				dryRunResults,
				[]string{
					strconv.Itoa(
						record.RowNumber,
					),
					uhid,
					auditDate,
					"0",
					"BLOCKED",
					err.Error(),
				},
			)

			continue
		}

		// ----------------------------------------------------
		// READY
		// ----------------------------------------------------

		valid++

		fmt.Printf(
			"Row %-3d UHID %-15s READY  fields=%d\n",
			record.RowNumber,
			uhid,
			len(payload),
		)

		dryRunResults = append(
			dryRunResults,
			[]string{
				strconv.Itoa(record.RowNumber),
				uhid,
				auditDate,
				strconv.Itoa(len(payload)),
				"READY",
				"",
			},
		)

		// ----------------------------------------------------
		// REAL SUBMISSION
		//
		// HTTP request happens ONLY here.
		// ----------------------------------------------------

		if !*submitMode {
			continue
		}

		// Respect --limit.
		if *limit > 0 && submitted >= *limit {
			continue
		}

		fmt.Printf(
			"    SUBMITTING row %d...\n",
			record.RowNumber,
		)

		// ----------------------------------------------------
		// SEND HTTP REQUEST
		// ----------------------------------------------------

		err = submitForm(payload)

		if err != nil {

			fmt.Printf(
				"    SUBMISSION FAILED: %v\n",
				err,
			)

			if logErr := writeSubmissionLog(
				SubmissionResult{
					ExcelRow: record.RowNumber,
					UHID:     uhid,
					Date:     auditDate,
					Status:   "FAILED",
					Error:    err.Error(),
				},
			); logErr != nil {

				fmt.Printf(
					"    WARNING: could not write submission log: %v\n",
					logErr,
				)
			}

			continue
		}

		// ----------------------------------------------------
		// SUCCESS
		// ----------------------------------------------------

		submitted++

		fmt.Printf(
			"    SUBMITTED successfully\n",
		)

		if err := writeSubmissionLog(
			SubmissionResult{
				ExcelRow: record.RowNumber,
				UHID:     uhid,
				Date:     auditDate,
				Status:   "SUCCESS",
				Error:    "",
			},
		); err != nil {

			fmt.Printf(
				"    WARNING: could not write submission log: %v\n",
				err,
			)
		}
	}

	// --------------------------------------------------------
	// WRITE DRY-RUN CSV
	// --------------------------------------------------------

	if err := writeDryRunReport(
		dryRunResults,
	); err != nil {
		panic(err)
	}

	fmt.Println()

	fmt.Println(
		"Dry-run report written to: dry_run.csv",
	)

	// --------------------------------------------------------
	// FINAL SUMMARY
	// --------------------------------------------------------

	fmt.Println()

	fmt.Println(
		"========================================",
	)

	fmt.Println(
		"SUMMARY",
	)

	fmt.Println(
		"========================================",
	)

	fmt.Printf(
		"Ready     : %d\n",
		valid,
	)

	fmt.Printf(
		"Blocked   : %d\n",
		invalid,
	)

	fmt.Printf(
		"Submitted : %d\n",
		submitted,
	)

	fmt.Printf(
		"Total     : %d\n",
		len(rows)-1,
	)

	// --------------------------------------------------------
	// HTTP STATUS
	// --------------------------------------------------------

	if *submitMode {

		fmt.Printf(
			"HTTP submissions sent: %d\n",
			submitted,
		)

	} else {

		fmt.Println(
			"NO HTTP REQUESTS WERE SENT.",
		)
	}

	fmt.Println(
		"========================================",
	)
}
