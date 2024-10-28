# go-telegram-bot

## What is this?

This is a simple Telegram bot that can be used to send messages to a Telegram chat. It is written in Go and uses the telegram bot api.

Currently its not _that_ useful, but I plan to add more features to it in the future, if I have some ideas.

Currently available commands:

- `/start` - Start conversation with the bot (Telegram standard command)
- `/help` - Get help message
- `/echo` - Echo the message sent to the bot
- `/ping` - Get the current latency of the bot
- `/uptime` - Get the uptime of the bot and server
- `/version` - Get the version of the bot
- `/ip-address` - Get information about an IP address
- `/mensa` - Get the meals of the canteen. (Of my university)

## How to use it?

1. Just text with my bot and hope its online (The username is @melo_the_bot)
2. Clone the repository and run the bot yourself

## How to run the bot?

1. Clone the repository
2. Create a new bot with the BotFather on Telegram
3. Copy the token of the bot
4. Create a new file called `.env` in the root directory of the project
5. Add the following line to the `.env` file: `TELEGRAM_BOT_TOKEN=YOUR_BOT_TOKEN`
6. Add a line to the `.env` file: `STAGE=prod`
7. Run the bot with `go run main.go` or use `docker-compose up --build` to run the bot in a docker container
