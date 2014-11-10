package jsonwriter

import (
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"
)

var (
	quote        = []byte(`"`)
	keyStart     = quote
	keyEnd       = []byte(`":`)
	escapedQuote = []byte(`\"`)
	escapedSlash = []byte(`\\`)
	escapedBS    = []byte(`\b`)
	escapedFF    = []byte(`\f`)
	escapedNL    = []byte(`\n`)
	escapedLF    = []byte(`\r`)
	escapedTab   = []byte(`\t`)
)

type Writer struct {
	W io.Writer
}

// Creates a JsonWriter that writes to the provided io.Writer
func New(w io.Writer) *Writer {
	return &Writer{w}
}

// Writes a key. The key is placed within quotes and ends
// with a colon
func (w *Writer) Key(key string) {
	w.W.Write(keyStart)
	w.writeString(key)
	w.W.Write(keyEnd)
}

// value can be a string, byte, u?int(8|16|32|64)?, float(32|64)?,
// time.Time, bool or nil
func (w *Writer) Value(value interface{}) {
	switch t := value.(type) {
	case uint8:
		w.W.Write([]byte(strconv.FormatUint(uint64(t), 10)))
	case uint16:
		w.W.Write([]byte(strconv.FormatUint(uint64(t), 10)))
	case uint32:
		w.W.Write([]byte(strconv.FormatUint(uint64(t), 10)))
	case uint:
		w.W.Write([]byte(strconv.FormatUint(uint64(t), 10)))
	case uint64:
		w.W.Write([]byte(strconv.FormatUint(t, 10)))
	case int8:
		w.W.Write([]byte(strconv.FormatInt(int64(t), 10)))
	case int16:
		w.W.Write([]byte(strconv.FormatInt(int64(t), 10)))
	case int32:
		w.W.Write([]byte(strconv.FormatInt(int64(t), 10)))
	case int:
		w.W.Write([]byte(strconv.FormatInt(int64(t), 10)))
	case int64:
		w.W.Write([]byte(strconv.FormatInt(t, 10)))
	case float32:
		w.W.Write([]byte(strconv.FormatFloat(float64(t), 'g', -1, 32)))
	case float64:
		w.W.Write([]byte(strconv.FormatFloat(t, 'g', -1, 64)))
	case string:
		w.W.Write(quote)
		w.writeString(t)
		w.W.Write(quote)
	default:
		panic(fmt.Sprintf("unsuported valued type %v", value))
	}
}

// writes a key: value
// This is the same as calling WriteKey(key) followe by WriteValue(value)
func (w *Writer) KeyValue(key string, value interface{}) {
	w.Key(key)
	w.Value(value)
}

func (w *Writer) writeString(s string) {
	start, end := 0, 0
	var special []byte
L:
	for i, r := range s {
		switch r {
		case '"':
			special = escapedQuote
		case '\\':
			special = escapedSlash
		case '\b':
			special = escapedBS
		case '\f':
			special = escapedFF
		case '\n':
			special = escapedNL
		case '\r':
			special = escapedLF
		case '\t':
			special = escapedTab
		default:
			end += utf8.RuneLen(r)
			continue L
		}

		if end > start {
			w.W.Write([]byte(s[start:end]))
		}
		w.W.Write(special)
		start = i + 1
		end = start
	}
	if end > start {
		w.W.Write([]byte(s[start:end]))
	}
}
