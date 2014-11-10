package jsonwriter

import (
	"bytes"
	. "github.com/karlseguin/expect"
	"testing"
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

func assertValue(value interface{}, expected string) {
	w, b := n()
	w.Value(value)
	Expect(b.String()).To.Equal(expected)
}

func n() (*Writer, *bytes.Buffer) {
	b := new(bytes.Buffer)
	return New(b), b
}
