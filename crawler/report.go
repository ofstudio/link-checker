package crawler

import (
	"bytes"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"sort"
	"strconv"
	"time"
)

type visitType struct {
	code       int
	isExternal bool
}

type byHref []linkType

func (a byHref) Len() int           { return len(a) }
func (a byHref) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byHref) Less(i, j int) bool { return a[i].Href < a[j].Href }

type Report struct {
	subject   string
	href      string
	createdAt time.Time
	visits    map[string]visitType // HTTP-коды посещенных URL
	referrers map[string][]string  // Ссылающиеся страницы
	aliases   map[string]string    // Алиасы урлов [реальный]исходный
}

func newReport(subject, href string) *Report {
	return &Report{
		subject:   subject,
		href:      href,
		createdAt: time.Now(),
		visits:    make(map[string]visitType),
		referrers: make(map[string][]string),
		aliases:   make(map[string]string),
	}
}

func (r *Report) visit(code int, url string, isExternal bool) {
	r.visits[url] = visitType{
		code:       code,
		isExternal: isExternal,
	}
}

func (r *Report) refers(url, who string) {
	for i := range r.referrers[url] {
		if r.referrers[url][i] == who {
			return
		}
	}
	r.referrers[url] = append(r.referrers[url], who)
}

func (r *Report) alias(real, original string) {
	r.aliases[real] = original
}

func (r *Report) Html() ([]byte, error) {

	d := templateData{
		Subject:   r.subject,
		Href:      r.href,
		CreatedAt: r.createdAt,
		Summary:   summaryType{},
	}

	// Summary
	codes := make(map[int]int)
	for _, visit := range r.visits {
		codes[visit.code] += 1
	}
	d.Summary.Keys = append(d.Summary.Keys, summaryKey{Status: statusInfo, Value: "Total URL"})
	d.Summary.Vals = append(d.Summary.Vals, len(r.visits))
	for code, num := range codes {
		var status = getStatusByCode(code)
		d.Summary.Keys = append(d.Summary.Keys, summaryKey{
			Status: status,
			Value:  strconv.Itoa(code),
		})
		d.Summary.Vals = append(d.Summary.Vals, num)
	}

	// Errors and Externals
	for href, visit := range r.visits {
		refs := r.referrers[r.aliases[href]]
		e := linkType{
			Code:   visit.code,
			Status: getStatusByCode(visit.code),
			Title:  getTitle(href),
			Href:   href,
		}
		for i := range refs {
			e.Refs = append(e.Refs, refs[i])
		}

		// Externals
		if visit.isExternal {
			d.Externals = append(d.Externals, e)
		}

		// Errors
		if visit.code < 100 || visit.code >= 400 {
			d.Errors = append(d.Errors, e)
		}
	}

	sort.Sort(byHref(d.Errors))
	sort.Sort(byHref(d.Externals))

	// Render template
	b := bytes.Buffer{}
	err := tmpl.ExecuteTemplate(&b, "report.html", d)
	if err != nil {
		return nil, err
	}
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: true,
		KeepDocumentTags:    true,
		KeepEndTags:         true,
		KeepQuotes:          true,
		KeepWhitespace:      true,
	})
	return m.Bytes("text/html", b.Bytes())
}

// trim long urls
func getTitle(s string) string {
	if len(s) > 70 {
		return s[0:70] + "..."
	}
	return s
}

func getStatusByCode(code int) statusType {
	if code >= 100 && code < 400 {
		return statusSuccess
	}
	return statusError
}
