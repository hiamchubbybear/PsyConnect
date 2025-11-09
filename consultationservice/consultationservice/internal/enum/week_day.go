package enum

import (
	"encoding/json"
	"errors"
	"strings"
)

type DayOfWeek int

const (
	Monday DayOfWeek = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func (d DayOfWeek) String() string {
	switch d {
	case Monday:
		return "Monday"
	case Tuesday:
		return "Tuesday"
	case Wednesday:
		return "Wednesday"
	case Thursday:
		return "Thursday"
	case Friday:
		return "Friday"
	case Saturday:
		return "Saturday"
	case Sunday:
		return "Sunday"
	default:
		return "Unknown"
	}
}

func (d *DayOfWeek) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	case "monday":
		*d = Monday
	case "tuesday":
		*d = Tuesday
	case "wednesday":
		*d = Wednesday
	case "thursday":
		*d = Thursday
	case "friday":
		*d = Friday
	case "saturday":
		*d = Saturday
	case "sunday":
		*d = Sunday
	default:
		return errors.New("invalid day of week: " + s)
	}
	return nil
}

func (d DayOfWeek) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}
