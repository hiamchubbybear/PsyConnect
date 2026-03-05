package toml

import (
	"fmt"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2/unstable"
)


type LocalDate struct {
	Year  int
	Month int
	Day   int
}


func (d LocalDate) AsTime(zone *time.Location) time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, zone)
}


func (d LocalDate) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}


func (d LocalDate) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}


func (d *LocalDate) UnmarshalText(b []byte) error {
	res, err := parseLocalDate(b)
	if err != nil {
		return err
	}
	*d = res
	return nil
}



type LocalTime struct {
	Hour       int 
	Minute     int 
	Second     int 
	Nanosecond int 
	Precision  int 
}





func (d LocalTime) String() string {
	s := fmt.Sprintf("%02d:%02d:%02d", d.Hour, d.Minute, d.Second)

	if d.Precision > 0 {
		s += fmt.Sprintf(".%09d", d.Nanosecond)[:d.Precision+1]
	} else if d.Nanosecond > 0 {
		
		
		s += strings.Trim(fmt.Sprintf(".%09d", d.Nanosecond), "0")
	}

	return s
}


func (d LocalTime) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}


func (d *LocalTime) UnmarshalText(b []byte) error {
	res, left, err := parseLocalTime(b)
	if err == nil && len(left) != 0 {
		err = unstable.NewParserError(left, "extra characters")
	}
	if err != nil {
		return err
	}
	*d = res
	return nil
}


type LocalDateTime struct {
	LocalDate
	LocalTime
}


func (d LocalDateTime) AsTime(zone *time.Location) time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, d.Hour, d.Minute, d.Second, d.Nanosecond, zone)
}


func (d LocalDateTime) String() string {
	return d.LocalDate.String() + "T" + d.LocalTime.String()
}


func (d LocalDateTime) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}


func (d *LocalDateTime) UnmarshalText(data []byte) error {
	res, left, err := parseLocalDateTime(data)
	if err == nil && len(left) != 0 {
		err = unstable.NewParserError(left, "extra characters")
	}
	if err != nil {
		return err
	}

	*d = res
	return nil
}
