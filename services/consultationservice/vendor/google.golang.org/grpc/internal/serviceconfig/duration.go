

package serviceconfig

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)





type Duration time.Duration

func (d Duration) String() string {
	return fmt.Sprint(time.Duration(d))
}


func (d Duration) MarshalJSON() ([]byte, error) {
	ns := time.Duration(d).Nanoseconds()
	sec := ns / int64(time.Second)
	ns = ns % int64(time.Second)

	var sign string
	if sec < 0 || ns < 0 {
		sign, sec, ns = "-", -1*sec, -1*ns
	}

	
	
	str := fmt.Sprintf("%s%d.%09d", sign, sec, ns)
	str = strings.TrimSuffix(str, "000")
	str = strings.TrimSuffix(str, "000")
	str = strings.TrimSuffix(str, ".000")
	return []byte(fmt.Sprintf("\"%ss\"", str)), nil
}


func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if !strings.HasSuffix(s, "s") {
		return fmt.Errorf("malformed duration %q: missing seconds unit", s)
	}
	neg := false
	if s[0] == '-' {
		neg = true
		s = s[1:]
	}
	ss := strings.SplitN(s[:len(s)-1], ".", 3)
	if len(ss) > 2 {
		return fmt.Errorf("malformed duration %q: too many decimals", s)
	}
	
	
	hasDigits := false
	var sec, ns int64
	if len(ss[0]) > 0 {
		var err error
		if sec, err = strconv.ParseInt(ss[0], 10, 64); err != nil {
			return fmt.Errorf("malformed duration %q: %v", s, err)
		}
		
		const maxProtoSeconds = 315_576_000_000
		if sec > maxProtoSeconds {
			return fmt.Errorf("out of range: %q", s)
		}
		hasDigits = true
	}
	if len(ss) == 2 && len(ss[1]) > 0 {
		if len(ss[1]) > 9 {
			return fmt.Errorf("malformed duration %q: too many digits after decimal", s)
		}
		var err error
		if ns, err = strconv.ParseInt(ss[1], 10, 64); err != nil {
			return fmt.Errorf("malformed duration %q: %v", s, err)
		}
		for i := 9; i > len(ss[1]); i-- {
			ns *= 10
		}
		hasDigits = true
	}
	if !hasDigits {
		return fmt.Errorf("malformed duration %q: contains no numbers", s)
	}

	if neg {
		sec *= -1
		ns *= -1
	}

	
	const maxSeconds = math.MaxInt64 / int64(time.Second)
	const maxNanosAtMaxSeconds = math.MaxInt64 % int64(time.Second)
	const minSeconds = math.MinInt64 / int64(time.Second)
	const minNanosAtMinSeconds = math.MinInt64 % int64(time.Second)

	if sec > maxSeconds || (sec == maxSeconds && ns >= maxNanosAtMaxSeconds) {
		*d = Duration(math.MaxInt64)
	} else if sec < minSeconds || (sec == minSeconds && ns <= minNanosAtMinSeconds) {
		*d = Duration(math.MinInt64)
	} else {
		*d = Duration(sec*int64(time.Second) + ns)
	}
	return nil
}
