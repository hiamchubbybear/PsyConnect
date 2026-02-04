package domain

type KafkaNotificationEvent struct {
	Key   string      `json:"key"`
	Value EventDetail `json:"value"`
}

type EventDetail struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type IncomingCallPayload struct {
	SessionID    string `json:"session_id"`
	CallerID     string `json:"caller_id"`
	CallerName   string `json:"caller_name"`
	CallerAvatar string `json:"caller_avatar"`
	SignalData   string `json:"signal_data,omitempty"`
}
