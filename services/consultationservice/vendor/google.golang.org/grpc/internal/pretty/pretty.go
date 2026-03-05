


package pretty

import (
	"bytes"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/protoadapt"
)

const jsonIndent = "  "




func ToJSON(e any) string {
	if ee, ok := e.(protoadapt.MessageV1); ok {
		e = protoadapt.MessageV2Of(ee)
	}

	if ee, ok := e.(protoadapt.MessageV2); ok {
		mm := protojson.MarshalOptions{
			Indent:    jsonIndent,
			Multiline: true,
		}
		ret, err := mm.Marshal(ee)
		if err != nil {
			
			
			
			return fmt.Sprintf("%+v", ee)
		}
		return string(ret)
	}

	ret, err := json.MarshalIndent(e, "", jsonIndent)
	if err != nil {
		return fmt.Sprintf("%+v", e)
	}
	return string(ret)
}




func FormatJSON(b []byte) string {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", jsonIndent)
	if err != nil {
		return string(b)
	}
	return out.String()
}
