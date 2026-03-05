





//go:build !cse
// +build !cse

package mongocrypt


type Error struct {
	Code    int32
	Message string
}


func (Error) Error() string {
	panic(cseNotSupportedMsg)
}
