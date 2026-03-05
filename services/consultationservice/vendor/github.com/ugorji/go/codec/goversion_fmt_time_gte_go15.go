


//go:build go1.5
// +build go1.5

package codec

import "time"

func fmtTime(t time.Time, fmt string, b []byte) []byte {
	return t.AppendFormat(b, fmt)
}
