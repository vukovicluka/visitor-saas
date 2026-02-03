package web

import (
	"embed"
	_ "embed"
)

//go:embed tracker/tracker.js
var TrackerJS []byte

//go:embed static
var StaticFS embed.FS