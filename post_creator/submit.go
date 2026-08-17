package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func Submit(values url.Values) error {
	const formURL = "https://docs.google.com/forms/d/e/1FAIpQLScuaiAFmILMbBua-wxXuQqh-3_uJAg_bxmSbFNix8kiw5LiGQ/formResponse"

	resp, err := http.Post(
		formURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Google Forms returned %s", resp.Status)
	}

	return nil
}
