package main

import (
	"log"
	"fmt"
	"strings"
	"github.com/james-bowman/jeannie"
	"github.com/james-bowman/sentence"
	"github.com/james-bowman/slack"
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

func processMessage(message *slack.Message)  {
	var response string
	
	log.Printf("%s-> %s", message.From, message.Text)
	
	var sentences []string
	for key, _ := range answerers {
		sentences = append(sentences, key)
	}
	
	answer, err := sentence.Recognise(sentenceRecognitionToken, message.Text, sentences)
		
	if answer != "" {
		response = answerers[answer](message.Text)
	}
	
	if response == "" {
		if err == nil {
			// if sentence not matched to a supported request then fallback to asking Jeannie
			response, err = jeannie.AskQuestion(jeannieToken, message.Text)
		}
	}
	
	if err != nil {
		response = "I seem to be having a problem: " + err.Error()
	}
	
	log.Printf("me-> %s", response)
	
	message.Respond(response)
}