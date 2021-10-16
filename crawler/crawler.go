package crawler

import "C"
import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"net/url"
	"strings"
	"time"
)

var ua = "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) " +
	"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36"

func Run(title, href string, delay time.Duration) (*Report, error) {
	u, err := url.Parse(href)
	if err != nil {
		return nil, err
	}

	collector := colly.NewCollector(
		colly.UserAgent(ua),
		colly.IgnoreRobotsTxt(),
	)
	collector.SetRequestTimeout(time.Second * 60)
	err = collector.Limit(&colly.LimitRule{Delay: delay, DomainGlob: "*"})
	if err != nil {
		return nil, err
	}

	report := newReport(title, u.String())

	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		switch {
		case e.Request.URL.Host != u.Host: // Don't scan external resource
			return
		case strings.HasPrefix(link, "mailto:"):
			return
		case strings.HasPrefix(link, "tel:"):
			return
		case strings.HasPrefix(link, "javascript:"):
			return
		}
		report.refers(e.Request.AbsoluteURL(link), e.Request.URL.String())
		ctx := colly.NewContext()
		ctx.Put("orgUrl", e.Request.AbsoluteURL(link))
		_ = collector.Request("GET", e.Request.AbsoluteURL(link), nil, ctx, nil)
	})

	collector.OnResponse(func(r *colly.Response) {
		if orgUrl := r.Ctx.Get("orgUrl"); orgUrl != "" {
			report.alias(r.Request.URL.String(), orgUrl)
		}
		report.visit(r.StatusCode, r.Request.URL.String(), r.Request.URL.Host != u.Host)
		fmt.Printf("%d - %s\n", r.StatusCode, r.Request.URL)
	})

	collector.OnError(func(r *colly.Response, err error) {
		if orgUrl := r.Ctx.Get("orgUrl"); orgUrl != "" {
			report.alias(r.Request.URL.String(), orgUrl)
		}
		report.visit(r.StatusCode, r.Request.URL.String(), r.Request.URL.Host != u.Host)
		fmt.Printf("%d - %s\n", r.StatusCode, r.Request.URL)
	})

	err = collector.Visit(u.String())
	if err != nil {
		return nil, err
	}
	return report, nil
}
