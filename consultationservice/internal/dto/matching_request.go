package dto

type Match struct {
	Gender      string   `json:"gender" bson:"gender"`
	Address     string   `json:"address" bson:"address"`
	Languages   []string `json:"languages" bson:"languages"`
	Issues      []string `json:"issues" bson:"issues"`
	Preferences struct {
		ConsultantGender string `json:"consultant_gender" bson:"consultant_gender"`
		Mode             string `json:"mode" bson:"mode"`
		Budget           int    `json:"budget" bson:"budget"`
	} `json:"preferences" bson:"preferences"`
	Availability struct {
		Days      []string `json:"days" bson:"days"`
		TimeSlots []string `json:"time_slots" bson:"time_slots"`
	} `json:"availability" bson:"availability"`
}
