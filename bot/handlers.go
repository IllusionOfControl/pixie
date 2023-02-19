package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleStartHelpCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Hi, i can turn your photo to pixel art. Try i sending me some pic!"+
			"Commands:\n"+
			"/start - Starts the bot and displays a welcome message.\n"+
			"/help - Displays a help message with information about how to use the bot.\n"+
			"/pixilizer - Starts the pixilize operation, requests parameters, and then returns the processed image\n"+
			"/palettizer - Starts the palettize operation, requests parameters, and then returns the processed image\n"+
			"/cancel - Reset current operation",
	)

	return msg
}

func handleDefaultCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Sorry, but command not found",
	)

	return msg
}
