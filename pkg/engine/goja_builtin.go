package engine

import (
	"fmt"
	"github.com/itchyny/gojq"
	"reflect"
	"strings"

	"github.com/zekrotja/rogu/log"
)

func (t *Goja) builtin_assert(v bool, msg ...string) {
	if v {
		return
	}

	mesg := "assertion failed"
	if len(msg) != 0 {
		mesg = fmt.Sprintf("%s: %s", mesg, strings.Join(msg, " "))
	}

	panic(t.rt.ToValue(mesg))
}

func (t *Goja) builtin_assert_eq(value any, expected any, msg ...string) {
	if reflect.DeepEqual(value, expected) {
		return
	}

	part := "unexpected value"
	if len(msg) != 0 {
		part = strings.Join(msg, " ")
	}

	mesg := fmt.Sprintf("assertion failed: %s: expected `%v` != received `%v`", part, expected, value)

	panic(t.rt.ToValue(mesg))
}

func (t *Goja) builtin_debug(msg ...string) {
	log.Debug().Msg(strings.Join(msg, " "))
}

func (t *Goja) builtin_info(msg ...string) {
	log.Info().Msg(strings.Join(msg, " "))
}

func (t *Goja) builtin_warn(msg ...string) {
	log.Warn().Msg(strings.Join(msg, " "))
}

func (t *Goja) builtin_error(msg ...string) {
	log.Error().Msg(strings.Join(msg, " "))
}

func (t *Goja) builtin_fatal(msg ...string) {
	log.Fatal().Msg(strings.Join(msg, " "))
}

func (t *Goja) builtin_print(msg ...string) {
	fmt.Print(strings.Join(msg, " "))
}

func (t *Goja) builtin_println(msg ...string) {
	fmt.Println(strings.Join(msg, " "))
}

func (t *Goja) builtin_debugf(format string, v ...any) {
	log.Debug().Msgf(format, v...)
}

func (t *Goja) builtin_infof(format string, v ...any) {
	log.Info().Msgf(format, v...)
}

func (t *Goja) builtin_warnf(format string, v ...any) {
	log.Warn().Msgf(format, v...)
}

func (t *Goja) builtin_errorf(format string, v ...any) {
	log.Error().Msgf(format, v...)
}

func (t *Goja) builtin_fatalf(format string, v ...any) {
	log.Fatal().Msgf(format, v...)
}

func (t *Goja) builtin_printf(format string, v ...any) {
	fmt.Printf(format, v...)
}

func (t *Goja) builtin_jq(object any, src string) []any {
	query, err := gojq.Parse(src)
	if err != nil {
		panic(t.rt.ToValue(fmt.Sprintf("command parsing failed: %s", err.Error())))
	}

	var results []any
	iter := query.Run(object)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
				break
			}
			panic(t.rt.ToValue(err.Error()))
		}
		results = append(results, v)
	}

	return results
}
