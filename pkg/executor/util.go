package executor

import (
	"encoding/json"
	"encoding/xml"
	"strings"

	"github.com/studio-b12/goat/pkg/clr"
	"github.com/studio-b12/goat/pkg/errs"
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

func parseBody(data []byte, responseType string) (any, error) {
	if responseType == "json" || strings.Contains(responseType, "application/json") {
		var bodyJson any
		err := json.Unmarshal(data, &bodyJson)
		if err != nil {
			return nil, errs.WithPrefix("failed unmarshalling json:", err)
		}
		return bodyJson, nil
	} else if responseType == "xml" || strings.Contains(responseType, "text/xml") {
		var bodyXml any
		err := xml.Unmarshal(data, &bodyXml)
		if err != nil {
			return nil, errs.WithPrefix("failed unmarshalling xml:", err)
		}
		return bodyXml, nil
	} else {
		// In the default case, return data as string
		return string(data), nil
	}
}
