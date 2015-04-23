# Talbot

Talbot is a basic chat bot (robot) written in Go that listens and responds to [Slack](https://slack.com) messages.

To make use of Talbot, you will need to setup a user account or a ['bot' user](https://api.slack.com/bot-users) (a special kind of user account specifically designed for programatic access to APIs) and obtain an [authentication token](https://api.slack.com/web#basics) from Slack.

Once you have a user account and an authentication token you can connect Talbot to slack.  Simply paste the authentication token into the slack.token file in the same directory as the executable.

Talbot can be given any name within [Slack](https://slack.com) - he will use and respond to whatever nick/username is registered for the bot account within Slack.  For the purposes of the examples below, we have assumed talbot has been used as the nick in Slack.

Talbot can be extended by registering additional actions as follows.

``` go
package main

import "github.com/james-bowman/talbot/brain"

func init() {
	brain.Register(brain.Action{
    		Regex: regexp.MustCompile("(?i)open the pod bay doors"),
	    	Usage: "open the pod bay doors",
			Description: "Opens the bay doors", 
    		Answerer: func(dummy string) string {
    			return "I'm sorry Dave, I can't do that right now."
			},
	})
}
```

This registers an action that will be executed when any user directs the message `open the pod bay doors` at the bot i.e.

    @talbot: open the pod bay doors

in an open channel or simply

    open the pod bay doors

as a Direct Message to Talbot.

Out of the box, Talbot supports the following commands:

- `help` - Reply with instructions
- `ping` - Reply with pong
- `hello` - Reply with a random greeting
- `the rules` - Make sure talbot knows the rules

