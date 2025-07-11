package enum
type ConsultationMode string

const (
	ModeOnline  ConsultationMode = "online"
	ModeOffline ConsultationMode = "offline"
	ModeChat    ConsultationMode = "chat"
	ModeCall    ConsultationMode = "call"
)
