package bot

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"image"
	"log"
	"pixie/pixilizer"
)

type PixieBot struct {
	Bot *tgbotapi.BotAPI
}

type PixieBotConfig struct {
	Token string
	Debug bool
}

func NewBot(config *PixieBotConfig) (*PixieBot, error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, err
	}
	bot.Debug = config.Debug

	return &PixieBot{Bot: bot}, err
}

func (pixie PixieBot) Polling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := pixie.Bot.GetUpdatesChan(u)

	for update := range updates {
		var response tgbotapi.Chattable

		if update.CallbackQuery != nil {
			response, _ = handleCallbackQuery(update, pixie.Bot)
		}
		if update.Message != nil {
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "help":
					response, _ = handleStartHelpCommand(update, pixie.Bot)
				case "start":
					response, _ = handleStartHelpCommand(update, pixie.Bot)
				default:
					response, _ = handleDefaultCommand(update, pixie.Bot)
				}
			}

			if update.Message.Photo != nil {
				response, _ = handleMessagePhoto(update, pixie.Bot)
			}
		}

		if _, err := pixie.Bot.Send(response); err != nil {
			panic(err)
		}

	}
}

func handleStartHelpCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) (tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Hi, i can turn your photo to pixel art. Try i sending me some pic!",
	)

	return msg, nil
}

func handleDefaultCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) (tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Sorry, but command not found",
	)

	return msg, nil
}

func handleMessagePhoto(update tgbotapi.Update, bot *tgbotapi.BotAPI) (tgbotapi.Chattable, error) {
	log.Printf("Got photo from %d", update.Message.Chat.ID)

	if len(update.Message.Photo) == 0 {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Please, load a photo!",
		)
		return msg, nil
	}

	photo := (update.Message.Photo)[len(update.Message.Photo)-1]
	fileConfig := tgbotapi.FileConfig{FileID: photo.FileID}
	file, err := bot.GetFile(fileConfig)
	if err != nil {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Error while loading picture :(!",
		)
		return msg, nil
	}

	fileUrl, _ := bot.GetFileDirectURL(file.FileID)
	photoData, err := DownloadFile(fileUrl)
	if err != nil {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Error while downloading picture :(!",
		)
		return msg, nil
	}

	img, _, _ := image.Decode(bytes.NewReader(photoData))
	newImage := pixilizer.PixilizeImage(img, 2)
	buf, _ := ConvertImageToBytes(newImage)

	msg := tgbotapi.NewPhoto(
		update.Message.Chat.ID,
		&tgbotapi.FileBytes{Name: "photo", Bytes: buf.Bytes()},
	)

	msg.ReplyMarkup = likeDislikeKeyboard

	return msg, nil
}

func handleCallbackQuery(update tgbotapi.Update, bot *tgbotapi.BotAPI) (tgbotapi.Chattable, error) {
	switch update.CallbackQuery.Data {
	case "like":
		likeKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Like", "pass"),
			),
		)
		msg := tgbotapi.NewEditMessageReplyMarkup(
			update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			likeKeyboard,
		)
		return msg, nil
	case "dislike":
		msg := tgbotapi.NewEditMessageReplyMarkup(
			update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			dislikeKeyboard,
		)
		return msg, nil
	case "pass":

		msg := tgbotapi.NewEditMessageReplyMarkup(
			update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			likeDislikeKeyboard,
		)
		return msg, nil
	default:
		return nil, fmt.Errorf("callback handler not found")
	}
}
