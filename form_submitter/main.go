package main

import (
	"net/http"
	"net/url"
	"strings"
)

func main() {
	data := url.Values{}

	data.Set(
		"entry.1388635514",
		"TEST123",
	)
	data.Set(
		"entry.1551094033",
		"ICU",
	)
	data.Set(
		"entry.1784678943",
		"YES",
	)

	http.Post(
		"https://docs.google.com/forms/d/e/FORM_ID/formResponse",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
}
