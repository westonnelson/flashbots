// Discord webhook helpers
// https://discord.com/developers/docs/resources/webhook#execute-webhook
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/metachris/flashbots/blockcheck"
)

type DiscordWebhookPayload struct {
	Content string `json:"content"`
}

func SendBundleOrderErrorToDiscord(b *blockcheck.BlockCheck) error {
	msg := b.Sprint(false, true)
	return SendToDiscord(msg)
}

func SendToDiscord(msg string) error {
	url := os.Getenv("DISCORD_WEBHOOK")
	if len(url) == 0 {
		return errors.New("no DISCORD_WEBHOOK env variable found")
	}

	// Trim msg to 2k characters if necessary
	if len(msg) > 2000 {
		if strings.Contains(msg, "```") {
			msg = msg[0:1994] + "...```"
		} else {
			msg = msg[0:1997] + "..."
		}
	}

	discordPayload := DiscordWebhookPayload{Content: msg}
	payloadBytes, err := json.Marshal(discordPayload)
	if err != nil {
		return err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	defer res.Body.Close()
	log.Println("response status:", res.Status)

	if res.StatusCode >= 300 {
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		bodyString := string(bodyBytes)
		log.Println(bodyString)
	}
	return nil
}
