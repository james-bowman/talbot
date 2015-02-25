package main

import (
	"log"
	"github.com/james-bowman/slack"
)

const (
//  James' slack token
//	slackJamesToken = "xoxp-3028015569-3028015571-3620637852-ebd10f"

	slackToken = "xoxb-3652031857-rWSOPDYm7p0WzzMwSq2ECqiC"
	
	jeannieToken = "g4hcu5vxqcmshoiWPwKWrhFOX1j8p1rd8FOjsngsVHYvkCBQYv"
	
	sentenceRecognitionToken = jeannieToken
	
)

func main() {
	Init()

	conn, err := slack.Connect(slackToken)
	
	if err != nil {
		log.Fatal(err)
	}
		
	slack.EventProcessor(conn, processMessage)	
}
