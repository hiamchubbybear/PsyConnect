package redis

import "time"


func NewCmdResult(val interface{}, err error) *Cmd {
	var cmd Cmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewSliceResult(val []interface{}, err error) *SliceCmd {
	var cmd SliceCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewStatusResult(val string, err error) *StatusCmd {
	var cmd StatusCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewIntResult(val int64, err error) *IntCmd {
	var cmd IntCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewDurationResult(val time.Duration, err error) *DurationCmd {
	var cmd DurationCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewBoolResult(val bool, err error) *BoolCmd {
	var cmd BoolCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewStringResult(val string, err error) *StringCmd {
	var cmd StringCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewFloatResult(val float64, err error) *FloatCmd {
	var cmd FloatCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewStringSliceResult(val []string, err error) *StringSliceCmd {
	var cmd StringSliceCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewBoolSliceResult(val []bool, err error) *BoolSliceCmd {
	var cmd BoolSliceCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewFloatSliceResult(val []float64, err error) *FloatSliceCmd {
	var cmd FloatSliceCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewMapStringStringResult(val map[string]string, err error) *MapStringStringCmd {
	var cmd MapStringStringCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewMapStringIntCmdResult(val map[string]int64, err error) *MapStringIntCmd {
	var cmd MapStringIntCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewTimeCmdResult(val time.Time, err error) *TimeCmd {
	var cmd TimeCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewZSliceCmdResult(val []Z, err error) *ZSliceCmd {
	var cmd ZSliceCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewZWithKeyCmdResult(val *ZWithKey, err error) *ZWithKeyCmd {
	var cmd ZWithKeyCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewScanCmdResult(keys []string, cursor uint64, err error) *ScanCmd {
	var cmd ScanCmd
	cmd.page = keys
	cmd.cursor = cursor
	cmd.SetErr(err)
	return &cmd
}


func NewClusterSlotsCmdResult(val []ClusterSlot, err error) *ClusterSlotsCmd {
	var cmd ClusterSlotsCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewGeoLocationCmdResult(val []GeoLocation, err error) *GeoLocationCmd {
	var cmd GeoLocationCmd
	cmd.locations = val
	cmd.SetErr(err)
	return &cmd
}


func NewGeoPosCmdResult(val []*GeoPos, err error) *GeoPosCmd {
	var cmd GeoPosCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewCommandsInfoCmdResult(val map[string]*CommandInfo, err error) *CommandsInfoCmd {
	var cmd CommandsInfoCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewXMessageSliceCmdResult(val []XMessage, err error) *XMessageSliceCmd {
	var cmd XMessageSliceCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewXStreamSliceCmdResult(val []XStream, err error) *XStreamSliceCmd {
	var cmd XStreamSliceCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}


func NewXPendingResult(val *XPending, err error) *XPendingCmd {
	var cmd XPendingCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}
