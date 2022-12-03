// Copyright (c) 2022 Tailscale Inc & AUTHORS All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var sess *discordgo.Session

func init() {
	godotenv.Load()
}

type incomingWebhook struct {
	Timestamp time.Time         `json:"timestamp"`
	Version   int               `json:"version"`
	Type      string            `json:"type"`
	Tailnet   string            `json:"tailnet"`
	Message   string            `json:"message"`
	Data      map[string]string `json:"data"`
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received webhook")
	secret := os.Getenv("TS_WEBHOOK_SECRET")
	events, err := verifyWebhookSignature(r, secret)
	if err != nil {
		fmt.Printf("handleWebhook verifyWebhookSignature: %v\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	fmt.Printf("handleWebhook received %d events\n", len(events))
	for _, event := range events {
		fmt.Printf("handleWebhook event: %+v\n", event)

		embed := &discordgo.MessageEmbed{}

		// Take embed title from the event message
		embed.Title = event.Message

		// Add the event data to the embed description
		embed.Description = ""

		for key, value := range event.Data {
			embed.Description += fmt.Sprintf("**%s**: %s\n", key, value)
		}

		// Add the event timestamp to the embed timestamp
		embed.Timestamp = event.Timestamp.Format(time.RFC3339)

		// Add the event tailnet to the embed fields
		embed.Fields = []*discordgo.MessageEmbedField{
			{
				Name:  "Type",
				Value: event.Type,
			},
			{
				Name:  "Tailnet",
				Value: event.Tailnet,
			},
		}

		// Send the embed to the Discord channel
		_, err := sess.ChannelMessageSendEmbed(os.Getenv("DISCORD_CHANNEL_ID"), embed)

		if err != nil {
			fmt.Printf("handleWebhook ChannelMessageSendEmbed: %v\n", err)
		}
	}
}

func main() {
	var err error
	sess, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))

	if err != nil {
		log.Fatal(err)
	}

	// Open a WS connection to Discord and do an identify so we can send messages.
	err = sess.Open()

	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	log.Printf("Listening for webhooks on port %s...\n", port)
	http.HandleFunc("/", handleWebhook)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
