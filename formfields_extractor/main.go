package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type FormField struct {
	Question string   `json:"question"`
	EntryID  string   `json:"entry_id"`
	Type     int      `json:"type"`
	Options  []string `json:"options,omitempty"`
}

func main() {
	data, err := os.ReadFile("audit_form.html")
	if err != nil {
		panic(err)
	}

	fields := make([]FormField, 0)

	html := string(data)
	const marker = "var FB_PUBLIC_LOAD_DATA_ = "
	start := strings.Index(html, marker)
	if start == -1 {
		panic("form data not found")
	}

	start += len(marker)
	end := strings.Index(html[start:], ";")
	raw := html[start : start+end]

	var formData []any
	err = json.Unmarshal([]byte(raw), &formData)
	if err != nil {
		panic(err)
	}

	// Navigate to questions
	level1 := formData[1].([]any)
	level2 := level1[1].([]any)
	questions := level2

	fmt.Println("Total elements:", len(questions))

	for _, q := range questions {
		question := q.([]any)
		if len(question) < 5 {
			continue
		}

		title, ok := question[1].(string)
		if !ok {
			continue
		}

		fmt.Println("--------------------")
		fmt.Println("Question:", title)
		fmt.Println("Type:", question[3])

		// Extract entry ID
		entryData, ok := question[4].([]any)
		if !ok || len(entryData) == 0 {
			fmt.Println("Entry ID: None")
			continue
		}

		questionType := int(question[3].(float64))

		var options []string

		if questionType == 2 {

			entryData := question[4].([]interface{})

			first := entryData[0].([]interface{})

			optionData, ok := first[1].([]interface{})

			if ok {

				for _, opt := range optionData {

					option := opt.([]interface{})

					options = append(
						options,
						option[0].(string),
					)
				}
			}
		}

		firstEntry, ok := entryData[0].([]any)
		if !ok || len(firstEntry) == 0 {
			fmt.Println("Entry ID: None")
			continue
		}

		entryID := firstEntry[0]

		field := FormField{
			Question: title,
			EntryID:  strconv.FormatInt(int64(entryID.(float64)), 10),
			Type:     questionType,
			Options:  options,
		}
		fields = append(fields, field)
	}

	jsonData, err := json.MarshalIndent(fields, "", "  ")

	if err != nil {
		panic(err)
	}

	os.WriteFile(
		"form_fields.json",
		jsonData,
		0644,
	)

	fmt.Println("form_fields.json created")
}
