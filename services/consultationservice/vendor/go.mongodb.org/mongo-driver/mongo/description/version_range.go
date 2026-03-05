





package description

import "fmt"


type VersionRange struct {
	Min int32
	Max int32
}


func NewVersionRange(min, max int32) VersionRange {
	return VersionRange{Min: min, Max: max}
}



func (vr VersionRange) Includes(v int32) bool {
	return v >= vr.Min && v <= vr.Max
}


func (vr *VersionRange) Equals(other *VersionRange) bool {
	if vr == nil && other == nil {
		return true
	}
	if vr == nil || other == nil {
		return false
	}
	return vr.Min == other.Min && vr.Max == other.Max
}


func (vr VersionRange) String() string {
	return fmt.Sprintf("[%d, %d]", vr.Min, vr.Max)
}
