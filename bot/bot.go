package bot

import (
	"bytes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"image"
	"pixie/pixie"
	"strconv"
)

func Polling(api *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := api.GetUpdatesChan(u)

	for update := range updates {
		var response tgbotapi.Chattable

		if update.CallbackQuery != nil {
			response = processCallbackQuery(update, api)
		}
		if update.Message != nil {
			if update.Message.IsCommand() {
				response = processUserCommand(update, api)
			} else {
				response = processUserMessage(update, api)
			}
		}

		if _, err := api.Send(response); err != nil {
			panic(err)
		}
	}
}

func NewBot(token string, debug bool) (*tgbotapi.BotAPI, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	api.Debug = debug

	return api, err
}

func processUserMessage(update tgbotapi.Update, api *tgbotapi.BotAPI) tgbotapi.Chattable {
	userState := GetUserState(update.Message.Chat.ID)

	switch userState.State {
	case PixilizerAskedForCount:
		if _, err := strconv.Atoi(update.Message.Text); err != nil {
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Input value is not a number",
			)
			return msg
		} else {
			userState.Context["count"] = update.Message.Text
			userState.State = PixilizerAskedForImage

			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Send the photo for processing",
			)
			return msg
		}
	case PixilizerAskedForImage:
		if len(update.Message.Photo) == 0 {
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Please, load a photo!",
			)
			return msg
		} else {
			userState.Context["imageId"] = (update.Message.Photo)[len(update.Message.Photo)-1].FileID

			fileUrl, _ := api.GetFileDirectURL(userState.Context["imageId"])
			photoData, err := DownloadFile(fileUrl)
			if err != nil {
				msg := tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Error while downloading picture :(",
				)
				return msg
			}

			image, _, _ := image.Decode(bytes.NewReader(photoData))

			dimensions, _ := strconv.Atoi(userState.Context["count"])
			newImage := pixie.PixilizeImage(image, dimensions)

			buf, _ := ConvertImageToBytes(newImage)

			msg := tgbotapi.NewPhoto(
				update.Message.Chat.ID,
				&tgbotapi.FileBytes{Name: "photo", Bytes: buf.Bytes()},
			)

			msg.ReplyMarkup = likeDislikeKeyboard

			userState.State = StartState
			userState.ClearContext()
			return msg
		}
	case PalettizerAskedForDimension:
		if _, err := strconv.Atoi(update.Message.Text); err != nil {
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Input value is not a number",
			)
			return msg
		} else {
			userState.Context["count"] = update.Message.Text
			userState.State = PalettizerAskedForImage

			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Send the photo for processing",
			)
			return msg
		}
	case PalettizerAskedForImage:
		if len(update.Message.Photo) == 0 {
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Please, load a photo!",
			)
			return msg
		} else {
			userState.Context["imageId"] = (update.Message.Photo)[len(update.Message.Photo)-1].FileID

			fileUrl, _ := api.GetFileDirectURL(userState.Context["imageId"])
			photoData, err := DownloadFile(fileUrl)
			if err != nil {
				msg := tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Error while downloading picture :(",
				)
				return msg
			}

			image, _, _ := image.Decode(bytes.NewReader(photoData))

			colors, _ := strconv.Atoi(userState.Context["count"])
			newImage, err := pixie.Palettize(image, colors, true)
			if err != nil {
				msg := tgbotapi.NewMessage(
					update.Message.Chat.ID,
					err.Error(),
				)
				return msg
			}
			buf, _ := ConvertImageToBytes(newImage)

			msg := tgbotapi.NewPhoto(
				update.Message.Chat.ID,
				&tgbotapi.FileBytes{Name: "photo", Bytes: buf.Bytes()},
			)

			msg.ReplyMarkup = likeDislikeKeyboard

			userState.State = StartState
			userState.ClearContext()
			return msg
		}
	}
	return nil
}

func processUserCommand(update tgbotapi.Update, api *tgbotapi.BotAPI) tgbotapi.Chattable {
	userState := GetUserState(update.Message.Chat.ID)

	switch update.Message.Command() {
	case "help", "start":
		return handleStartHelpCommand(update, api)
	case "cancel":
		userState.State = StartState
		userState.ClearContext()

		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Hi, i can turn your photo to pixel art. Try i sending me some pic!",
		)

		return msg
	case "pixilizer":
		userState.State = PixilizerAskedForCount
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Input pixel dimension",
		)
		return msg
	case "palletizer":
		userState.State = PalettizerAskedForDimension
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Input colors count",
		)
		return msg
	default:
		return handleDefaultCommand(update, api)
	}
}

func processCallbackQuery(update tgbotapi.Update, bot *tgbotapi.BotAPI) tgbotapi.Chattable {
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
		return msg
	case "dislike":
		msg := tgbotapi.NewEditMessageReplyMarkup(
			update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			dislikeKeyboard,
		)
		return msg
	case "pass":
		msg := tgbotapi.NewEditMessageReplyMarkup(
			update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			likeDislikeKeyboard,
		)
		return msg
	default:
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Callback handler not found",
		)
		return msg
	}
}
