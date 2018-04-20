package main

import (
	"context"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/orijtech/giphy/v1"
	"log"
	"os"
	"strings"
)

var token string
var slackApi *slack.Client
var giphyApi *giphy.Client

// Initializes api clients for slack and giphy
func init() {
	token = os.Getenv("SLACK_TOKEN")
	slackApi = slack.New(token)
	giphyApi, _ = giphy.NewClientFromEnvOrDefault()

}

func main() {

	slackApi.SetDebug(true)
	rtm := slackApi.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)

				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
					respond(rtm, ev, prefix)
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				//Take no action
			}
		}
	}
}

// The response message to an incoming message event
func respond(rtm *slack.RTM, msg *slack.MessageEvent, prefix string) {
	var response string
	text := msg.Text
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	requestSingleGif := map[string]bool{
		"gimme more":  true,
		"meow": true,
	}

	requestMultipleGifs := map[string]bool{
		"gimme more!":  true,
		"meow!": true,
	}

	if requestSingleGif[text] {
		gif := randomCat("cat")
		response = "Meow! " + gif.URL
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	}

	if requestMultipleGifs[text] {

		for i := 0; i < 5; i++ {
			gif := randomCat("cat")
			response = fmt.Sprintf("Meow #%d ! %s", i, gif.BitlyURL)
			rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
		}
	}
}

// Posts a random cat gif url from giphy to the channel
func randomCat(query string) *giphy.Giph {
	gif, err := giphyApi.RandomGIF(context.TODO(), &giphy.Request{
		Tag:    query,
		Rating: giphy.RatingPG,
	})
	if err != nil {
		log.Fatal(err)
	}
	return gif
}
