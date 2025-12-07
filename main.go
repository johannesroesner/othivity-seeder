package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var jwtToken string
var targetUrl string

var profileIds []string

type Response struct {
	Id string `json:"id"`
}

func main() {

	if len(os.Args) < 3 {
		log.Fatal("please provide a valid jwt token and a target url")
	}

	jwtToken = os.Args[1]
	targetUrl = os.Args[2]

	client := http.Client{}

	seedProfiles(client)
}

type Profile struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Language  string `json:"language"`
	Theme     string `json:"theme"`
}

func seedProfiles(client http.Client) {
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

	var profiles []Profile
	err = json.NewDecoder(file).Decode(&profiles)
	if err != nil {
		log.Fatal("error parsing profile.json: ", err)
	}

	for _, index := range profiles {
		body, err := json.Marshal(index)
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
			log.Fatalf("error creating profile, status code: %d", response.StatusCode)
		}

		var responseData Response
		err = json.NewDecoder(response.Body).Decode(&responseData)
		if err != nil {
			log.Fatal("error parsing response: ", err)
		}

		profileIds = append(profileIds, responseData.Id)

		time.Sleep(100 * time.Millisecond)
	}

}
