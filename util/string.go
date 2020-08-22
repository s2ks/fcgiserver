package util

import (
	//"fmt"
	"strings"
)

func ByteSubstituteMap(raw []byte, keyval map[string]string, delim string) ([]byte, error) {

	dest := make([]byte, len(raw))
	copy(dest, raw)

	for key, val := range keyval {
		pattern := make([]byte, len(key)+(2*len(delim)))
		w := copy(pattern[0:], []byte(delim))
		w += copy(pattern[w:], []byte(key))
		w += copy(pattern[w:], []byte(delim))

		keycount := strings.Count(string(dest), string(pattern))

		if keycount < 1 {
			continue
		}

		var buf []byte

		/* if val is longer than pattern we need a buffer big enough to fit the final result */
		if len(val) > w {
			buf = make([]byte, len(dest)+(keycount*(len(val)-w)))
		} else {
		/* otherwise just use dest*/
			buf = dest
		}

		for i := 0; i < keycount; i++ {
			keystart := strings.Index(string(dest), string(pattern))
			keyend := keystart + len(pattern)

			start := dest[0:keystart]
			end := dest[keyend:]

			w = copy(buf[0:], start)
			w += copy(buf[w:], val)
			w += copy(buf[w:], end)

			dest = make([]byte, w)

			copy(dest[0:], buf[0:w])
		}
	}

	return dest, nil
}

func SubstituteMap(str string, keyval map[string]string, delim string) (string, error) {
	raw := []byte(str)

	buf, err := ByteSubstituteMap(raw, keyval, delim)

	return string(buf), err
}

func Substitute(str, key, val string, delim string) (string, error) {
	raw := []byte(str)
	varmap := make(map[string]string)

	varmap[key] = val

	buf, err := ByteSubstituteMap(raw, varmap, delim)

	return string(buf), err
}

/*
func Substitute(str, key, val string, delim byte) (string, error) {
	raw := []byte(str)
	rawlen := len(raw)

	dest := make([]byte, rawlen + (len(val) - len(key)))
	written := 0

	keystart := -1
	keyend := -1

	for i := 0; i < rawlen; i++ {
		if raw[i] == delim && keystart < 0 {
			keystart = i
		} else if raw[i] == delim {
			keyend = i
			break
		}
	}

	if keystart < 0 {
		return str, fmt.Errorf("No variable found")
	} else if keyend < 0 {
		return str, fmt.Errorf("Syntax error at: \"%v\"", str)
	}

	start := raw[0:keystart]
	end := raw[keyend + 1:]
	keyname := raw[keystart + 1:keyend]
	keyval := []byte(val)

	if key == string(keyname) {
		written += copy(dest[written:], start)
		written += copy(dest[written:], keyval)
		written += copy(dest[written:], end)
		return string(dest[0:written]), nil
	} else {
		return str, fmt.Errorf("Unknown variable \"%%%v%%\"", keyname)
	}
}*/
