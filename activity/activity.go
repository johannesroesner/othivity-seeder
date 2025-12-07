package activity

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/johannesroesner/othivity-seeder/address"
	"github.com/johannesroesner/othivity-seeder/club"
	"github.com/johannesroesner/othivity-seeder/profile"
)

type activity struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Date        string   `json:"date"`
	Language    string   `json:"language"`
	GroupSize   int      `json:"groupSize"`
	OrganizerId string   `json:"organizerId"`
	ImageUrl    string   `json:"imageUrl"`
	Tags        []string `json:"tags"`
	StartedBy   string   `json:"startedBy"`
	TakePart    []string `json:"takePart"`
	Street      string   `json:"street"`
	HouseNumber string   `json:"houseNumber"`
	City        string   `json:"city"`
	PostalCode  string   `json:"postalCode"`
}

type responseBody struct {
	Id string `json:"id"`
}

var Ids []string

func Seed(client http.Client, jwtToken string, targetUrl string) {
	file, err := os.Open("./data/activity.json")
	if err != nil {
		log.Fatal("error opening activity.json: ", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("error closing activity.json: ", err)
		}
	}(file)

	var activities []activity
	err = json.NewDecoder(file).Decode(&activities)
	if err != nil {
		log.Fatal("error parsing activity.json: ", err)
	}

	for _, activity := range activities {
		generatedAddress := address.Generate()
		activity.Street = generatedAddress.Street
		activity.HouseNumber = generatedAddress.HouseNumber
		activity.City = generatedAddress.City
		activity.PostalCode = generatedAddress.PostalCode

		activity.Date = time.Now().AddDate(0, 0, randomInt(2, 60)).Format(time.RFC3339)

		if randomInt(0, 10) > 8 {
			activity.OrganizerId = club.Ids[rand.Intn(len(club.Ids))]
		}

		groupSize := randomInt(5, 12)
		activity.GroupSize = groupSize

		startedBy, takePart := getStarterAndParticipants(groupSize)
		activity.StartedBy = startedBy
		activity.TakePart = takePart

		body, err := json.Marshal(activity)
		if err != nil {
			log.Fatal("error marshalling activity: ", err)
		}

		request, err := http.NewRequest(http.MethodPost, targetUrl+"/api/activities", bytes.NewBuffer(body))
		if err != nil {
			log.Fatal("error creating request: ", err)
		}

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+jwtToken)

		response, err := client.Do(request)
		if err != nil {
			log.Fatal("error executing request: ", err)
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatal("error closing response body: ", err)
			}
		}(response.Body)

		if response.StatusCode != http.StatusCreated {
			log.Fatalf("error creating activity, status code: %d and title = %s", response.StatusCode, activity.Title)
		}

		var responseData responseBody
		err = json.NewDecoder(response.Body).Decode(&responseData)
		if err != nil {
			log.Fatal("error parsing response: ", err)
		}

		Ids = append(Ids, responseData.Id)

		time.Sleep(100 * time.Millisecond)
	}
}

func getStarterAndParticipants(groupSize int) (string, []string) {
	var startedBy string
	var takePart []string

	picked := profile.Ids[rand.Intn(len(profile.Ids))]
	startedBy = picked
	takePart = append(takePart, picked)

	takePartCount := randomInt(2, groupSize-1)
	for i := 0; i < takePartCount; i++ {
		picked := profile.Ids[rand.Intn(len(profile.Ids))]
		if startedBy != picked && !contains(takePart, picked) {
			takePart = append(takePart, picked)
		}
	}

	return startedBy, takePart
}

func contains(members []string, picked string) bool {
	for _, member := range members {
		if member == picked {
			return true
		}
	}
	return false
}

func randomInt(min int, max int) int {
	return rand.Intn(max-min+1) + min
}
