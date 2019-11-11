package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nlopes/slack"
	"pikabot"
)

func main() {
	pika := pikabot.CreatePikaSlash(os.Getenv("SLACK_SIGNING_SECRET"), os.Getenv("SLACK_PIKA_SCOREFILE"))
	http.HandleFunc("/", pika.SlashHandler)
	log.Println("[INFO] Server listening")
	go http.ListenAndServe(":8080", nil)

	log.Println("[INFO] RTM connection started")
	api := slack.New(os.Getenv("SLACK_OAUTH_TOKEN"), slack.OptionDebug(true))
	pikaDrive := pikabot.PikaDrive{
		Client: api,
		CredFilepath: fmt.Sprintf("/%s/credentials.json", os.Getenv("PIKA_CONFIG_DIR")),
		TokenFilepath: fmt.Sprintf("/%s/token.json", os.Getenv("PIKA_CONFIG_DIR")),
	}

	channel, err := api.GetGroupInfo("GBVS0KZ39") //channel
	if err != nil {
		log.Println("[DEBUG] Failed to get channel info", err)
	} else {
		log.Println("[DEBUG] Channel name:", channel.Name)
	}

	rtm := api.NewRTM()
	go rtm.ManageConnection()
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
