package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	bot "pixie/bot"
	"strconv"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadEnv()
	botConfig := &bot.PixieBotConfig{}

	if value, exists := os.LookupEnv("TOKEN"); exists {
		botConfig.Token = value
	} else {
		fmt.Println("Error: TOKEN environment variable not set.")
		fmt.Println("Please set the TOKEN environment variable before running this program.")
		os.Exit(1)
	}

	if value, exists := os.LookupEnv("DEBUG"); exists {
		if debug, err := strconv.ParseBool(value); err != nil {
			botConfig.Debug = debug
		} else {
			fmt.Errorf(err.Error())
		}
	} else {
		botConfig.Debug = false
	}

	bot, err := bot.NewBot(botConfig)
	if err != nil {
		panic(err)
	}

	bot.Polling()
}
