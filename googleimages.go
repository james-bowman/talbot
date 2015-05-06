package main // import "github.com/james-bowman/talbot"

import (
	"encoding/json"
	"fmt"
	"github.com/james-bowman/talbot/brain"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

func init() {
	brain.Register(brain.Action{
		Regex:       regexp.MustCompile("(?i)(image|img)( me)? (.*)"),
		Usage:       "image me <search expression>",
		Description: "Queries Google Images for _search expression_ and returns random result",
		Answerer: func(message string) string {
			reExpression, err := regexp.Compile(`(?i)(image|img)( me)? (.*)`)

			if err != nil {
				log.Printf("Error compiling Regex to obtain query expression: %s", err)
				return ""
			}

			searchExpression := reExpression.FindStringSubmatch(message)

			if len(searchExpression) > 0 {
				return imageSearch(searchExpression[3], false, false)
			}
			return ""
		},
	})

	brain.Register(brain.Action{
		Regex:       regexp.MustCompile("(?i)(animate)( me)? (.*)"),
		Usage:       "animate me <search expression>",
		Description: "The sames as `image me` except requests an animated gif matching _search expression_",
		Answerer: func(message string) string {
			reExpression, err := regexp.Compile(`(?i)(animate)( me)? (.*)`)

			if err != nil {
				log.Printf("Error compiling Regex to obtain query expression: %s", err)
				return ""
			}

			searchExpression := reExpression.FindStringSubmatch(message)

			if len(searchExpression) > 0 {
				return imageSearch(searchExpression[3], true, false)
			}
			return ""
		},
	})

	brain.Register(brain.Action{
		Regex:       regexp.MustCompile("(?i)(?:mo?u)?s?ta(?:s|c)h(?:e|ify)?(?: me)? (.*)"),
		Usage:       "mustache me <search expression or URL>",
		Description: "Queries Google Images for _search expression_ and adds a mustache or simply mustachifies the image at _url_",
		Answerer: func(message string) string {
			mustachify := "http://mustachify.me/rand?src=%s"

			reExpression, err := regexp.Compile(`(?i)(?:mo?u)?s?ta(?:s|c)h(?:e|ify)?(?: me)? (.*)`)
			if err != nil {
				log.Printf("Error compiling Regex to obtain query expression: %s", err)
				return ""
			}

			searchExpression := reExpression.FindStringSubmatch(message)
			expr := searchExpression[1]

			reURL, err := regexp.Compile(`(?i)^<(https?:\/\/.*)>`)
			if err != nil {
				log.Printf("Error compiling Regex to obtain URL: %s", err)
				return ""
			}

			if reURL.MatchString(expr) {
				return fmt.Sprintf(mustachify, url.QueryEscape(reURL.FindStringSubmatch(expr)[1]))
			}

			return fmt.Sprintf(mustachify, url.QueryEscape(imageSearch(expr, false, true)))
		},
	})
}

func imageSearch(expr string, animated bool, faces bool) string {
	googleURL, err := url.Parse("http://ajax.googleapis.com/ajax/services/search/images")
	if err != nil {
		log.Printf("Error parsing Google Images URL: %s", err)
		return ""
	}

	q := googleURL.Query()
	q.Set("v", "1.0")
	q.Set("rsz", "8")
	q.Set("safe", "active")

	if animated {
		q.Set("imgtype", "animated")
	}

	if faces {
		q.Set("imgtype", "face")
	}

	q.Set("q", expr)

	googleURL.RawQuery = q.Encode()

	resp, err := http.Get(googleURL.String())
	if err != nil {
		log.Printf("Error calling url '%s' : %s ", googleURL, err)
		return "Sorry I had a problem finding that image from Google"
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Error reading results from HTTP Request '%s': %s", googleURL, err)
		return "Sorry I had a problem finding that image from Google"
	}

	fmt.Printf("Response from Google: %s", body)

	var results map[string]interface{}
	err = json.Unmarshal(body, &results)
	if err != nil {
		log.Printf("%T\n%s\n%#v\n", err, err, err)
		switch v := err.(type) {
		case *json.SyntaxError:
			log.Println(string(body[v.Offset-40 : v.Offset]))
		}
		log.Printf("%s", body)
		return "Sorry I had a problem finding that image from Google"
	}

	imageList, ok := results["responseData"].(map[string]interface{})["results"]

	var image string
	var images []interface{}
	if ok {
		images = imageList.([]interface{})

		if len(images) > 0 {
			rand.Seed(time.Now().Unix())

			i := images[rand.Intn(len(images))].(map[string]interface{})
			image = i["unescapedUrl"].(string)
		} else {
			image = images[0].(map[string]interface{})["unescapedUrl"].(string)
		}

		reURL, err := regexp.Compile(`(?i).(png|jpe?g|gif)$`)
		if err != nil {
			log.Printf("Error compiling Regex to parse returned URL '%s': %s", image, err)
			return ""
		}

		if !reURL.MatchString(image) {
			image = image + ".png"
		}
	}

	return image
}
