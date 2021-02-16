package test

import (
	"github.com/joho/godotenv"
	"log"
)

func TestSetup() {
	// Do something here.
	err := godotenv.Load("../.env.test")
	if err != nil {
		log.Fatalf("Error loading .env.test file; msg %s", err)
	}
}
