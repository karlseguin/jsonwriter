package jsonwriter

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"
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
}

// Creates a JsonWriter that writes to the provided io.Writer
func New(w io.Writer) *Writer {
	return &Writer{
		W:     w,
		first: true,
	}
}

// Starts the writing process by creating an object.
// Should only be called once
func (w *Writer) RootObject(f func()) {
	w.W.Write(startObject)
	f()
	w.W.Write(endObject)
}

// Starts the writing process by creating an array.
// Should only be called once
func (w *Writer) RootArray(f func()) {
	w.Separator()
	w.array = true
	w.first = true
	w.W.Write(startArray)
	f()
	w.W.Write(endArray)
}

// Star an object with the given key
func (w *Writer) Object(key string, f func()) {
	w.Key(key)
	w.first = true
	w.W.Write(startObject)
	f()
	w.first = false
	w.W.Write(endObject)
}

// Star an array with the given key
func (w *Writer) Array(key string, f func()) {
	w.Key(key)
	w.first, w.array = true, true
	w.W.Write(startArray)
	f()
	w.first, w.array = false, false
	w.W.Write(endArray)
}

// Star an object within an array (a keyless object)
func (w *Writer) ArrayObject(f func()) {
	w.Separator()
	w.first = true
	w.array = false
	w.W.Write(startObject)
	f()
	w.W.Write(endObject)
	w.array = true
}

// Writes an array with the given array of values
// (values must be an array or a slice)
func (w *Writer) ArrayValues(key string, values interface{}) {
	w.Key(key)
	w.first, w.array = true, true
	w.W.Write(startArray)

	vo := reflect.ValueOf(values)
	kind := vo.Kind()
	if kind == reflect.Array || kind == reflect.Slice {
		for i, l := 0, vo.Len(); i < l; i++ {
			w.Value(vo.Index(i).Interface())
		}
	}
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
// time.Time, bool, nil or encoding/json.Marshaller
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
		w.Int(int(t))
	case uint16:
		w.Int(int(t))
	case uint32:
		w.Int(int(t))
	case uint:
		w.W.Write([]byte(strconv.FormatUint(uint64(t), 10)))
	case uint64:
		w.W.Write([]byte(strconv.FormatUint(t, 10)))
	case int8:
		w.Int(int(t))
	case int16:
		w.Int(int(t))
	case int32:
		w.Int(int(t))
	case int:
		w.Int(int(t))
	case int64:
		w.W.Write([]byte(strconv.FormatInt(t, 10)))
	case float32:
		w.W.Write([]byte(strconv.FormatFloat(float64(t), 'g', -1, 32)))
	case float64:
		w.W.Write([]byte(strconv.FormatFloat(t, 'g', -1, 64)))
	case json.Marshaler:
		b, _ := t.MarshalJSON()
		w.W.Write(b)
	case string:
		w.String(t)
	case time.Time:
		w.W.Write(quote)
		w.writeString(t.Format(time.RFC3339Nano))
		w.W.Write(quote)
	default:
		panic(fmt.Sprintf("unsuported valued type %v", value))
	}
}

func (w *Writer) Int(value int) {
	w.W.Write([]byte(strconv.Itoa(value)))
}

func (w *Writer) String(value string) {
	w.W.Write(quote)
	w.writeString(value)
	w.W.Write(quote)
}

// Writes raw data to an array as-is, leaving encoding up to the caller
func (w *Writer) RawValue(data []byte) {
	if w.array {
		w.Separator()
	}

	w.W.Write(data)
}

// Writes the value as-is, leaving delimiters and encoding up to the caller
func (w *Writer) Raw(data []byte) {
	w.W.Write(data)
}

// writes a key: value
// This is the same as calling WriteKey(key) followe by WriteValue(value)
func (w *Writer) KeyValue(key string, value interface{}) {
	w.Key(key)
	w.Value(value)
}

// writes a key: value
// write a key'd raw value
func (w *Writer) KeyRaw(key string, data []byte) {
	w.Key(key)
	w.W.Write(data)
}

// writes a key: value where value is a string
// This is an optimized Write method
func (w *Writer) KeyString(key string, value string) {
	w.Key(key)
	w.W.Write(quote)
	w.writeString(value)
	w.W.Write(quote)
}

// writes a key: value where value is a int
// This is an optimized Write method
func (w *Writer) KeyInt(key string, value int) {
	w.Key(key)
	w.W.Write([]byte(strconv.Itoa(value)))
}

// writes a key: value where value is a bool
// This is an optimized Write method
func (w *Writer) KeyBool(key string, value bool) {
	w.Key(key)
	if value {
		w.W.Write(_true)
	} else {
		w.W.Write(_false)
	}
}

// writes a key: value where value is a bool
// This is an optimized Write method
func (w *Writer) KeyFloat(key string, value float32) {
	w.Key(key)
	w.W.Write([]byte(strconv.FormatFloat(float64(t), 'g', -1, 32)))
}

// writes a key: value where value is a bool
// This is an optimized Write method
func (w *Writer) KeyFloat64(key string, value float64) {
	w.Key(key)
	w.W.Write([]byte(strconv.FormatFloat(value, 'g', -1, 64)))
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

func (w *Writer) Reset() {
	w.depth = 0
	w.first = true
	w.array = false
}
