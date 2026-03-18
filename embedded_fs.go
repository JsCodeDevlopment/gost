package gost

import "embed"

//go:embed src/* locales/* docs/* main.go.tpl go.mod .env.example docker-compose.yml
var TemplateFS embed.FS
