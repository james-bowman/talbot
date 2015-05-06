package main // import "github.com/james-bowman/talbot"

import (
	"github.com/james-bowman/talbot/brain"
	"math/rand"
	"regexp"
	"time"
)

func init() {
	brain.Register(brain.Action{
		Regex:       regexp.MustCompile("(?i)ping"),
		Usage:       "ping",
		Description: "Reply with pong",
		Answerer: func(dummy string) string {
			return "pong"
		},
	})

	brain.Register(brain.Action{
		Regex:       regexp.MustCompile("(?i)^(hello|hi|hiya|howdy|bonjour|bon dia|hallo|salut|aloha|hola|hey|yo)($| )"),
		Usage:       "hello",
		Description: "Reply with a random greeting",
		Answerer:    greeting,
	})

	brain.Register(brain.Action{
		Regex:       regexp.MustCompile("(?i)rules"),
		Usage:       "the rules",
		Description: "Make sure I know the rules",
		Answerer:    rules,
	})
}

func rules(dummy string) string {
	return "1. A robot may not injure a human being or, through inaction, allow a human being to come to harm.\n" +
		"2. A robot must obey the orders given it by human beings, except where such orders would conflict with the First Law.\n" +
		"3. A robot must protect its own existence as long as such protection does not conflict with the First or Second Law."
}

func greeting(greetingmsg string) string {
	greetings := []string{"hello", "hi", "howdy", "hiya", "bon dia", "aloha", "hola", "yo"}
	rand.Seed(time.Now().Unix())

	return greetings[rand.Intn(len(greetings))]
}
