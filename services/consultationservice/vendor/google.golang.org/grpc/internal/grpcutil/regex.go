

package grpcutil

import "regexp"


func FullMatchWithRegex(re *regexp.Regexp, text string) bool {
	if len(text) == 0 {
		return re.MatchString(text)
	}
	re.Longest()
	rem := re.FindString(text)
	return len(rem) == len(text)
}
