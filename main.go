package main

import (
	"log"
	"os"
)

var jwtToken string

func main() {

	if len(os.Args) < 2 {
		log.Fatal("please provide a valid jwt token")
	}
	jwtToken = os.Args[1]

}
