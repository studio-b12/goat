// Package clr allows to format strings using
// shell format codes conditionally, so that
// it can be easily used with any logger.
//
// Formatting can simply be disabled on the
// Colorizer
package clr

import (
	"fmt"
	"strings"
)

// FormatCode describes a text format.
type FormatCode int

const (
	FormatReset     = FormatCode(0)
	FormatBold      = FormatCode(1)
	FormatUnderline = FormatCode(4)

	ColorFGBlack  = FormatCode(30)
	ColorFGRed    = FormatCode(31)
	ColorFGGreen  = FormatCode(32)
	ColorFGBrown  = FormatCode(33)
	ColorFGBlue   = FormatCode(34)
	ColorFGPurple = FormatCode(35)
	ColorFGCyan   = FormatCode(36)
	ColorFGGray   = FormatCode(37)

	ColorBGBlack  = FormatCode(40)
	ColorBGRed    = FormatCode(41)
	ColorBGGreen  = FormatCode(42)
	ColorBGBrown  = FormatCode(43)
	ColorBGBlue   = FormatCode(44)
	ColorBGPurple = FormatCode(45)
	ColorBGCyan   = FormatCode(46)
	ColorBGGray   = FormatCode(47)
)

type formatWrapper struct {
	value any
	codes []FormatCode
}

// Format takes the given value and wraps
// it to a colorWrapper with the given
// format codes. The colorWrapper can then
// be passed into the Print method of a
// Printer.
func Format(v any, codes ...FormatCode) formatWrapper {
	return formatWrapper{value: v, codes: codes}
}

// Printer is used to print raw and format
// wrapped values.
type Printer struct {
	enabled bool
}

// SetEnable specifies whether or not the
// printer should format given values.
func (c *Printer) SetEnable(v bool) {
	c.enabled = v
}

// Print takes one or more values and prints
// then to a string.
//
// When a formatWrapper is passed, the format
// will be applied depending on the Printers
// enabled state.
func (c *Printer) Print(values ...any) string {
	var sb strings.Builder

	for _, val := range values {
		if wrapper, ok := val.(formatWrapper); ok {
			if c.enabled && len(wrapper.codes) > 0 {
				fmt.Fprint(&sb, "\x1b[")
				for i, code := range wrapper.codes {
					fmt.Fprintf(&sb, "%d", code)
					if i < len(wrapper.codes)-1 {
						fmt.Fprint(&sb, ";")
					}
				}
				fmt.Fprintf(&sb, "m%v\x1b[0m", wrapper.value)
				continue
			}

			val = wrapper.value
		}

		fmt.Fprintf(&sb, "%v", val)
	}

	return sb.String()
}
