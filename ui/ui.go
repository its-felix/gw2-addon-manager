package ui

import (
	"embed"
	"net/http"
)

//go:embed dist/*
var uiFS embed.FS

func FS() http.FileSystem {
	return http.FS(uiFS)
}
