package enum

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
