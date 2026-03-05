

package logr




type contextKey struct{}


type notFoundError struct{}

func (notFoundError) Error() string {
	return "no logr.Logger was present"
}

func (notFoundError) IsNotFound() bool {
	return true
}
