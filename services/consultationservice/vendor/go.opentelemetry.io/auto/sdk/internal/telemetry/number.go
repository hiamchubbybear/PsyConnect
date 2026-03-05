


package telemetry

import (
	"encoding/json"
	"strconv"
)



type protoInt64 int64


func (i *protoInt64) Int64() int64 { return int64(*i) }


func (i *protoInt64) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		parsedInt, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		*i = protoInt64(parsedInt)
	} else {
		var parsedInt int64
		if err := json.Unmarshal(data, &parsedInt); err != nil {
			return err
		}
		*i = protoInt64(parsedInt)
	}
	return nil
}



type protoUint64 uint64


func (i *protoUint64) Uint64() uint64 { return uint64(*i) }


func (i *protoUint64) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		parsedUint, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		*i = protoUint64(parsedUint)
	} else {
		var parsedUint uint64
		if err := json.Unmarshal(data, &parsedUint); err != nil {
			return err
		}
		*i = protoUint64(parsedUint)
	}
	return nil
}
