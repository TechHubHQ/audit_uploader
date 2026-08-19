package main

import (
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func launchBrowser() {
	l := launcher.New().
		Headless(false).
		MustLaunch()

	browser := rod.New().
		ControlURL(l).
		MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://forms.gle/jeCB1AuEmiW1Uzr39")
	page.MustWindowMaximize()
	page.MustWaitLoad()

	time.Sleep(5 * time.Second)
}
