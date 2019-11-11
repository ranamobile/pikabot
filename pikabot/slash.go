package pikabot

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"unicode"

	"github.com/nlopes/slack"
)

type PikaSlash struct {
	SlackSigningSecret string
	ScoreFilepath      string
}

func CreatePikaSlash(secret string, scorepath string) PikaSlash {
	return PikaSlash{
		SlackSigningSecret: secret,
		ScoreFilepath:      scorepath,
	}
}

func (pika *PikaSlash) ParseMarkCommand(command slack.SlashCommand) *slack.Msg {
	response := []rune(command.Text)

	for index, character := range response {
		if rand.Int()%2 == 0 {
			response[index] = unicode.ToUpper(character)
		} else {
			response[index] = unicode.ToLower(character)
		}
	}

	return &slack.Msg{
		Text:         string(response),
		ResponseType: slack.ResponseTypeInChannel,
	}
}

func (pika *PikaSlash) ParseScoreCommand(command slack.SlashCommand) *slack.Msg {
	scorelist := CreateScoreList()
	scorelist.Read(pika.ScoreFilepath)

	name := command.Text[:len(command.Text)-2]
	count := command.Text[len(command.Text)-2:]
	text := "invalid command"

	if command.Text == "top" {
		sort.Sort(sort.Reverse(scorelist))

		var topScores []string
		for index := 0; index < 3; index++ {
			topScores = append(topScores, fmt.Sprintf("%s: %d", scorelist.Scores[index].Name, scorelist.Scores[index].Count))
		}
		text = strings.Join(topScores, "\n")
	} else if count == "++" {
		score := scorelist.Increment(name)
		text = fmt.Sprintf("%s: %d", score.Name, score.Count)
	} else if count == "--" {
		score := scorelist.Decrement(name)
		text = fmt.Sprintf("%s: %d", score.Name, score.Count)
	}

	log.Println("Writing scores to file:", pika.ScoreFilepath)
	scorelist.Write(pika.ScoreFilepath)

	return &slack.Msg{
		Text:         text,
		ResponseType: slack.ResponseTypeInChannel,
	}
}

func (pika *PikaSlash) SlashHandler(response http.ResponseWriter, request *http.Request) {
	verifier, err := slack.NewSecretsVerifier(request.Header, pika.SlackSigningSecret)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	request.Body = ioutil.NopCloser(io.TeeReader(request.Body, &verifier))
	command, err := slack.SlashCommandParse(request)
	if err != nil {
		log.Println("Slack error:", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		log.Println("Slack error:", err)
		response.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Println("Command:", command.Command, command.Text)
	var params *slack.Msg

	switch command.Command {
	case "/mark":
		params = pika.ParseMarkCommand(command)
	case "/score":
		params = pika.ParseScoreCommand(command)
	}

	if params != nil {
		data, err := json.Marshal(params)
		if err == nil {
			response.Header().Set("Content-Type", "application/json")
			response.Write(data)
			return
		}
	}

	response.WriteHeader(http.StatusInternalServerError)
}
