package crawler

import (
	"embed"
	"html/template"
	"time"
)

type statusType string

const (
	statusInfo    statusType = "info"
	statusSuccess statusType = "success"
	statusError   statusType = "error"
)

type summaryKey struct {
	Status statusType
	Value  string
}

type summaryType struct {
	Keys []summaryKey
	Vals []int
}

type linkType struct {
	Code   int
	Status statusType
	Title  string
	Href   string
	Refs   []string
}

type templateData struct {
	Subject   string
	Href      string
	CreatedAt time.Time
	Summary   summaryType
	Errors    []linkType
	Externals []linkType
}

//go:embed "templates"
var templatesFS embed.FS
var tmpl = template.Must(template.ParseFS(templatesFS, "templates/*.*"))
