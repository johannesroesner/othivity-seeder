package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/johannesroesner/othivity-seeder/activity"
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

	fmt.Println("seeding profiles...")
	profile.Seed(client, jwtToken, targetUrl)
	fmt.Println("seeding profiles done")

	fmt.Println("seeding clubs...")
	club.Seed(client, jwtToken, targetUrl)
	fmt.Println("seeding clubs done")

	fmt.Println("seeding activities...")
	activity.Seed(client, jwtToken, targetUrl)
	fmt.Println("seeding activities done")
}
