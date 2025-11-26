package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	ProfileId         string   `json:"profile_id,omitempty" bson:"profile_id"`
	Address           string   `json:"address,omitempty" bson:"address" required:"true"`
	Languages         []string `json:"languages,omitempty" bson:"languages" required:"true"`
	IssueDetail       []string `json:"issue_detail,omitempty" bson:"issue_detail"`
	ConsultationModes []string `json:"consultation_modes,omitempty" bson:"consultation_modes" required:"true"`
	RangePrice        int      `json:"rage_price,omitempty" bson:"rage_price" required:"true"`
	Availability      struct {
		Days      []string `json:"days,omitempty" bson:"days" required:"true"`
		TimeSlots []string `json:"time_slots,omitempty" bson:"time_slots" required:"true"`
	} `json:"availability,omitempty" bson:"availability"`
	PreferredTherapistGender   string   `json:"preferred_therapist_gender,omitempty" bson:"preferred_therapist_gender"`
	ExperienceLevel            string   `json:"experience_level,omitempty" bson:"experience_level"`
	TherapistSpecialization    []string `json:"specialization,omitempty" bson:"therapist_specialization"`
	UrgencyLevel               string   `json:"urgency_level,omitempty" bson:"urgency_level"`
	SessionDuration            int      `json:"session_duration,omitempty" bson:"session_duration"`
	PreferredTherapistLanguage []string `json:"preferred_therapist_language,omitempty" bson:"preferred_therapist_language"`
	IsFlexibleWithSchedule     bool     `json:"is_flexible_with_schedule,omitempty" bson:"is_flexible_with_schedule"`
}
type ProfessionalTitle struct {
	Code    string `json:"code" bson:"code"`
	Display string `json:"display" bson:"display"`
}

type Degree struct {
	Type        string `json:"type" bson:"type"`
	Field       string `json:"field" bson:"field"`
	Institution string `json:"institution" bson:"institution"`
	Year        int    `json:"year" bson:"year"`
}

type Certification struct {
	Name   string `json:"name" bson:"name"`
	Issuer string `json:"issuer" bson:"issuer"`
	Year   int    `json:"year" bson:"year"`
}

type ProfessionalInfo struct {
	Title           ProfessionalTitle `json:"title" bson:"title"`
	Degrees         []Degree          `json:"degrees" bson:"degrees"`
	Certifications  []Certification   `json:"certifications" bson:"certifications"`
	ExperienceYears int               `json:"experience_years" bson:"experience_years"`
}

type Therapist struct {
	ProfileId         string   `json:"profile_id,omitempty" bson:"profile_id"`
	Address           string   `json:"address,omitempty" bson:"address"`
	Languages         []string `json:"languages,omitempty" bson:"languages"`
	Specialization    []string `json:"specialization,omitempty" bson:"specialization"`
	ConsultationModes []string `json:"consultation_modes,omitempty" bson:"consultation_modes"`
	Experience        int      `json:"experience,omitempty" bson:"experience"`
	Rating            float64  `json:"rating,omitempty" bson:"rating"`
	Currency          string   `json:"currency" bson:"currency"`
	RagePrice         int      `json:"rage_price,omitempty" bson:"rage_price"`
	IsAvailable       bool     `json:"is_available,omitempty" bson:"is_available"`
	Availability      struct {
		Days      []string `json:"days,omitempty" bson:"days"`
		TimeSlots []string `json:"time_slots,omitempty" bson:"time_slots"`
	} `json:"availability,omitempty" bson:"availability"`
	CurrentSession   []string         `json:"current_session,omitempty" bson:"current_session"`
	MatchedClients   []string         `json:"matched_clients,omitempty" bson:"matched_clients"`
	AvatarOverride   string           `json:"avatar_override" bson:"avatar_override"`
	Name             string           `json:"name" bson:"name"`
	ProfessionalInfo ProfessionalInfo `json:"professional_info,omitempty" bson:"professional_info"`
}

// --- Pools dữ liệu để random ---
var (
	names = []string{
		"Dr. Alice Nguyen", "Dr. Bob Tran", "Dr. Carol Pham", "Dr. David Le",
		"Dr. Emma Vo", "Dr. Frank Hoang", "Dr. Grace Dang", "Dr. Henry Do",
		"Dr. Ivy Truong", "Dr. Jack Bui", "Dr. Kelly Mai", "Dr. Liam Chau",
	}

	languagePool = []string{
		"English", "Vietnamese", "French", "German", "Japanese",
		"Korean", "Spanish", "Chinese", "Italian", "Portuguese",
	}

	specializationPool = []string{
		"Anxiety", "Depression", "Stress Management", "PTSD",
		"Family Therapy", "Child Psychology", "Addiction",
		"Career Counseling", "Couples Therapy", "Eating Disorders",
	}

	degreeFields = []string{
		"Psychology", "Clinical Psychology", "Counseling",
		"Neuroscience", "Social Work", "Education Psychology",
	}

	institutions = []string{
		"Harvard University", "Stanford University", "University of Toronto",
		"University of Oxford", "University of Melbourne", "National University of Singapore",
		"Hanoi Medical University", "Ho Chi Minh City University of Education",
	}

	certifications = []string{
		"Cognitive Behavioral Therapy", "Mindfulness-Based Stress Reduction",
		"Family Systems Therapy", "Trauma-Informed Care", "Dialectical Behavior Therapy",
		"Positive Psychology Coaching", "Child Development Specialist",
	}

	issuers = []string{
		"APA", "Harvard Medical School", "WHO", "Stanford Medicine",
		"Oxford Psychiatry Dept", "Vietnam National University",
	}
)

func randomSubset(pool []string, min, max int) []string {
	rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	size := min + rand.Intn(max-min+1)
	if size > len(pool) {
		size = len(pool)
	}
	return pool[:size]
}

func randomDegrees(count int) []Degree {
	degrees := []Degree{}
	for i := 0; i < count; i++ {
		degrees = append(degrees, Degree{
			Type:        []string{"Bachelor", "Master", "PhD"}[rand.Intn(3)],
			Field:       degreeFields[rand.Intn(len(degreeFields))],
			Institution: institutions[rand.Intn(len(institutions))],
			Year:        2005 + rand.Intn(15),
		})
	}
	return degrees
}

func randomCertifications(count int) []Certification {
	certs := []Certification{}
	for i := 0; i < count; i++ {
		certs = append(certs, Certification{
			Name:   certifications[rand.Intn(len(certifications))],
			Issuer: issuers[rand.Intn(len(issuers))],
			Year:   2010 + rand.Intn(13),
		})
	}
	return certs
}

func main() {
	rand.Seed(time.Now().UnixNano())
	const MONGO_URI = "mongodb://root:123456@localhost:27017/?authSource=admin"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())
	clientCollection := client.Database("consultationservice").Collection("clients")

	therapistCollection := client.Database("consultationservice").Collection("therapists")

	var therapists []interface{}
	for i := 0; i < 10; i++ {
		name := names[rand.Intn(len(names))]
		t := Therapist{
			ProfileId:         uuid.New().String(),
			Address:           fmt.Sprintf("%d Example Street, District %d, Ho Chi Minh City", rand.Intn(200), rand.Intn(12)+1),
			Languages:         randomSubset(languagePool, 2, 4),
			Specialization:    randomSubset(specializationPool, 2, 5),
			ConsultationModes: randomSubset([]string{"Online", "Offline"}, 1, 2),
			Experience:        rand.Intn(15) + 1,
			Rating:            3.5 + rand.Float64()*1.5, // 3.5 -> 5.0
			Currency:          "USD",
			RagePrice:         40 + rand.Intn(60),
			IsAvailable:       rand.Intn(2) == 0,
			AvatarOverride:    fmt.Sprintf("https://i.pinimg.com/736x/22/3f/e7/223fe7f4bb42917325a8375ecf7a7abf.jpg"),
			Name:              name,
			ProfessionalInfo: ProfessionalInfo{
				Title: ProfessionalTitle{
					Code:    "PSY",
					Display: "Psychologist",
				},
				Degrees:         randomDegrees(rand.Intn(2) + 1),
				Certifications:  randomCertifications(rand.Intn(3) + 1),
				ExperienceYears: rand.Intn(20),
			},
		}

		// Availability
		t.Availability.Days = randomSubset([]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}, 2, 4)
		t.Availability.TimeSlots = randomSubset([]string{"09:00-11:00", "13:00-15:00", "15:00-17:00", "19:00-21:00"}, 1, 3)

		t.CurrentSession = []string{}
		t.MatchedClients = []string{}

		therapists = append(therapists, t)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := therapistCollection.InsertMany(ctx, therapists)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted IDs:", result.InsertedIDs)

	c := Client{
		ProfileId:                  uuid.New().String(),
		Address:                    "456 Nguyen Trai, District 5, Ho Chi Minh City",
		Languages:                  []string{"Vietnamese", "English"},
		IssueDetail:                []string{"Stress", "Anxiety"},
		ConsultationModes:          []string{"Online"},
		RangePrice:                 100,
		PreferredTherapistGender:   "Female",
		ExperienceLevel:            "Senior",
		TherapistSpecialization:    []string{"Anxiety", "Stress Management"},
		UrgencyLevel:               "High",
		SessionDuration:            60,
		PreferredTherapistLanguage: []string{"Vietnamese"},
		IsFlexibleWithSchedule:     true,
	}
	c.Availability.Days = []string{"Tuesday", "Thursday"}
	c.Availability.TimeSlots = []string{"18:00-20:00"}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resultClient, err := clientCollection.InsertOne(ctx, c)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Client inserted :", resultClient.InsertedID)

}
