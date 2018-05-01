package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
	"github.com/orijtech/giphy/v1"
)

var (
	token    string
	slackApi *slack.Client
	giphyApi *giphy.Client
)

// Initializes api clients for slack and giphy
func init() {
	var err error
	token, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		log.Panicf("'SLACK_TOKEN' is not set")
	}
	slackApi = slack.New(token)
	giphyApi, err = giphy.NewClientFromEnvOrDefault()
	if err != nil {
		log.Panicf("Can not init giphy client: %s", err)
	}
}

func routeMessage(msg slack.RTMEvent, rtm *slack.RTM) error {
	fmt.Print("Event Received: %+v", msg)

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
		return errors.New("Invalid credentials")

	default:
		//Take no action
	}
	return nil
}

func main() {

	slackApi.SetDebug(true)
	rtm := slackApi.NewRTM()
	go rtm.ManageConnection()
	var err error
	for err == nil {
		err = routeMessage(<-rtm.IncomingEvents, rtm)
	}
}

func unifyString(text, prefix string) string {
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	return strings.ToLower(text)
}

// The response message to an incoming message event
func respond(rtm *slack.RTM, msg *slack.MessageEvent, prefix string) {
	var response string
	text := unifyString(msg.Text, prefix)

	requestSingleGif := map[string]bool{
		"gimme more": true,
		"meow":       true,
	}

	requestMultipleGifs := map[string]bool{
		"gimme more!": true,
		"meow!":       true,
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
