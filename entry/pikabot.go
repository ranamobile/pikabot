package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nlopes/slack"
	"github.com/ranamobile/pikabot"
)

// This will start up the pikabot and supports slash commands and
// parse messages with file attachments.  The file attachments will
// be automatically uploaded to the associated Google Drive account.
func main() {
	// Create the pikabot slash command handler and configure it to
	// listen on an HTTP endpoint on port 8080.
	pika := pikabot.CreatePikaSlash(
		os.Getenv("SLACK_SIGNING_SECRET"),
		fmt.Sprintf("/%s/pikascores", os.Getenv("PIKA_CONFIG_DIR")))
	http.HandleFunc("/", pika.SlashHandler)
	log.Println("[INFO] Server listening")
	go http.ListenAndServe(":8080", nil)

	// Configure the real-time messaging (RTM) connection with slack
	// and the message event handler (PikaDrive).
	log.Println("[INFO] RTM connection started")
	api := slack.New(os.Getenv("SLACK_OAUTH_TOKEN"), slack.OptionDebug(true))
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	pikaDrive := pikabot.PikaDrive{
		Client:        api,
		CredFilepath:  fmt.Sprintf("/%s/credentials.json", os.Getenv("PIKA_CONFIG_DIR")),
		TokenFilepath: fmt.Sprintf("/%s/token.json", os.Getenv("PIKA_CONFIG_DIR")),
	}

	// Listen for message events forever and parse the message
	// events with associated files to be copied to Google Drive.
	for msg := range rtm.IncomingEvents {
		log.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			log.Printf("Message: %v\n", ev)
			if len(ev.Files) > 0 {
				pikaDrive.CopyFilesToGdrive(ev)
			}
		}
	}
}
