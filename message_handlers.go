package main

import (
	"encoding/json"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"io"
	"log"
	"main/requests"
	"main/systemInfo"
	"main/tools"
	"os/exec"
	"strconv"
	"time"
)

func help(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, "Available commands:\n"+
		"/help - Show this message\n"+
		"/echo - Repeat the message\n"+
		"/ping - Show Ping of the bot\n"+
		"/uptime - Get the uptime of the bot and server\n"+
		"/ip-address - Get information about an IP address\n"+
		"/mensa - Get the meals of the canteen", nil)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

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
	out, err := exec.Command("sh", "./systemInfo/shell-scripts/ping.sh").Output()
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
	if len(ctx.Args()) > 2 {
		_, err := ctx.EffectiveMessage.Reply(b, "Too many arguments provided", nil)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		return nil
	}

	if len(ctx.Args()) <= 1 {
		serverUptime, err := systemInfo.GetUptime("server")
		if err != nil {
			log.Fatal(err)
		}
		containerUptime, err := systemInfo.GetUptime("container")
		if err != nil {
			log.Fatal(err)
		}
		_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Server uptime: %s\nContainer uptime: %s", systemInfo.FormatDurationHumanReadable(serverUptime), systemInfo.FormatDurationHumanReadable(containerUptime)), nil)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		return nil
	}

	if len(ctx.Args()) == 2 {
		switch ctx.Args()[1] {
		case "server":
			serverUptime, err := systemInfo.GetUptime("server")
			if err != nil {
				log.Fatal(err)
			}
			_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Server uptime: %s", systemInfo.FormatDurationHumanReadable(serverUptime)), nil)
			if err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
			return nil
		case "container":
			containerUptime, err := systemInfo.GetUptime("container")
			if err != nil {
				log.Fatal(err)
			}
			_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Container uptime: %s", systemInfo.FormatDurationHumanReadable(containerUptime)), nil)
			if err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
			return nil
		default:
			_, err := ctx.EffectiveMessage.Reply(b, "Invalid service provided", nil)
			if err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
			return nil
		}
	}

	_, err := ctx.EffectiveMessage.Reply(b, "Please provide a service to get the uptime of", nil)
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
	Ip                 string  `json:"ip"`
	Network            string  `json:"network"`
	Version            string  `json:"version"`
	City               string  `json:"city"`
	Region             string  `json:"region"`
	RegionCode         string  `json:"region_code"`
	Country            string  `json:"country"`
	CountryName        string  `json:"country_name"`
	CountryCode        string  `json:"country_code"`
	CountryCodeIso3    string  `json:"country_code_iso3"`
	CountryCapital     string  `json:"country_capital"`
	CountryTld         string  `json:"country_tld"`
	ContinentCode      string  `json:"continent_code"`
	InEu               bool    `json:"in_eu"`
	Postal             string  `json:"postal"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Timezone           string  `json:"timezone"`
	UtcOffset          string  `json:"utc_offset"`
	CountryCallingCode string  `json:"country_calling_code"`
	Currency           string  `json:"currency"`
	CurrencyName       string  `json:"currency_name"`
	Languages          string  `json:"languages"`
	CountryArea        float64 `json:"country_area"`
	CountryPopulation  int     `json:"country_population"`
	Asn                string  `json:"asn"`
	Org                string  `json:"org"`
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
	headers := []requests.HeaderEntry{
		{"Accept", "*/*"},
		{"User-Agent", "gotgbot"},
		{"Connection", "keep-alive"},
	}

	log.Printf("Requesting %s", url)

	// send the request
	res, err := requests.SendRequest(method, url, headers)
	if err != nil {
		log.Printf("Failed to send request\n %s", err)
		_, err := ctx.EffectiveMessage.Reply(b, "Failed to send request", nil)
		if err != nil {
			panic("failed to send message: " + err.Error())
		}
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
			ipLookupResponse.CountryName,
			ipLookupResponse.CountryCode,
			ipLookupResponse.CountryCodeIso3,
			ipLookupResponse.CountryCapital,
			ipLookupResponse.CountryTld,
			ipLookupResponse.ContinentCode,
			strconv.FormatBool(ipLookupResponse.InEu),
			ipLookupResponse.Postal,
			strconv.FormatFloat(ipLookupResponse.Latitude, 'f', -1, 64),
			strconv.FormatFloat(ipLookupResponse.Longitude, 'f', -1, 64),
			ipLookupResponse.Timezone,
			ipLookupResponse.UtcOffset,
			ipLookupResponse.CountryCallingCode,
			ipLookupResponse.Currency,
			ipLookupResponse.CurrencyName,
			ipLookupResponse.Languages,
			strconv.FormatFloat(ipLookupResponse.CountryArea, 'f', -1, 64),
			strconv.Itoa(ipLookupResponse.CountryPopulation),
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

func mensa(b *gotgbot.Bot, ctx *ext.Context) error {

	if tools.StringInSlice("-help", ctx.Args()) || tools.StringInSlice("-h", ctx.Args()) || tools.StringInSlice("--help", ctx.Args()) || tools.StringInSlice("-?", ctx.Args()) {
		_, err := ctx.EffectiveMessage.Reply(b, "Usage: /mensa [-all] [-mon] [-tue] [-wed] [-thu] [-fri] [-sat] [-wee]\n\nWill show the menu of this weekday\n\n-all: Show all days\n-mon: Show Monday\n-tue: Show Tuesday\n-wed: Show Wednesday\n-thu: Show Thursday\n-fri: Show Friday\n-sat: Show Saturday\n-wee: Show weekdays (Monday - Friday)", nil)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		return nil
	}

	var days []string

	if len(ctx.Args()) == 1 {
		days = append(days, time.Wednesday.String())
	} else {

		if tools.StringInSlice("-mon", ctx.Args()) {
			days = append(days, "Monday")
		}
		if tools.StringInSlice("-tue", ctx.Args()) {
			days = append(days, "Tuesday")
		}
		if tools.StringInSlice("-wed", ctx.Args()) {
			days = append(days, "Wednesday")
		}
		if tools.StringInSlice("-thu", ctx.Args()) {
			days = append(days, "Thursday")
		}
		if tools.StringInSlice("-fri", ctx.Args()) {
			days = append(days, "Friday")
		}
		if tools.StringInSlice("-sat", ctx.Args()) {
			days = append(days, "Saturday")
		}
		if tools.StringInSlice("-all", ctx.Args()) {
			days = append(days, "Monday")
			days = append(days, "Tuesday")
			days = append(days, "Wednesday")
			days = append(days, "Thursday")
			days = append(days, "Friday")
			days = append(days, "Saturday")
		}
		if tools.StringInSlice("-wee", ctx.Args()) {
			days = append(days, "Monday")
			days = append(days, "Tuesday")
			days = append(days, "Wednesday")
			days = append(days, "Thursday")
			days = append(days, "Friday")
		}
	}

	// prepare the message
	message := tools.PrepareMessageForMeals(tools.LoadMeals(), false, days...)

	// send the message
	_, err := ctx.EffectiveMessage.Reply(b, message, nil)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}
