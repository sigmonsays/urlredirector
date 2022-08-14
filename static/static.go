package static

import (
	"embed"
	_ "embed"
)

//go:embed *
var Files embed.FS
