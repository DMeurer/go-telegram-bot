package main

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"io"
	"log"
	"net/http"
	"os/exec"
)

// echo replies to a messages with its own contents.
func echo(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, ctx.EffectiveMessage.Text, nil)
	if err != nil {
		return fmt.Errorf("failed to echo message: %w", err)
	}
	return nil
}

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, "Hello, I am a bot. I am here to help you.\nIf i decided to switch the repo to public, the code can be found here:\nhttps://github.com/DMeurer/go-telegram-bot", nil)
	if err != nil {
		return fmt.Errorf("failed to echo message: %w", err)
	}
	return nil
}

func ping(b *gotgbot.Bot, ctx *ext.Context) error {
	// get ping of the bot by pinging 8.8.8.8
	out, err := exec.Command("sh", "./shell-scripts/ping.sh").Output()
	if err != nil {
		log.Fatal(err)
	}
	_, err = ctx.EffectiveMessage.Reply(b, string(out), nil)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func uptime(b *gotgbot.Bot, ctx *ext.Context) error {
	// get uptime of the bot by executing uptime command
	out, err := exec.Command("sh", "./shell-scripts/uptime.sh").Output()
	if err != nil {
		log.Fatal(err)
	}
	_, err = ctx.EffectiveMessage.Reply(b, string(out), nil)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func devDebug(b *gotgbot.Bot, ctx *ext.Context) error {
	// get uptime of the bot by executing uptime command
	messageToSend := ""
	for _, arg := range ctx.Args() {
		messageToSend += arg + " "
	}
	_, err := ctx.EffectiveMessage.Reply(b, messageToSend, nil)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func apiIpLookup(b *gotgbot.Bot, ctx *ext.Context) error {
	// build the url (ctx.Args()[1] is the ip address)
	url := "https://ipapi.co/" + ctx.Args()[1] + "/json/"
	method := "GET"

	log.Printf("Requesting %s", url)

	// create a new http client
	client := &http.Client{}

	// build the request with the url and method
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// add headers
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "gotgbot")
	req.Header.Add("Connection", "keep-alive")

	// send the request
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)

	// read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// send the response body to the chat as a reply
	_, err = ctx.EffectiveMessage.Reply(b, string(body), nil)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil

}
