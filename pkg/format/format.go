package format

import (
	"strconv"
	"time"
)

func ThousandSepartor(n int64, separator byte) string {
	in := strconv.FormatInt(n, 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits--
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = separator
		}
	}

}

func GetCurrentTimeInFullFormat() string {
	return time.Now().Format("Monday, 2 January 2006 - 03:04 PM")
}

func GetCurrentTimeInCompactFormat() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
