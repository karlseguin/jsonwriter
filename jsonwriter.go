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
	null         = []byte("null")
	_true        = []byte("true")
	_false       = []byte("false")
	comma        = []byte(",")
	keyEnd       = []byte(`":`)
	startObject  = []byte("{")
	endObject    = []byte("}")
	startArray   = []byte("[")
	endArray     = []byte("]")
	escapedQuote = []byte(`\"`)
	escapedSlash = []byte(`\\`)
	escapedBS    = []byte(`\b`)
	escapedFF    = []byte(`\f`)
	escapedNL    = []byte(`\n`)
	escapedLF    = []byte(`\r`)
	escapedTab   = []byte(`\t`)
)

type Writer struct {
	depth int
	first bool
	array bool
	W     io.Writer
	end   []byte
}

// Creates a JsonWriter that writes to the provided io.Writer
func New(w io.Writer) *Writer {
	return &Writer{
		W:     w,
		first: true,
	}
}

func (w *Writer) RootObject(f func()) {
	w.ORoot()
	f()
	w.ERoot()
}

func (w *Writer) RootArray(f func()) {
	w.ARoot()
	f()
	w.ERoot()
}

func (w *Writer) Object(key string, f func()) {
	w.SObject(key)
	f()
	w.EObject()
}

func (w *Writer) Array(key string, f func()) {
	w.SArray(key)
	f()
	w.EArray()
}

// Starts a root object
func (w *Writer) ORoot() {
	w.end = endObject
	w.W.Write(startObject)
}

// Starts an array object
func (w *Writer) ARoot() {
	w.end = endArray
	w.array = true
	w.W.Write(startArray)
}

// Ends the root (used for both ORoot and ARoots)
func (w *Writer) ERoot() {
	w.array = false
	w.W.Write(w.end)
}

// Starts an object with the specified key
func (w *Writer) SObject(key string) {
	w.Key(key)
	w.first = true
	w.W.Write(startObject)
}

// Ends the object
func (w *Writer) EObject() {
	w.W.Write(endObject)
}

// Starts an array with the specified key
func (w *Writer) SArray(key string) {
	w.Key(key)
	w.first, w.array = true, true
	w.W.Write(startArray)
}

// Ends the array
func (w *Writer) EArray() {
	w.array = false
	w.W.Write(endArray)
}

// Writes a key. The key is placed within quotes and ends
// with a colon
func (w *Writer) Key(key string) {
	w.Separator()
	w.W.Write(keyStart)
	w.writeString(key)
	w.W.Write(keyEnd)
}

// value can be a string, byte, u?int(8|16|32|64)?, float(32|64)?,
// time.Time, bool or nil
func (w *Writer) Value(value interface{}) {
	if w.array {
		w.Separator()
	}

	if value == nil {
		w.W.Write(null)
		return
	}

	switch t := value.(type) {
	case bool:
		if t == true {
			w.W.Write(_true)
		} else {
			w.W.Write(_false)
		}
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

func (w *Writer) Separator() {
	if w.first == false {
		w.W.Write(comma)
	} else {
		w.first = false
	}
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
