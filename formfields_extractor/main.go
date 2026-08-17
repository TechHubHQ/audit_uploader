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

	html := string(data)

	const marker = "var FB_PUBLIC_LOAD_DATA_ = "

	start := strings.Index(html, marker)
	if start == -1 {
		panic("form data not found")
	}

	start += len(marker)

	end := strings.Index(html[start:], ";")
	if end == -1 {
		panic("end of form data not found")
	}

	raw := html[start : start+end]

	var formData []any

	err = json.Unmarshal([]byte(raw), &formData)
	if err != nil {
		panic(err)
	}

	// Navigate to questions.
	level1 := formData[1].([]any)
	level2 := level1[1].([]any)
	questions := level2

	fmt.Println("Total elements:", len(questions))

	fields := make([]FormField, 0)

	for _, q := range questions {

		// Each question is itself an []any.
		question, ok := q.([]any)
		if !ok {
			continue
		}

		// Need at least question ID, title, type, etc.
		if len(question) < 5 {
			continue
		}

		// Question title.
		title, ok := question[1].(string)
		if !ok {
			continue
		}

		// Question type.
		questionType, ok := question[3].(float64)
		if !ok {
			continue
		}

		// Type 8 is a section/header, not an actual form field.
		if int(questionType) == 8 {
			continue
		}

		// Extract entry information.
		entryData, ok := question[4].([]any)
		if !ok || len(entryData) == 0 {
			continue
		}

		firstEntry, ok := entryData[0].([]any)
		if !ok || len(firstEntry) == 0 {
			continue
		}

		// Entry ID.
		entryIDFloat, ok := firstEntry[0].(float64)
		if !ok {
			continue
		}

		entryID := strconv.FormatInt(int64(entryIDFloat), 10)

		field := FormField{
			Question: title,
			EntryID:  entryID,
			Type:     int(questionType),
			Options:  getOptions(question),
		}

		fmt.Println("--------------------")
		fmt.Println("Question:", field.Question)
		fmt.Println("Type:", field.Type)
		fmt.Println("Entry ID:", field.EntryID)

		if len(field.Options) > 0 {
			fmt.Println("Options:", field.Options)
		}

		fields = append(fields, field)
	}

	jsonData, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("form_fields.json", jsonData, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"\nform_fields.json created with %d fields\n",
		len(fields),
	)
}
