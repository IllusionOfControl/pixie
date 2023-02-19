package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"pixie/bot"
	"strconv"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	var (
		token string
		debug bool
	)
	loadEnv()

	if value, exists := os.LookupEnv("TOKEN"); exists {
		token = value
	} else {
		fmt.Println("Error: TOKEN environment variable not set.")
		fmt.Println("Please set the TOKEN environment variable before running this program.")
		os.Exit(1)
	}

	if value, exists := os.LookupEnv("DEBUG"); exists {
		if value, err := strconv.ParseBool(value); err != nil {
			debug = value
		} else {
			fmt.Errorf(err.Error())
		}
	} else {
		debug = false
	}

	pixieBot, err := bot.NewBot(token, debug)
	if err != nil {
		panic(err)
	}

	bot.Polling(pixieBot)
}
