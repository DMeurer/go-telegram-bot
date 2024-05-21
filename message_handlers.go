package main

import (
	"encoding/json"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"io"
	"log"
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

type ipLookupApiResponse struct {
	Ip                   string
	Network              string
	Version              string
	City                 string
	Region               string
	RegionCode           string
	Country              string
	Country_Name         string
	Country_Code         string
	Country_Code_Iso3    string
	Country_Capital      string
	Country_Tld          string
	Continent_Code       string
	In_Eu                string
	Postal               string
	Latitude             json.Number
	Longitude            json.Number
	Timezone             string
	Utc_Offset           string
	Country_Calling_Code string
	Currency             string
	Currency_Name        string
	Languages            string
	Country_Area         string
	Country_Population   string
	Asn                  string
	Org                  string
}

func apiIpLookup(b *gotgbot.Bot, ctx *ext.Context) error {
	// check if ctx.Args() is longer than 1
	if len(ctx.Args()) < 2 {
		_, err := ctx.EffectiveMessage.Reply(b, "Please provide an IP address to lookup", nil)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		return nil
	}

	// build the url (ctx.Args()[1] is the ip address)
	url := "https://ipapi.co/" + ctx.Args()[1] + "/json/"
	method := "GET"
	headers := []headerEntry{
		{"Accept", "*/*"},
		{"User-Agent", "gotgbot"},
		{"Connection", "keep-alive"},
	}

	log.Printf("Requesting %s", url)

	// send the request
	res, err := sendRequest(method, url, headers)

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

	// parse the response body
	var ipLookupResponse ipLookupApiResponse
	err = json.Unmarshal(body, &ipLookupResponse)
	if err != nil {
		log.Print("Failed to parse response")
		log.Print(err)
		_, err := ctx.EffectiveMessage.Reply(b, "Failed to parse response\n"+string(body), nil)
		if err != nil {
			panic("failed to send message: " + err.Error())
		}
	}

	// build the response message
	responseMessage := ""

	if len(ctx.Args()) == 2 {
		responseMessage += fmt.Sprintf(
			"Requested IP: %s\nCity: %s\nRegion: %s\nCountry: %s\nTimezone: %s\nLanguages: %s\nASN: %s\nORG: %s",
			ipLookupResponse.Ip,
			ipLookupResponse.City,
			ipLookupResponse.Region,
			ipLookupResponse.Country,
			ipLookupResponse.Timezone,
			ipLookupResponse.Languages,
			ipLookupResponse.Asn,
			ipLookupResponse.Org,
		)
		responseMessage += fmt.Sprintf("\n\nTo get more information, use /ip-address <ip-address> *")
	} else if len(ctx.Args()) == 3 {
		responseMessage += fmt.Sprintf(
			"Requested IP: %s\nNetwork: %s\nVersion: %s\nCity: %s\nRegion: %s\nRegion Code: %s\nCountry: %s\nCountry Name: %s\nCountry Code: %s\nCountry Code ISO3: %s\nCountry Capital: %s\nCountry TLD: %s\nContinent Code: %s\nIn EU: %s\nPostal: %s\nLatitude: %s\nLongitude: %s\nTimezone: %s\nUTC Offset: %s\nCountry Calling Code: %s\nCurrency: %s\nCurrency Name: %s\nLanguages: %s\nCountry Area: %s\nCountry Population: %s\nASN: %s\nORG: %s",
			ipLookupResponse.Ip,
			ipLookupResponse.Network,
			ipLookupResponse.Version,
			ipLookupResponse.City,
			ipLookupResponse.Region,
			ipLookupResponse.RegionCode,
			ipLookupResponse.Country,
			ipLookupResponse.Country_Name,
			ipLookupResponse.Country_Code,
			ipLookupResponse.Country_Code_Iso3,
			ipLookupResponse.Country_Capital,
			ipLookupResponse.Country_Tld,
			ipLookupResponse.Continent_Code,
			ipLookupResponse.In_Eu,
			ipLookupResponse.Postal,
			ipLookupResponse.Latitude,
			ipLookupResponse.Longitude,
			ipLookupResponse.Timezone,
			ipLookupResponse.Utc_Offset,
			ipLookupResponse.Country_Calling_Code,
			ipLookupResponse.Currency,
			ipLookupResponse.Currency_Name,
			ipLookupResponse.Languages,
			ipLookupResponse.Country_Area,
			ipLookupResponse.Country_Population,
			ipLookupResponse.Asn,
			ipLookupResponse.Org,
		)
	} else {
		responseMessage += fmt.Sprintf("Invalid number of arguments.\n\nUsage: /ip-address <ip-address> [more-info]")
	}

	// send the response body to the chat as a reply
	_, err = ctx.EffectiveMessage.Reply(b, responseMessage, nil)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}
