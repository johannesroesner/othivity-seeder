package main

import (
	"log"
	"net/http"
	"os"

	"github.com/johannesroesner/othivity-seeder/address"
	"github.com/johannesroesner/othivity-seeder/club"
	"github.com/johannesroesner/othivity-seeder/profile"
)

func main() {

	if len(os.Args) < 3 {
		log.Fatal("please provide a valid jwt token and a target url")
	}

	var jwtToken = os.Args[1]
	var targetUrl = os.Args[2]

	client := http.Client{}

	address.Init()

	profile.Seed(client, jwtToken, targetUrl)
	club.Seed(client, jwtToken, targetUrl)
}
