package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type FormField struct {
	Question string   `json:"question"`
	EntryID  string   `json:"entry_id"`
	Type     int      `json:"type"`
	Options  []string `json:"options,omitempty"`
}

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

// Excel value → string.
func cellValue(v string) string {
	return strings.TrimSpace(v)
}

// Find the Form field corresponding to an Excel header.
func findFormField(
	excelHeader string,
	fields []FormField,
	mapping map[string]string,
) (FormField, bool) {

	normalizedExcel := normalize(excelHeader)

	// First use explicit semantic mapping.
	formQuestion, exists := mapping[excelHeader]

	if exists {
		normalizedExcel = normalize(formQuestion)
	}

	// Otherwise assume the normalized headers match.
	for _, field := range fields {
		if normalize(field.Question) == normalizedExcel {
			return field, true
		}
	}

	return FormField{}, false
}

// Convert Excel's date representation into the Google Form date value.
func normalizeDate(value string) (string, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return "", nil
	}

	// Our current workbook preview gives values like:
	// 08-04-26
	//
	// This means 08-Apr-2026.
	layouts := []string{
		"01-02-06",
		"02-01-2006",
		"02/01/06",
		"02/01/2006",
		"2006-01-02",
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, value)
		if err == nil {
			return t.Format("2006-01-02"), nil
		}
	}

	return "", fmt.Errorf("cannot parse date %q", value)
}

// Find the exact option spelling used by Google Forms.
//
// Example:
//
// Excel: YES
// Form:  [Yes No]
//
// Result:
//
// Yes
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

func main() {

	// ------------------------------------------------------------
	// Load form schema
	// ------------------------------------------------------------

	fields, err := loadFormFields()
	if err != nil {
		panic(err)
	}

	// ------------------------------------------------------------
	// Load semantic mapping
	// ------------------------------------------------------------

	mapping, err := loadMapping()
	if err != nil {
		panic(err)
	}

	// ------------------------------------------------------------
	// Open Excel
	// ------------------------------------------------------------

	file, err := excelize.OpenFile("IP Rx AUDIT.xlsx")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	rows, err := file.GetRows("AUGUST")
	if err != nil {
		panic(err)
	}

	if len(rows) < 2 {
		panic("AUGUST has no audit rows")
	}

	headers := rows[0]
	row := rows[1]

	// ------------------------------------------------------------
	// Build payload
	// ------------------------------------------------------------

	values := url.Values{}

	for i, excelHeader := range headers {

		excelHeader = strings.TrimSpace(excelHeader)

		if excelHeader == "" {
			continue
		}

		// Ignore the two NA columns.
		if strings.EqualFold(excelHeader, "NA") {
			continue
		}

		if i >= len(row) {
			continue
		}

		rawValue := cellValue(row[i])

		// Don't submit blank Excel cells.
		//
		// This is especially important for conditional count
		// questions.
		if rawValue == "" {
			continue
		}

		field, ok := findFormField(
			excelHeader,
			fields,
			mapping,
		)

		if !ok {
			panic(fmt.Sprintf(
				"no Form field found for Excel column %q",
				excelHeader,
			))
		}

		value := rawValue

		// Date field.
		if field.Type == 9 {
			value, err = normalizeDate(rawValue)
			if err != nil {
				panic(err)
			}
		}

		// Choice field.
		if field.Type == 2 {
			value, err = normalizeChoice(
				rawValue,
				field.Options,
			)

			if err != nil {
				panic(fmt.Sprintf(
					"%s: %v",
					excelHeader,
					err,
				))
			}
		}

		values.Set(
			"entry."+field.EntryID,
			value,
		)
	}

	// ------------------------------------------------------------
	// PRINT ONLY
	// ------------------------------------------------------------

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("GOOGLE FORM PAYLOAD PREVIEW")
	fmt.Println("========================================")

	for key, values := range values {
		fmt.Printf(
			"%s = %q\n",
			key,
			values[0],
		)
	}

	err = Submit(values)
	if err != nil {
		panic(err)
	}

	fmt.Println("Submission successful")
	fmt.Println()
	fmt.Println("Total submitted fields:", len(values))
	fmt.Println()
}
