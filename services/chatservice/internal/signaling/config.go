package signaling

// WebRTC Configuration
type WebRTCConfig struct {
	ICEServers []ICEServer `json:"iceServers"`
}

type ICEServer struct {
	URLs       []string `json:"urls"`
	Username   string   `json:"username,omitempty"`
	Credential string   `json:"credential,omitempty"`
}

// GetDefaultConfig returns default STUN/TURN server configuration
func GetDefaultConfig() WebRTCConfig {
	return WebRTCConfig{
		ICEServers: []ICEServer{
			{
				// Google's public STUN servers
				URLs: []string{
					"stun:stun.l.google.com:19302",
					"stun:stun1.l.google.com:19302",
				},
			},
			// Add TURN server if needed (requires credentials)
			// {
			// 	URLs:       []string{"turn:your-turn-server.com:3478"},
			// 	Username:   "username",
			// 	Credential: "password",
			// },
		},
	}
}

// GetConfig returns WebRTC configuration for clients
func (s *SignalingService) GetConfig() WebRTCConfig {
	return GetDefaultConfig()
}
