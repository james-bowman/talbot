package main

import (
	"log"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

const (
	Meeting = "Meeting"
	Quiet = "Quiet"
)

type Room struct {
	Floor int			`json:"floor"`
	Ref string			`json:"ref"`
	Name string			`json:"name"`
	Type string			`json:"type"`
	Position string		`json:"position"`
	Bookable bool		`json:"bookable"`
	Configurable bool	`json:"configurable"`
	Seats int			`json:"seats"`
	Screen bool			`json:"screen"`
	VC bool				`json:"vc"`
	Camera bool			`json:"camera"`
	Whiteboard bool		`json:"whiteboard"`
}

var Rooms []Room

func InitRooms() error {
	roomFile, err := ioutil.ReadFile("rooms.json")
	if err != nil {
		log.Printf("Error opening rooms file: %s", err.Error())
		return err
	}
	
	err = json.Unmarshal(roomFile, &Rooms)
	
	return err
}

func String(room Room) string {
	text := fmt.Sprintf("%s is on the *%dth floor* on the *%s* side of the building.\n", room.Name, room.Floor, room.Position)
	text = fmt.Sprintf("%sIt has *%d seats*", text, room.Seats)

	if room.Screen {
		text = text + ", a screen"
	}
	
	if room.VC {
		text = text + ", VC facilities"
	}
	
	if room.Whiteboard {
		text = text + ", a whiteboard"
	}
	
	return text
}




