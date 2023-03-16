package engine

import (
	"fmt"
	"strings"

	"github.com/zekrotja/rogu/log"
)

func (t *Goja) builtin_assert(v bool, msg ...string) {
	mesg := "assertion failed"
	if len(msg) != 0 {
		mesg = fmt.Sprintf("%s: %s", mesg, strings.Join(msg, " "))
	}

	if !v {
		panic(t.rt.ToValue(mesg))
	}
}

func (t *Goja) builtin_debug(msg string) {
	log.Debug().Msg(msg)
}

func (t *Goja) builtin_info(msg string) {
	log.Info().Msg(msg)
}

func (t *Goja) builtin_warn(msg string) {
	log.Warn().Msg(msg)
}

func (t *Goja) builtin_error(msg string) {
	log.Error().Msg(msg)
}

func (t *Goja) builtin_fatal(msg string) {
	log.Fatal().Msg(msg)
}

func (t *Goja) builtin_print(msg string) {
	fmt.Print(msg)
}

func (t *Goja) builtin_println(msg string) {
	fmt.Println(msg)
}
