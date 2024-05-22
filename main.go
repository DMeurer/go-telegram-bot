package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	// Telegram API
	"github.com/joho/godotenv"
)

// This bot is as basic as it gets - it simply repeats everything you say.
func main() {
	// Get token from the environment variable
	err := godotenv.Load()
	if err != nil {
		panic("Error loading env file.")
	}

	// determine stage
	stage := os.Getenv("STAGE")
	token := ""
	if stage == "" {
		panic("STAGE environment variable is empty")
	} else if stage == "dev" {
		token = os.Getenv("TELEGRAM_APITOKEN_DEV")
		if token == "" {
			panic("TELEGRAM_APITOKEN_DEV environment variable is empty")
		}
	} else if stage == "prod" {
		token = os.Getenv("TELEGRAM_APITOKEN")
		if token == "" {
			panic("TELEGRAM_APITOKEN environment variable is empty")
		}
	} else {
		panic("STAGE environment variable is not set to dev or prod")
	}

	// Create bot from environment value.
	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	// Default
	dispatcher.AddHandler(handlers.NewCommand("start", start))

	// Actual commands
	dispatcher.AddHandler(handlers.NewCommand("echo", echo))
	dispatcher.AddHandler(handlers.NewCommand("ip-address", apiIpLookup))

	// Debugging and utility commands
	dispatcher.AddHandler(handlers.NewCommand("ping", ping))
	dispatcher.AddHandler(handlers.NewCommand("uptime", uptime))
	dispatcher.AddHandler(handlers.NewCommand("dev", devDebug))
	dispatcher.AddHandler(handlers.NewCommand("version", version))

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started...\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

func version(b *gotgbot.Bot, ctx *ext.Context) error {
	// TODO: Change for every commit
	_, err := ctx.EffectiveMessage.Reply(b, "Currently on Version v1.0.2", nil)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}
