package club

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/johannesroesner/othivity-seeder/address"
	"github.com/johannesroesner/othivity-seeder/profile"
)

type club struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	AccessLevel string   `json:"accessLevel"`
	ImageUrl    string   `json:"imageUrl"`
	Admins      []string `json:"admins"`
	Members     []string `json:"members"`
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
	file, err := os.Open("./data/club.json")
	if err != nil {
		log.Fatal("error opening club.json: ", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("error closing club.json: ", err)
		}
	}(file)

	var clubs []club
	err = json.NewDecoder(file).Decode(&clubs)
	if err != nil {
		log.Fatal("error parsing club.json: ", err)
	}

	for _, club := range clubs {
		generatedAddress := address.Generate()
		club.Street = generatedAddress.Street
		club.HouseNumber = generatedAddress.HouseNumber
		club.City = generatedAddress.City
		club.PostalCode = generatedAddress.PostalCode

		admins, members := getAdminsAndMembers()

		club.Admins = admins
		club.Members = members

		body, err := json.Marshal(club)
		if err != nil {
			log.Fatal("error marshalling club: ", err)
		}

		request, err := http.NewRequest(http.MethodPost, targetUrl+"/api/clubs", bytes.NewBuffer(body))
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
			log.Fatalf("error creating club, status code: %d and clubName %s", response.StatusCode, club.Name)
		}

		var responseData responseBody
		err = json.NewDecoder(response.Body).Decode(&responseData)
		if err != nil {
			log.Fatal("error parsing response: ", err)
		}

		Ids = append(Ids, responseData.Id)
	}
}

func getAdminsAndMembers() ([]string, []string) {
	var admins []string
	var members []string

	adminCount := randomInt(1, 3)

	for i := 0; i < adminCount; i++ {
		picked := profile.Ids[rand.Intn(len(profile.Ids))]
		admins = append(admins, picked)
		members = append(members, picked)
	}

	memberCount := randomInt(3, 10)
	for i := 0; i < memberCount; i++ {
		picked := profile.Ids[rand.Intn(len(profile.Ids))]
		if !contains(members, picked) {
			members = append(members, picked)
		}
	}

	return admins, members
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
