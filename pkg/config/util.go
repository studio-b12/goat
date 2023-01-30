package config

import "strings"

// ParseKVArgs takes a list of key value pairs which are then
// applied to the map m. Key value paris are split by '='. Sub
// levels in maps are separated by '.'.
//
// Example:
//
//	m := make(map[string]any)
//	kv := []string{"foo.bar=bazz"}
//	ParseKVArgs(kv, m) // m = { "foo": { "bar": "bazz" } }
func ParseKVArgs(args []string, m map[string]any) {
	for _, kvPair := range args {
		kv := strings.SplitN(kvPair, "=", 2)
		key, val := kv[0], kv[1]

		keyPath := strings.Split(key, ".")

		subm := m
		var ok bool
		for i := 0; i < len(keyPath)-1; i++ {
			subm, ok = m[keyPath[i]].(map[string]any)
			if !ok {
				subm = make(map[string]any)
				m[keyPath[i]] = subm
			}
		}

		subm[keyPath[len(keyPath)-1]] = val
	}
}
