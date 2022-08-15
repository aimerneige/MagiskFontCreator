package internal

import (
	_ "embed" // import embedded assets
)

//go:embed html/index.html
var aboutHTML string
