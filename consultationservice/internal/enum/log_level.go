package enum

type LogLevel string

const (
	LOG   LogLevel = "Log"
	AUDIT LogLevel = "Audit"
	WARN  LogLevel = "Warn"
	ERROR LogLevel = "Error"
)
