





//go:build !cse
// +build !cse

package mongocrypt


type KmsContext struct{}


func (kc *KmsContext) HostName() (string, error) {
	panic(cseNotSupportedMsg)
}


func (kc *KmsContext) Message() ([]byte, error) {
	panic(cseNotSupportedMsg)
}


func (kc *KmsContext) KMSProvider() string {
	panic(cseNotSupportedMsg)
}



func (kc *KmsContext) BytesNeeded() int32 {
	panic(cseNotSupportedMsg)
}


func (kc *KmsContext) FeedResponse([]byte) error {
	panic(cseNotSupportedMsg)
}
