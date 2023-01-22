package embedded

import (
	_ "embed"
)

var (
	//go:embed embeds/new_gurlfile.gurl
	NewGurlfile []byte
)
