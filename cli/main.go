package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ko6bxl/cm2bot"
)

func main() {

	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal(err)
	}

	cm2bot.Run(os.Getenv("KO6_TOKEN"))
}
