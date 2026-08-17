package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const formURL = "https://docs.google.com/forms/d/e/1FAIpQLScuaiAFmILMbBua-wxXuQqh-3_uJAg_bxmSbFNix8kiw5LiGQ/formResponse"

func main() {
	data := url.Values{}

	// Deliberately obvious test values.
	data.Set("entry.1388635514", "AUTOMATION-TEST")
	data.Set("entry.1551094033", "TEST-LOCATION")
	data.Set("entry.16127784", "TEST-DOCTOR")
	data.Set("entry.2039185993", "TEST-DEPARTMENT")

	// Choice field.
	data.Set("entry.1784678943", "Yes")

	fmt.Println("Submitting test response...")

	resp, err := http.Post(
		formURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("HTTP status:", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response body length:", len(body))
}
