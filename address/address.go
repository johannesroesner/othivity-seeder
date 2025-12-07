package address

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
)

type regensburgAddresses struct {
	Street  string `json:"street"`
	Min     int    `json:"min"`
	Max     int    `json:"max"`
	City    string `json:"city"`
	ZipCode string `json:"zipCode"`
}

type Address struct {
	Street      string `json:"street"`
	HouseNumber string `json:"houseNumber"`
	City        string `json:"city"`
	PostalCode  string `json:"postalCode"`
}

var regensburgAddressData []regensburgAddresses

func Init() {
	file, err := os.Open("./data/address.json")
	if err != nil {
		log.Fatal("error opening address.json: ", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("error closing address.json: ", err)
		}
	}(file)

	err = json.NewDecoder(file).Decode(&regensburgAddressData)
	if err != nil {
		log.Fatal("error parsing address.json: ", err)
	}
}

func Generate() Address {
	regensburgAddress := regensburgAddressData[randomInt(0, len(regensburgAddressData)-1)]
	houseNumber := randomInt(regensburgAddress.Min, regensburgAddress.Max)
	return Address{
		Street:      regensburgAddress.Street,
		HouseNumber: string(houseNumber),
		City:        regensburgAddress.City,
		PostalCode:  regensburgAddress.ZipCode,
	}
}

func randomInt(min int, max int) int {
	return rand.Intn(max-min+1) + min
}
