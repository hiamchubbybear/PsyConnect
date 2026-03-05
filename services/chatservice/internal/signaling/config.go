package signaling


type WebRTCConfig struct {
	ICEServers []ICEServer `json:"iceServers"`
}

type ICEServer struct {
	URLs       []string `json:"urls"`
	Username   string   `json:"username,omitempty"`
	Credential string   `json:"credential,omitempty"`
}


func GetDefaultConfig() WebRTCConfig {
	return WebRTCConfig{
		ICEServers: []ICEServer{
			{
				
				URLs: []string{
					"stun:stun.l.google.com:19302",
					"stun:stun1.l.google.com:19302",
				},
			},
			
			
			
			
			
			
		},
	}
}


func (s *SignalingService) GetConfig() WebRTCConfig {
	return GetDefaultConfig()
}
