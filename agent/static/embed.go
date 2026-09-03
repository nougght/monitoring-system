package static

import "embed"

// local web ui
//
//go:embed files/*
var StaticFiles embed.FS
