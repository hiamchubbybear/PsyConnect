





package options

import (
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


type RunCmdOptions struct {
	
	
	ReadPreference *readpref.ReadPref
}


func RunCmd() *RunCmdOptions {
	return &RunCmdOptions{}
}


func (rc *RunCmdOptions) SetReadPreference(rp *readpref.ReadPref) *RunCmdOptions {
	rc.ReadPreference = rp
	return rc
}





func MergeRunCmdOptions(opts ...*RunCmdOptions) *RunCmdOptions {
	rc := RunCmd()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.ReadPreference != nil {
			rc.ReadPreference = opt.ReadPreference
		}
	}

	return rc
}
