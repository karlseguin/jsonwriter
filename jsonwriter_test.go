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

func (_ WriterTests) WritesAnInt() {
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

func (_ WriterTests) WritesAFloat() {
	assertValue(float32(1.2393), "1.2393")
	assertValue(float32(-49493.443), "-49493.44")
	assertValue(float64(99499449.23949), "9.949944923949e+07")
	assertValue(float64(-3290123.94994), "-3.29012394994e+06")
}

func (_ WriterTests) WritesAString() {
	assertValue(`abc`, `"abc"`)
	assertValue(`ab"cd`, `"ab\"cd"`)
	assertValue(`💣`, `"💣"`)
	assertValue("\\it's\n\tOver\r9000!\\ 💣 💣 💣", `"\\it's\n\tOver\r9000!\\ 💣 💣 💣"`)
}

func (_ WriterTests) WritesABool() {
	assertValue(true, "true")
	assertValue(false, "false")
}

func (_ WriterTests) WritesANull() {
	assertValue(nil, "null")
}

func (_ WriterTests) WritesATime() {
	assertValue(time.Unix(1415677601, 9), `"2014-11-11T10:46:41.000000009+07:00"`)
}

func (_ WriterTests) SimpleObject() {
	w, b := n()
	w.RootObject(func() {
		w.KeyValue("spice", "flow")
	})
	Expect(b.String()).To.Equal(`{"spice":"flow"}`)
}

func (_ WriterTests) MultiValueObject() {
	w, b := n()
	w.RootObject(func() {
		w.KeyValue("spice", "flow")
		w.KeyValue("over", 9000)
	})
	Expect(b.String()).To.Equal(`{"spice":"flow","over":9000}`)
}

func (_ WriterTests) NestedObject1() {
	w, b := n()
	w.RootObject(func() {
		w.KeyValue("power", 9000)
		w.Object("atreides", func() {
			w.KeyValue("name", "leto")
			w.KeyValue("sister", "ghanima")
		})
	})

	Expect(b.String()).To.Equal(JSON(`{
		"power": 9000,
		"atreides": {
			"name": "leto",
			"sister": "ghanima"
		}
	}`))
}

func (_ WriterTests) NestedObject2() {
	w, b := n()
	w.RootObject(func() {
		w.KeyValue("power", 9000)
		w.Object("atreides", func() {
			w.KeyValue("name", "leto")
			w.KeyValue("sister", "ghanima")
			w.Object("enemies", func() {
				w.Array("sorted", func() {
					w.Value("harkonnen")
					w.Value("corrino")
				})
			})
		})
	})

	Expect(b.String()).To.Equal(JSON(`{
		"power": 9000,
		"atreides": {
			"name": "leto",
			"sister": "ghanima",
			"enemies": {
				"sorted": ["harkonnen", "corrino"]
			}
		}
	}`))
}

func (_ WriterTests) ArrayObject() {
	w, b := n()
	w.RootObject(func() {
		w.Array("scores", func() {
			w.ArrayObject(func() {
				w.KeyValue("points", 32)
			})
		})
	})

	Expect(b.String()).To.Equal(JSON(`{
		"scores":[{"points":32}]
	}`))
}

func (_ WriterTests) ArrayObject2() {
	w, b := n()
	w.RootObject(func() {
		w.Array("scores", func() {
			w.ArrayObject(func() {
				w.KeyValue("points", 32)
				w.KeyValue("enabled", true)
			})
			w.ArrayObject(func() {
				w.KeyValue("points", 9002)
				w.KeyValue("enabled", false)
			})
			w.Value(nil)
		})
	})

	Expect(b.String()).To.Equal(JSON(`{
		"scores":[{"points":32, "enabled":true}, {"points": 9002, "enabled": false}, null]
	}`))
}

func (_ WriterTests) RootArray() {
	w, b := n()
	w.RootArray(func() {
		w.Value(1.2)
		w.Value(false)
		w.Value("\n")
	})
	Expect(b.String()).To.Equal(JSON(`[1.2, false, "\n"]`))
}

func (_ WriterTests) MarshalJSON() {
	w, b := n()
	w.RootObject(func() {
		w.KeyValue("c", new(Marshalable))
	})
	Expect(b.String()).To.Equal(JSON(`{"c":{"ok":true}}`))
}

func (_ WriterTests) Raw() {
	w, b := n()
	w.RootObject(func() {
		w.Raw([]byte(`"test":[true]`))
	})
	Expect(b.String()).To.Equal(JSON(`{"test":[true]}`))
}

func (_ WriterTests) ArrayValuesStrings() {
	w, b := n()
	w.RootObject(func() {
		w.ArrayValues("names", []string{"leto", "jessica", "paul"})
	})
	Expect(b.String()).To.Equal(JSON(`{"names":["leto", "jessica", "paul"]}`))
}

func (_ WriterTests) ArrayValuesInts() {
	w, b := n()
	w.RootObject(func() {
		w.ArrayValues("scores", []int{2, 49299, 9001})
	})
	Expect(b.String()).To.Equal(JSON(`{"scores":[2, 49299, 9001]}`))
}

func (_ WriterTests) BoolAfterArray() {
	w, b := n()
	w.RootObject(func() {
		w.Array("scores", func() {
			w.Raw([]byte("123"))
		})
		w.KeyBool("more", false)
	})
	Expect(b.String()).To.Equal(JSON(`{"scores":[123],"more":false}`))
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

type Marshalable struct {
}

func (*Marshalable) MarshalJSON() ([]byte, error) {
	return []byte(`{"ok":true}`), nil
}
