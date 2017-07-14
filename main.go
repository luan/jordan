package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

func main() {
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(false)

	rtm := api.NewRTM()
	go rtm.ManageConnection() // spawn slack bot

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				if ev.Msg.Type == "message" && ev.Msg.SubType != "message_deleted" && strings.Contains(ev.Msg.Text, "<!here|@here>") {
					reply := fmt.Sprintf("Hello <@%s>, please avoid using `@here` in this channel. Please mention the team member(s) listed as interrupt in the channel topic instead.", ev.Msg.User)
					rtm.SendMessage(rtm.NewOutgoingMessage(reply, ev.Msg.Channel))
				}

			case *slack.RTMError:
				logger.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				logger.Println("Invalid credentials")
				break Loop

			default:
			}
		}
	}
}
