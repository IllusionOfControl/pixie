package bot

import (
	"bytes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"image"
	"log"
	"pixie/pixie"
)

func handleStartHelpCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Hi, i can turn your photo to pixel art. Try i sending me some pic!",
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

func handleMessagePhoto(update tgbotapi.Update, bot *tgbotapi.BotAPI) tgbotapi.Chattable {
	log.Printf("Got photo from %d", update.Message.Chat.ID)

	if len(update.Message.Photo) == 0 {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Please, load a photo!",
		)
		return msg
	}

	photo := (update.Message.Photo)[len(update.Message.Photo)-1]
	fileConfig := tgbotapi.FileConfig{FileID: photo.FileID}
	file, err := bot.GetFile(fileConfig)
	if err != nil {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Error while loading picture :(!",
		)
		return msg
	}

	fileUrl, _ := bot.GetFileDirectURL(file.FileID)
	photoData, err := DownloadFile(fileUrl)
	if err != nil {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Error while downloading picture :(!",
		)
		return msg
	}

	img, _, _ := image.Decode(bytes.NewReader(photoData))
	//newImage := pixie.PixilizeImage(img, 2)
	palette := pixie.LoadPaletteFromImage(img, 16)
	newImage := pixie.DrawImageWithPalette(img, palette)
	buf, _ := ConvertImageToBytes(newImage)

	msg := tgbotapi.NewPhoto(
		update.Message.Chat.ID,
		&tgbotapi.FileBytes{Name: "photo", Bytes: buf.Bytes()},
	)

	msg.ReplyMarkup = likeDislikeKeyboard

	return msg
}
