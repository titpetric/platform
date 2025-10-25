package user

import (
	"embed"
)

//go:embed all:templates
var TemplateFS embed.FS
