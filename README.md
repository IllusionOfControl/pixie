# Golang Telegram Bot for Pixilizing Images
This is a Golang Telegram bot that can pixilize images. It uses the go-telegram-bot-api library to interact with the Telegram Bot API and the gocv library for image processing.

## Getting Started
Before starting, you will need to create a bot on Telegram and obtain its API token. You can follow the instructions on Telegram's documentation to do this.

### Installation
To install the bot and its dependencies, run the following command:

```sh
go get -u github.com/go-telegram-bot-api/telegram-bot-api/v5
go get -u github.com/joho/godotenv
go get -u github.com/mccutchen/palettor
go get -u github.com/nfnt/resize
```

### Docker
### Usage
To run the bot, use the following command:

```
go run main.go
```
Once the bot is running, you can interact with it by sending messages to it on Telegram. The bot accepts the following commands:

`/start` - Starts the bot and displays a welcome message.

`/help` - Displays a help message with information about how to use the bot.

`/pixilizer` - starts the pixilize operation, requests parameters, and then returns the processed image

`/palettizer` - starts the palettize operation, requests parameters, and then returns the processed image

`/cancel` - Reset current operation

To pixilize an image, the user needs to send an image file to the bot by selecting the file using the paperclip icon in the chat. The bot will then ask for the size of the pixels (in percentage), the number of clusters for K-means algorithm, and the color offset. The bot will then pixilize the image and send the result back to the user.

Contributors
This project was created by Illusion of control. Contributions are welcome, and you can submit pull requests or file issues on the project's GitHub repository.