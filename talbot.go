package main

import (
	"log"
	"strings"
	"encoding/json"
	"github.com/james-bowman/slack"
	"github.com/james-bowman/jeannie"
)


func processMessage(data slack.Event, sequence int) ([]byte, error) {
	if data.Type == "message" {
			if data.Channel[0] == 'D' {
				// process direct messages
				sequence++
				response, err := jeannie.AskQuestion(data.Text)
				if err != nil {
					response = err.Error()
				}
				e := slack.Event{Id: sequence, Type: "message", Channel: data.Channel, Text: response}
				bytes, err := json.Marshal(e)
				if err != nil {
					return nil, err
				}
				return bytes, nil
			} else {
				// process messages in open channels directed at Talbot
				if strings.HasPrefix(data.Text, "<@U03K60XR7>:") {
					sequence++
					
					response, err := jeannie.AskQuestion(data.Text[13:])
					if err != nil {
						response = err.Error()
					}
					
					e := slack.Event{Id: sequence, Type: "message", Channel: data.Channel, Text: "<@" + data.User + ">: " + response}
					bytes, err := json.Marshal(e)
					if err != nil {
						return nil, err
					}
					return bytes, nil
				}
			}
		}
	return nil, nil
}

func main() {
//  James' token
//	token := "xoxp-3028015569-3028015571-3620637852-ebd10f"

// 	Bot token
	token := "xoxb-3652031857-rWSOPDYm7p0WzzMwSq2ECqiC"

	for {
		conn, err := slack.Connect(token)
	
		if err != nil {
			log.Fatal(err)
		}
		
		conn.Start(processMessage)	
	}
}
