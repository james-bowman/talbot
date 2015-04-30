package brain

import (
	"fmt"
	"github.com/james-bowman/slack"
	"log"
	"regexp"
)

var tellCommand = regexp.MustCompile("(?i)tell (<.*?>)( .*)? <#(.*?)>? (.*)")

type Action struct {
	// indicates whether the action should be checked for all messages (true) or just those
	// directed at the bot (false).  Defaults to false.
	Hear bool

	// pattern to match in the message for the action's Answerer function to execute
	Regex *regexp.Regexp

	// usage example
	Usage string

	// textual help description for the action
	Description string

	// function to execute if the Regex matches
	Answerer actionFunc
}

func (a Action) String() string {
	return fmt.Sprintf("`%s` - %s", a.Usage, a.Description)
}

type actionFunc func(string) string

var defaultAction actionFunc

var actions actionList

type actionList []Action

func (a actionList) handle(message string, asked bool) string {
	var response string

	var action Action
	for _, action = range a {
		if action.Hear || asked {
			matches := action.Regex.FindStringSubmatch(message)

			if len(matches) > 0 {
				response = action.Answerer(message)
				break
			}
		}
	}

	return response
}

func Register(newAction Action) {
	log.Printf("\nRegistering Action: %s", newAction)
	actions = append(actions, newAction)
}

func RegisterDefault(action actionFunc) {
	log.Printf("\nRegistering Default Action: %#v", action)
	if defaultAction != nil {
		panic(fmt.Sprintf("Attempted to set default action failed because one has already been registered: %#v", defaultAction))
	}
	defaultAction = action
}

func init() {
	Register(Action{
		Regex:       regexp.MustCompile("(?i)help"),
		Usage:       "help",
		Description: "Reply with usage instructions",
		Answerer: func(dummy string) string {
			response := "I am a robot that listens to the team's chat and provides automated functions." +
				"  I currently support the following commands:\n"

			for _, value := range actions {
				if value.Usage != "" {
					response = fmt.Sprintf("%s\n\t%s", response, value)
				}
			}

			return response + "\n\nI can do some other things too - try asking me something!"
		},
	})
}

func OnHeardMessage(message *slack.Message) {
	response := actions.handle(message.Text, false)
	if response != "" {
		if err := message.Send(response); err != nil {
			// gulp!
			log.Printf("Error sending message: %s with message: '%s'", err, response)
		}
	}
}

func OnAskedMessage(message *slack.Message) {
	var response string

	log.Printf("%s-> %s", message.From, message.Text)

	matches := tellCommand.FindStringSubmatch(message.Text)
	if len(matches) > 0 {
		message.Tell(matches[3], matches[1]+": "+matches[4])
		return
	}

	response = actions.handle(message.Text, true)

	if response == "" {
		if defaultAction != nil {
			// if sentence not matched to a supported request then fallback to a default
			response = defaultAction(message.Text)
		} else {
			response = "I don't understand.\n_Ask me_ `help` _for the list of commands I currently understand_"
		}
	}

	log.Printf("me-> %s", response)

	if err := message.Respond(response); err != nil {
		// gulp!
		log.Printf("Error responding to message: %s\nwith Message: '%s'", err, response)
	}
}
