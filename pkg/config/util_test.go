package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseKVArgs(t *testing.T) {
	t.Run("empty-state", func(t *testing.T) {
		s := map[string]any{}
		kv := []string{
			"foo=bar",
			"bar=bazz=fuzz",
			"creds.token=basic 1237934hsdf98",
		}
		exp := map[string]any{
			"foo": "bar",
			"bar": "bazz=fuzz",
			"creds": map[string]any{
				"token": "basic 1237934hsdf98",
			},
		}

		ParseKVArgs(kv, s)
		assert.Equal(t, exp, s)
	})

	t.Run("prefilled-state", func(t *testing.T) {
		s := map[string]any{
			"hello": "world",
			"headers": map[string]any{
				"Content-Type": "application/json",
			},
		}
		kv := []string{
			"foo=bar",
			"bar=bazz=fuzz",
			"creds.token=basic 1237934hsdf98",
			"headers.User-Agent=test v123",
		}
		exp := map[string]any{
			"hello": "world",
			"foo":   "bar",
			"bar":   "bazz=fuzz",
			"creds": map[string]any{
				"token": "basic 1237934hsdf98",
			},
			"headers": map[string]any{
				"Content-Type": "application/json",
				"User-Agent":   "test v123",
			},
		}

		ParseKVArgs(kv, s)
		assert.Equal(t, exp, s)
	})
}
