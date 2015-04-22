package main

import (
	"log"
	"fmt"
	"io/ioutil"
	"github.com/james-bowman/slack"
	"github.com/james-bowman/talbot/brain"
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
