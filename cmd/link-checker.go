package main

import (
	"fmt"
	"github.com/ofstudio/link-checker/config"
	"github.com/ofstudio/link-checker/crawler"
	"os"
	"time"
)

const reportsDir = "reports/"

func main() {
	p := config.MustGetProfile("config.yaml")

	err := os.MkdirAll(reportsDir, 0755)
	if err != nil {
		handleError(err)
	}

	fmt.Printf("Checking %s: %s => %s\n\n", p.Id, p.Title, p.HomePage)

	report, err := crawler.Run(p.Title, p.HomePage, p.Delay)
	if err != nil {
		handleError(err)
	}

	html, err := report.Html()
	if err != nil {
		handleError(err)
	}

	err = os.WriteFile(reportsDir+reportFilename(p.Id), html, 0644)
	if err != nil {
		handleError(err)
	}
}

func reportFilename(id string) string {
	return id + "-report-" + time.Now().Format("2006-01-02") + ".html"
}

func handleError(err error) {
	fmt.Println("Fatal error: ", err)
	os.Exit(-1)
}
