// Package static contains static elements and embeds those in variable. Makes testing easier
package static

import (
	"embed"
	"html/template"
)

//go:embed *.html
var templateFS embed.FS

// Templates contains all html templates inside folder package
var Templates = template.Must(template.ParseFS(templateFS, "*.html"))
