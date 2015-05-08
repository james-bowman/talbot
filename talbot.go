package main // import "github.com/james-bowman/talbot"

import (
	"flag"
	"fmt"
	"github.com/james-bowman/slack"
	"github.com/james-bowman/talbot/brain"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	slackToken := getToken()

	conn, err := slack.Connect(slackToken)

	if err != nil {
		log.Fatal(err)
	}

	slack.EventProcessor(conn, brain.OnAskedMessage, brain.OnHeardMessage)
}

// getToken for authenticating with Slack.  Ordered lookup process trying first the command line,
// then environment variables and finally a config file
func getToken() string {
	var slackToken string

	// look for the token on the command line
	log.Println("Looking for Slack auth token as command line argument")
	flag.StringVar(&slackToken, "slacktoken", "", "Slack authentication token - if not specified, will look for an environment variable or config file")
	flag.Parse()

	if slackToken == "" {
		// if not specified look for it in an environment variable
		log.Println("Slack auth token not found - looking for environment variable")
		slackToken = os.Getenv("TALBOT_SLACK_TOKEN")
	}

	if slackToken == "" {
		// if not specified look for it in a config file
		log.Println("Slack auth token not found - looking for config file")
		slackTokenFileName := "slack.token"

		slackTokenFile, err := ioutil.ReadFile(slackTokenFileName)
		if err != nil {
			log.Panic(fmt.Sprintf("Error opening slack authentication token file %s: %s", slackTokenFileName, err))
		}

		slackToken = string(slackTokenFile)
	}

	if slackToken != "" {
		log.Println("Slack auth token found")
	}

	return slackToken
}
