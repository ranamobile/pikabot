package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nlopes/slack"
	"github.com/ranamobile/pikabot"
)

func CopyFilesToGdrive(api *slack.Client, event *slack.MessageEvent) {
	var channelName string
	channel, err := api.GetChannelInfo(event.Channel)
	if err == nil {
		channelName = channel.Name
	} else {
		group, err := api.GetGroupInfo(event.Channel)
		if err != nil {
			log.Println("[ERROR] Failed to get channel name:", err)
			return
		}
		channelName = group.Name
	}

	// connect to google drive
	gdrive := lib.CreateGoogleService(os.Getenv("SLACK_PIKA_CREDFILE"), os.Getenv("SLACK_PIKA_TOKENFILE"))
	slackDir, err := lib.CreateDir(gdrive, "Slack", "root")
	if err != nil {
		log.Println("[ERROR] Failed to create/get directory:", err)
		return
	}
	channelDir, err := lib.CreateDir(gdrive, channelName, slackDir.Id)
	if err != nil {
		log.Println("[ERROR] Failed to create/get directory:", err)
		return
	}

	for index, item := range event.Files {
		filename := fmt.Sprintf("%s-%d.jpg", event.Timestamp, index)
		filepath := fmt.Sprintf("/tmp/%s", filename)
		handler, err := os.Create(filepath)
		if err != nil {
			log.Println("[ERROR] Failed to save image: %s", err)
			return
		}
		api.GetFile(item.URLPrivateDownload, handler)
		handler.Close()

		handler, err = os.Open(filepath)
		if err != nil {
			log.Println("[ERROR] Failed to read image: %s", err)
			return
		}
		lib.CreateFile(gdrive, filename, "image/jpeg", handler, channelDir.Id)
	}
}

func main() {
	pika := lib.CreatePikaSlash(os.Getenv("SLACK_SIGNING_SECRET"), os.Getenv("SLACK_PIKA_SCOREFILE"))
	http.HandleFunc("/", pika.SlashHandler)
	log.Println("[INFO] Server listening")
	go http.ListenAndServe(":8080", nil)

	log.Println("[INFO] RTM connection started")
	api := slack.New(os.Getenv("SLACK_OAUTH_TOKEN"), slack.OptionDebug(true))

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
				CopyFilesToGdrive(api, ev)
			}
		}
	}

}
