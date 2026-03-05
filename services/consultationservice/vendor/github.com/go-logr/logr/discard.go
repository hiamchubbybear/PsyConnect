

package logr




func Discard() Logger {
	return New(nil)
}
