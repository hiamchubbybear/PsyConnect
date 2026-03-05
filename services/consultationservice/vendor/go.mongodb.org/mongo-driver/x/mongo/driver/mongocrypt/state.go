





package mongocrypt


type State int



const (
	StateError         State = 0
	NeedMongoCollInfo  State = 1
	NeedMongoMarkings  State = 2
	NeedMongoKeys      State = 3
	NeedKms            State = 4
	Ready              State = 5
	Done               State = 6
	NeedKmsCredentials State = 7
)


func (s State) String() string {
	switch s {
	case StateError:
		return "Error"
	case NeedMongoCollInfo:
		return "NeedMongoCollInfo"
	case NeedMongoMarkings:
		return "NeedMongoMarkings"
	case NeedMongoKeys:
		return "NeedMongoKeys"
	case NeedKms:
		return "NeedKms"
	case Ready:
		return "Ready"
	case Done:
		return "Done"
	case NeedKmsCredentials:
		return "NeedKmsCredentials"
	default:
		return "Unknown State"
	}
}
