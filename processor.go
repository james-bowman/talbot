package main

import (
	"log"
	"fmt"
	"strings"
	"github.com/james-bowman/slack"
	"github.com/james-bowman/jeannie"
	"github.com/james-bowman/sentence"
	"regexp"
	"strconv"
)

type answerer func(string) string

var answerers map[string]answerer

func Init() {
	err := InitRooms()
	
	if err != nil {
		log.Printf("Error initialising rooms: " + err.Error())
	}

	answerers = map[string]answerer{
    	"ping": ping,
    	"hello": greeting,
    	"move card 153 to done": minglemovecard,
    	"help": help,
    	"where is the IT Training Room": findRoom,
    	"the rules": rules,
    	"which rooms can seat 6 people and have a whiteboard": findRoomsByFeatures,
	}
}

func rules(dummy string) string {
	return "1. A robot may not injure a human being or, through inaction, allow a human being to come to harm.\n" +
		"2. A robot must obey the orders given it by human beings, except where such orders would conflict with the First Law.\n" +
		"3. A robot must protect its own existence as long as such protection does not conflict with the First or Second Law."
}

func ping(pingmsg string) string {
	return "pong"
}

func greeting(greetingmsg string) string {
	return greetingmsg
}

func minglemovecard(instruction string) string {
	return "Sorry I can't quite do this yet but it is coming soon."
}

func help(dummy string) string {
	response := "I am a robot that listens to the team's chat and provides automated functions." +
	"  I currently support the following commands:\n"
	
	for key, _ := range answerers {
		response = fmt.Sprintf("%s\n\t%s", response, key)
	}
	
	return response + "\n\nI can do some other things too - try asking me something!"
}

func findRoom(instruction string) string {
	var response string
	
	for _, eachRoom := range Rooms {
		if strings.Contains(strings.ToUpper(instruction), strings.ToUpper(eachRoom.Name)) {
			if len(response) > 0 {
				response = response + "\n\n"
			} 
			response = response + String(eachRoom)
		}
	}

	return response
}

func findRoomsByFeatures(instruction string) string {
	reSeats, err := regexp.Compile(`([0-9]+)`)
	
	if err != nil {
		log.Printf("Error compiling Regex to find seats: %s", err.Error())
	}
	
	seats := reSeats.FindStringSubmatch(instruction)
		
	var size int
	if len(seats) > 0 {
		size, err = strconv.Atoi(seats[1])
		if err != nil {
			log.Printf("Error converting %s number of seats string to int: %s", seats[1], err.Error())
			return ""
		}
	}

	vc := strings.Contains(strings.ToUpper(instruction), "VC")
	whiteboard := strings.Contains(strings.ToUpper(instruction), "WHITEBOARD")
	screen := strings.Contains(strings.ToUpper(instruction), "SCREEN")	
					
	var response string
	for _, eachRoom := range Rooms {
		if eachRoom.Seats >= size {
			if (!vc || (vc && eachRoom.VC)) && 
						(!whiteboard || (whiteboard && eachRoom.Whiteboard)) && 
						(!screen || (screen && eachRoom.Screen)) {
				if len(response) > 0 {
					response = response + "\n\n"
				} 
				response = response + String(eachRoom)
			}
		}
	}
	
	if response == "" {
		response = "I am afraid there aren't any in this building"
	}
	
	return response
}

func processMessage(user string, question string, replyto string) string {
	var response string
	
	log.Printf("%s-> %s", user, question)
	
	var sentences []string
	for key, _ := range answerers {
		sentences = append(sentences, key)
	}
	
	answer, err := sentence.Recognise(sentenceRecognitionToken, question, sentences)
		
	if answer != "" {
		response = answerers[answer](question)
	}
	
	if response == "" {
		if err == nil {
			// if sentence not matched to a supported request then fallback to asking Jeannie
			response, err = jeannie.AskQuestion(jeannieToken, question)
		}
	}
	
	if err != nil {
		response = "I seem to be having a problem: " + err.Error()
	}
	
	fullResponse := replyto + response
	log.Printf("me-> %s", fullResponse)
	
	return fullResponse
}

func filterMessage(data slack.Event) *slack.Event {
	var response string
	
	if data.Type == "message" && data.ReplyTo == 0 {
		if data.Channel[0] == 'D' {
			// process direct messages
			response = processMessage(data.User, data.Text, "")
		} else {
			// process messages in public channels directed at Talbot
			r, _ := regexp.Compile("^(<@U03K60XR7>|@?talbot):? (.+)")
				
			matches := r.FindStringSubmatch(data.Text)
				
			if len(matches) == 3 {
				response = processMessage(data.User, matches[2], "<@" + data.User + ">: ")
			} else {
				if strings.Contains(strings.ToUpper(data.Text), "BATCH") {
					response = "<@" + data.User + ">: Language Error: I don't understand the term 'Batch' please re-state using goal orientated language"
				}
			}
		}
	}
	if response != "" {
		return &slack.Event{Type: "message", Channel: data.Channel, Text: response}
	}
	return nil
}