





package readpref

import (
	"fmt"
	"strings"
)


type Mode uint8


const (
	_ Mode = iota
	
	
	
	PrimaryMode
	
	
	
	PrimaryPreferredMode
	
	
	SecondaryMode
	
	
	
	SecondaryPreferredMode
	
	
	NearestMode
)



func ModeFromString(mode string) (Mode, error) {
	switch strings.ToLower(mode) {
	case "primary":
		return PrimaryMode, nil
	case "primarypreferred":
		return PrimaryPreferredMode, nil
	case "secondary":
		return SecondaryMode, nil
	case "secondarypreferred":
		return SecondaryPreferredMode, nil
	case "nearest":
		return NearestMode, nil
	}
	return Mode(0), fmt.Errorf("unknown read preference %v", mode)
}


func (mode Mode) String() string {
	switch mode {
	case PrimaryMode:
		return "primary"
	case PrimaryPreferredMode:
		return "primaryPreferred"
	case SecondaryMode:
		return "secondary"
	case SecondaryPreferredMode:
		return "secondaryPreferred"
	case NearestMode:
		return "nearest"
	default:
		return "unknown"
	}
}


func (mode Mode) IsValid() bool {
	switch mode {
	case PrimaryMode,
		PrimaryPreferredMode,
		SecondaryMode,
		SecondaryPreferredMode,
		NearestMode:
		return true
	default:
		return false
	}
}
