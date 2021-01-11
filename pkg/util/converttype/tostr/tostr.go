package tostr

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

func FromInterface(v interface{}) string {
	res := fmt.Sprintf("%v", v)
	return res
}

func FromInt(v int) string {
	res := strconv.FormatInt(int64(v), 10)
	return res
}

func FromUint(v uint) string {
	res := strconv.FormatUint(uint64(v), 10)
	return res
}

func FromUint64(v uint64) string {
	res := strconv.FormatUint(v, 10)
	return res
}

func FromIOReadCloser(v io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(v)
	res := buf.String()
	return res
}
