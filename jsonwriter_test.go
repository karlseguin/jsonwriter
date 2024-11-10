package jsonwriter

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func Test_WritesAnInt(t *testing.T) {
	assertValue(t, uint8(1), "1")
	assertValue(t, uint16(2), "2")
	assertValue(t, uint32(232134), "232134")
	assertValue(t, uint64(432434), "432434")
	assertValue(t, uint(5), "5")
	assertValue(t, int8(-3), "-3")
	assertValue(t, int16(-16), "-16")
	assertValue(t, int32(-31), "-31")
	assertValue(t, int64(-4343), "-4343")
	assertValue(t, int(-59922), "-59922")
}

func Test_WritesAFloat(t *testing.T) {
	assertValue(t, float32(1.2393), "1.2393")
	assertValue(t, float32(-49493.443), "-49493.44")
	assertValue(t, float64(99499449.23949), "9.949944923949e+07")
	assertValue(t, float64(-3290123.94994), "-3.29012394994e+06")
}

func Test_WritesAString(t *testing.T) {
	assertValue(t, `abc`, `"abc"`)
	assertValue(t, `ab"cd`, `"ab\"cd"`)
	assertValue(t, `💣`, `"💣"`)
	assertValue(t, "\\it's\n\tOver\r9000!\\ 💣 💣 💣", `"\\it's\n\tOver\r9000!\\ 💣 💣 💣"`)

	for c := 0x00; c < 0x20; c++ {
		result := fmt.Sprintf("%c", c)

		switch c {
		case '\b':
			assertValue(t, result, `"\b"`)
		case '\t':
			assertValue(t, result, `"\t"`)
		case '\n':
			assertValue(t, result, `"\n"`)
		case '\f':
			assertValue(t, result, `"\f"`)
		case '\r':
			assertValue(t, result, `"\r"`)
		default:
			assertValue(t, result, fmt.Sprintf(`"\u%04x"`, c))
		}
	}
}

func Test_WritesABool(t *testing.T) {
	assertValue(t, true, "true")
	assertValue(t, false, "false")
}

func Test_WritesANull(t *testing.T) {
	assertValue(t, nil, "null")
}

func Test_WritesATime(t *testing.T) {
	assertValue(t, time.Unix(1415677601, 9).UTC(), `"2014-11-11T03:46:41.000000009Z"`)
}

func Test_WritesAReader(t *testing.T) {
	b := bytes.NewBuffer([]byte("1234"))
	assertValue(t, b, `"MTIzNA=="`)
}

func Test_SimpleObject(t *testing.T) {
	w, b := n()
	w.RootObject(func() {
		w.KeyValue("spice", "flow")
	})
	assertString(t, b.String(), `{"spice":"flow"}`)
}

func Test_MultiValueObject(t *testing.T) {
	w, b := n()
	w.RootObject(func() {
		w.KeyValue("spice", "flow")
		w.KeyValue("over", 9000)
	})
	assertString(t, b.String(), `{"spice":"flow","over":9000}`)
}

func Test_NestedObject1(t *testing.T) {
	w, b := n()
	w.RootObject(func() {
		w.KeyValue("power", 9000)
		w.Object("atreides", func() {
			w.KeyValue("name", "leto")
			w.KeyValue("sister", "ghanima")
		})
	})

	assertString(t, b.String(), `{"power":9000,"atreides":{"name":"leto","sister":"ghanima"}}`)
}

func Test_NestedObject2(t *testing.T) {
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

	assertString(t, b.String(), `{"power":9000,"atreides":{"name":"leto","sister":"ghanima","enemies":{"sorted":["harkonnen","corrino"]}}}`)
}

func Test_ArrayObject(t *testing.T) {
	w, b := n()
	w.RootObject(func() {
		w.Array("scores", func() {
			w.ArrayObject(func() {
				w.KeyValue("points", 32)
			})
		})
	})

	assertString(t, b.String(), `{"scores":[{"points":32}]}`)
}

func Test_ArrayObject2(t *testing.T) {
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

	assertString(t, b.String(), `{"scores":[{"points":32,"enabled":true},{"points":9002,"enabled":false},null]}`)
}

func Test_RootArray(t *testing.T) {
	w, b := n()
	w.RootArray(func() {
		w.Value(1.2)
		w.Value(false)
		w.Value("\n")
	})
	assertString(t, b.String(), `[1.2,false,"\n"]`)
}

func Test_MarshalJSON(t *testing.T) {
	w, b := n()
	w.RootObject(func() {
		w.KeyValue("c", new(Marshalable))
	})
	assertString(t, b.String(), `{"c":{"ok":true}}`)
}

func Test_RawValue1(t *testing.T) {
	w, b := n()
	w.RootArray(func() {
		w.RawValue([]byte(`"abc"`))
	})
	assertString(t, b.String(), `["abc"]`)
}

func Test_RawValue2(t *testing.T) {
	w, b := n()
	w.RootArray(func() {
		w.RawValue([]byte(`"abc"`))
		w.RawValue([]byte(`"def"`))
	})
	assertString(t, b.String(), `["abc","def"]`)
}

func Test_Raw(t *testing.T) {
	w, b := n()
	w.RootObject(func() {
		w.Raw([]byte(`"test":[true]`))
	})
	assertString(t, b.String(), `{"test":[true]}`)
}

func Test_ArrayValuesStrings(t *testing.T) {
	w, b := n()
	w.RootObject(func() {
		w.ArrayValues("names", []string{"leto", "jessica", "paul"})
	})
	assertString(t, b.String(), `{"names":["leto","jessica","paul"]}`)
}

func Test_ArrayValuesInts(t *testing.T) {
	w, b := n()
	w.RootObject(func() {
		w.ArrayValues("scores", []int{2, 49299, 9001})
	})
	assertString(t, b.String(), `{"scores":[2,49299,9001]}`)
}

func Test_BoolAfterArray(t *testing.T) {
	w, b := n()
	w.RootObject(func() {
		w.Array("scores", func() {
			w.Raw([]byte("123"))
		})
		w.KeyBool("more", false)
	})
	assertString(t, b.String(), `{"scores":[123],"more":false}`)
}

func Test_Subarray(t *testing.T) {
	w, b := n()
	w.RootObject(func() {
		w.Array("scores", func() {
			w.SubArray(func() {
				w.Value(1)
				w.Value(2)
			})
			w.SubArray(func() {
				w.Value(3)
			})
		})
	})
	assertString(t, b.String(), `{"scores":[[1,2],[3]]}`)
}

func assertValue(t *testing.T, value interface{}, expected string) {
	t.Helper()
	w, b := n()
	w.Value(value)
	assertString(t, b.String(), expected)
}

func assertString(t *testing.T, value string, expected string) {
	t.Helper()
	if value != expected {
		t.Errorf("Expected '%s', got '%s'", expected, value)
	}
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
