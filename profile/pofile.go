package profile

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type profile struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Language  string `json:"language"`
	Theme     string `json:"theme"`
}

type responseBody struct {
	Id string `json:"id"`
}

var Ids []string

func Seed(client http.Client, jwtToken string, targetUrl string) {
	file, err := os.Open("./data/profile.json")
	if err != nil {
		log.Fatal("error opening profile.json: ", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("error closing profile.json: ", err)
		}
	}(file)

	var profiles []profile
	err = json.NewDecoder(file).Decode(&profiles)
	if err != nil {
		log.Fatal("error parsing profile.json: ", err)
	}

	for _, profile := range profiles {
		body, err := json.Marshal(profile)
		if err != nil {
			log.Fatal("error marshalling profile: ", err)
		}

		request, err := http.NewRequest(http.MethodPost, targetUrl+"/api/profiles", bytes.NewBuffer(body))
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
			log.Fatalf("error creating profile, status code: %d and username: %s", response.StatusCode, profile.Username)
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
