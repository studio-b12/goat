package embedded

import (
	_ "embed"
)

var (
	//go:embed embeds/new_goatfile.goat
	NewGoatfile []byte
)
