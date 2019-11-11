package pikabot

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/nlopes/slack"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type PikaDrive struct {
	Client *slack.Client
	CredFilepath string
	TokenFilepath string
}

func CreateGoogleService(credential string, token string) (service *drive.Service) {
	content, err := ioutil.ReadFile(credential)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(content, drive.DriveFileScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := GetClient(config, token)

	service, err = drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	return service
}

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config, tokenFile string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := TokenFromFile(tokenFile)
	if err != nil {
		tok = GetTokenFromWeb(config)
		SaveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func SaveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func CreateDir(service *drive.Service, name string, parentId string) (*drive.File, error) {
	result, err := service.Files.List().
		Q(fmt.Sprintf("name='%s' and '%s' in parents", name, parentId)).
		Spaces("drive").
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	if len(result.Files) > 0 {
		return result.Files[0], nil
	} else {
		driveFile := &drive.File{
			Name:     name,
			MimeType: "application/vnd.google-apps.folder",
			Parents:  []string{parentId},
		}

		file, err := service.Files.Create(driveFile).Do()

		if err != nil {
			log.Println("Could not create dir: " + err.Error())
			return nil, err
		}
		return file, nil
	}
}

func CreateFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
	driveFile := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := service.Files.Create(driveFile).Media(content).Do()

	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	return file, nil
}

func (pika *PikaDrive) CopyFilesToGdrive(event *slack.MessageEvent) {
	var channelName string
	channel, err := pika.Client.GetChannelInfo(event.Channel)
	if err == nil {
		channelName = channel.Name
	} else {
		group, err := pika.Client.GetGroupInfo(event.Channel)
		if err != nil {
			log.Println("[ERROR] Failed to get channel name:", err)
			return
		}
		channelName = group.Name
	}

	// connect to google drive
	gdrive := CreateGoogleService(pika.CredFilepath, pika.TokenFilepath)
	slackDir, err := CreateDir(gdrive, "Slack", "root")
	if err != nil {
		log.Println("[ERROR] Failed to create/get directory:", err)
		return
	}
	channelDir, err := CreateDir(gdrive, channelName, slackDir.Id)
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
		pika.Client.GetFile(item.URLPrivateDownload, handler)
		handler.Close()

		handler, err = os.Open(filepath)
		if err != nil {
			log.Println("[ERROR] Failed to read image: %s", err)
			return
		}
		CreateFile(gdrive, filename, "image/jpeg", handler, channelDir.Id)
	}
}
