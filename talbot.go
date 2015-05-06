package main // import "github.com/james-bowman/talbot"

import (
	"fmt"
	"github.com/james-bowman/slack"
	"github.com/james-bowman/talbot/brain"
	"io/ioutil"
	"log"
)

func main() {
	slackTokenFileName := "slack.token"

	slackToken, err := ioutil.ReadFile(slackTokenFileName)
	if err != nil {
		log.Panic(fmt.Sprintf("Error opening slack authentication token file %s: %s", slackTokenFileName, err))
	}

	conn, err := slack.Connect(string(slackToken))

	if err != nil {
		log.Fatal(err)
	}

	slack.EventProcessor(conn, brain.OnAskedMessage, brain.OnHeardMessage)
}
