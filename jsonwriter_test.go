package jsonwriter

import (
	"bytes"
	. "github.com/karlseguin/expect"
	"testing"
	"time"
)

type WriterTests struct{}

func Test_Writer(t *testing.T) {
	Expectify(new(WriterTests), t)
}

func (_ *WriterTests) WritesAnInt() {
	assertValue(uint8(1), "1")
	assertValue(uint16(2), "2")
	assertValue(uint32(232134), "232134")
	assertValue(uint64(432434), "432434")
	assertValue(uint(5), "5")
	assertValue(int8(-3), "-3")
	assertValue(int16(-16), "-16")
	assertValue(int32(-31), "-31")
	assertValue(int64(-4343), "-4343")
	assertValue(int(-59922), "-59922")
}

func (_ *WriterTests) WritesAFloat() {
	assertValue(float32(1.2393), "1.2393")
	assertValue(float32(-49493.443), "-49493.44")
	assertValue(float64(99499449.23949), "9.949944923949e+07")
	assertValue(float64(-3290123.94994), "-3.29012394994e+06")
}

func (_ *WriterTests) WritesAString() {
	assertValue(`abc`, `"abc"`)
	assertValue(`ab"cd`, `"ab\"cd"`)
	assertValue(`💣`, `"💣"`)
	assertValue("\\it's\n\tOver\r9000!\\ 💣 💣 💣", `"\\it's\n\tOver\r9000!\\ 💣 💣 💣"`)
}

func (_ *WriterTests) WritesABool() {
	assertValue(true, "true")
	assertValue(false, "false")
}

func (_ *WriterTests) WritesANull() {
	assertValue(nil, "null")
}

func (_ *WriterTests) WritesATime() {
	assertValue(time.Unix(1415677601, 9), `"2014-11-11T10:46:41.000000009+07:00"`)
}

func (_ *WriterTests) SimpleObject() {
	w, b := n()
	w.ORoot()
	w.KeyValue("spice", "flow")
	w.ERoot()
	Expect(b.String()).To.Equal(`{"spice":"flow"}`)
}

func (_ *WriterTests) MultiValueObject() {
	w, b := n()
	w.ORoot()
	w.KeyValue("spice", "flow")
	w.KeyValue("over", 9000)
	w.ERoot()
	Expect(b.String()).To.Equal(`{"spice":"flow","over":9000}`)
}

func (_ *WriterTests) NestedObject1() {
	w, b := n()
	w.ORoot()
	w.KeyValue("spice", "flow")
	w.KeyValue("over", 9000)

	w.SObject("first")
	w.KeyValue("afraid", true)
	w.EObject()

	w.SObject("second")
	w.KeyValue("a", byte(1))
	w.KeyValue("b", 1.01)
	w.EObject()

	w.ERoot()
	Expect(b.String()).To.Equal(JSON(`{
		"spice":"flow",
		"over":9000,
		"first": {"afraid": true},
		"second": {"a": 1, "b": 1.01}
	}`))
}

func (_ *WriterTests) NestedObject2() {
	w, b := n()
	w.ORoot()
	w.SObject("first")
	w.SObject("second")
	w.SObject("third")
	w.EObject()
	w.EObject()
	w.EObject()
	w.ERoot()
	Expect(b.String()).To.Equal(JSON(`{
		"first":{
			"second":{
				"third":{}
			}
		}
	}`))
}

func (_ *WriterTests) NestedObject3() {
	w, b := n()
	w.ORoot()
	w.SObject("first")
	w.SObject("second")
	w.KeyValue("a", true)
	w.SObject("third")
	w.EObject()
	w.EObject()
	w.EObject()
	w.ERoot()
	Expect(b.String()).To.Equal(JSON(`{
		"first":{
			"second":{
				"a": true,
				"third":{}
			}
		}
	}`))
}

func (_ *WriterTests) RootArray() {
	w, b := n()
	w.ARoot()
	w.Value(1)
	w.Value("b\"")
	w.Value(nil)
	w.ERoot()
	Expect(b.String()).To.Equal(`[1,"b\"",null]`)
}

func (_ *WriterTests) NestedArray() {
	w, b := n()
	w.ORoot()
	w.SArray("scores")
	w.Value(3)
	w.Value(5)
	w.EArray()
	w.ERoot()
	Expect(b.String()).To.Equal(JSON(`{
		"scores": [3, 5]
	}`))
}

func (_ *WriterTests) AlternativeSyntaxObject() {
	w, b := n()

	w.RootObject(func() {
		w.KeyValue("power", 9000)
		w.Array("scores", func() {
			w.Value(1)
			w.Value(true)
			w.Value(nil)
		})
		w.Object("atreides", func() {
			w.KeyValue("name", "leto")
			w.KeyValue("sister", "ghanima")
		})
	})

	Expect(b.String()).To.Equal(JSON(`{
		"power": 9000,
		"scores":[1, true, null],
		"atreides": {
			"name": "leto",
			"sister": "ghanima"
		}
	}`))
}

func (_ *WriterTests) AlternativeSyntaxArray() {
	w, b := n()

	w.RootArray(func() {
		w.Value(1.2)
		w.Value(false)
		w.Value("\n")
	})

	Expect(b.String()).To.Equal(JSON(`[1.2, false, "\n"]`))
}


func assertValue(value interface{}, expected string) {
	w, b := n()
	w.Value(value)
	Expect(b.String()).To.Equal(expected)
}

func n() (*Writer, *bytes.Buffer) {
	b := new(bytes.Buffer)
	return New(b), b
}
