package executor

import (
	"strings"

	"github.com/studio-b12/goat/pkg/clr"
	"github.com/zekrotja/rogu/log"
)

func printSeparator(head string) {
	const lenSpacerTotal = 100

	lenSpacer := lenSpacerTotal - 2 - len(head)
	lenSpacerLeft := lenSpacer / 2
	lenSpacerRight := lenSpacerLeft
	if lenSpacer%2 > 0 {
		lenSpacerRight++
	}

	msg := clr.Print(clr.Format("%s %s %s", clr.ColorFGPurple))
	log.Info().Msgf(msg,
		strings.Repeat("-", lenSpacerLeft),
		head,
		strings.Repeat("-", lenSpacerRight))
}
