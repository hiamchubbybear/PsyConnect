package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Therapist struct {
	ProfileId         string   `json:"profile_id,omitempty" bson:"profile_id"`
	Address           string   `json:"address,omitempty" bson:"address"`
	Languages         []string `json:"languages,omitempty" bson:"languages"`
	Specialization    []string `json:"specialization,omitempty" bson:"specialization"`
	ConsultationModes []string `json:"consultation_modes,omitempty" bson:"consultation_modes"`
	Experience        int      `json:"experience,omitempty" bson:"experience"`
	Rating            float64  `json:"rating,omitempty" bson:"rating"`
	Currency          string   `json:"currency" bson:"currency"`
	RagePrice         int      `json:"rage_price,omitempty" bson:"rage-price"`
	IsAvailable       bool     `json:"is_available,omitempty" bson:"is_available"`
	Availability      struct {
		Days      []string `json:"days,omitempty" bson:"days"`
		TimeSlots []string `json:"time_slots,omitempty" bson:"time_slots"`
	} `json:"availability,omitempty" bson:"availability"`
	CurrentSession []string `json:"current_session,omitempty" bson:"current_session"`
	MatchedClients []string `json:"matched_clients,omitempty" bson:"matched_clients"`
}

func main() {
	const MONGO_URI = "mongodb://root:123456@localhost:27017"

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("consultationservice").Collection("therapists")

	var therapists []interface{}
	for i := 1; i <= 10; i++ {
		t := Therapist{
			ProfileId:         fmt.Sprintf("therapist_%d", i),
			Address:           fmt.Sprintf("123 Street %d, City", i),
			Languages:         []string{"English", "Vietnamese"},
			Specialization:    []string{"Anxiety", "Depression"},
			ConsultationModes: []string{"Online", "Offline"},
			Experience:        3 + i,
			Rating:            4.0 + float64(i)*0.1,
			Currency:          "USD",
			RagePrice:         50 + i*10,
			IsAvailable:       i%2 == 0,
		}
		t.Availability.Days = []string{"Monday", "Wednesday", "Friday"}
		t.Availability.TimeSlots = []string{"09:00-11:00", "14:00-16:00"}
		t.CurrentSession = []string{}
		t.MatchedClients = []string{}

		therapists = append(therapists, t)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertMany(ctx, therapists)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted IDs:", result.InsertedIDs)
}

