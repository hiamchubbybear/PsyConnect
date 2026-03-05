



package protoimpl

import (
	"google.golang.org/protobuf/internal/version"
)

const (
	
	
	MaxVersion = version.Minor

	
	
	
	GenVersion = 20

	
	
	MinVersion = 0
)



























type EnforceVersion uint




const (
	_ = EnforceVersion(GenVersion - MinVersion)
	_ = EnforceVersion(MaxVersion - GenVersion)
)
