package user

import (
	"embed"
)

//go:embed all:template
var TemplateFS embed.FS
