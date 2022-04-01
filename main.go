package main

import (
	"embed"

	"github.com/joernott/eyas/cmd"
	"github.com/joernott/eyas/server"
)

// It will add the specified files.
//go:embed index.html
// It will add all non-hidden file in images, css, and js.
//go:embed images css js

var Static embed.FS

func main() {
	server.Static = Static
	cmd.Execute()
}
