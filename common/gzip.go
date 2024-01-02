package common

import (
	"bytes"
	"compress/gzip"
)

func GzipDistance(a, b string) float64 {
	lenA := GetEncodeLen(a)
	lenB := GetEncodeLen(b)
	cx := GetEncodeLen(a + b)
	dist := (cx - min(lenA, lenB)) / max(lenA, lenB)
	return dist
}

func GetEncodeLen(s string) float64 {
	buf := new(bytes.Buffer)
	writer := gzip.NewWriter(buf)
	writer.Write([]byte(s))
	writer.Close()
	return float64(buf.Len())
}

type Entry struct {
	Dist float64
	Q    string
}
