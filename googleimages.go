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
	imageRegex := "(?i)(image|img)( me)? (.*)"
	animateRegex := "(?i)(animate)( me)? (.*)"
	mustacheRegex := "(?i)(?:mo?u)?s?ta(?:s|c)h(?:e|ify)?(?: me)? (.*)"

	brain.Register(brain.Action{
		Regex:       regexp.MustCompile(imageRegex),
		Usage:       "image me <search expression>",
		Description: "Queries Google Images for _search expression_ and returns random result",
		Answerer: func(message string) string {
			return processQueryAndSearch(message, imageRegex, false)
		},
	})

	brain.Register(brain.Action{
		Regex:       regexp.MustCompile(animateRegex),
		Usage:       "animate me <search expression>",
		Description: "The sames as `image me` except requests an animated gif matching _search expression_",
		Answerer: func(message string) string {
			return processQueryAndSearch(message, animateRegex, true)
		},
	})

	brain.Register(brain.Action{
		Regex:       regexp.MustCompile(mustacheRegex),
		Usage:       "mustache me <search expression or URL>",
		Description: "Queries Google Images for _search expression_ and adds a mustache or simply mustachifies the image at _url_",
		Answerer: func(message string) string {
			mustachify := "http://mustachify.me/rand?src=%s"

			reExpression, err := regexp.Compile(mustacheRegex)
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

func processQueryAndSearch(message string, regex string, animated bool) string {
	reExpression, err := regexp.Compile(regex)

	if err != nil {
		log.Printf("Error compiling Regex to obtain query expression: %s", err)
		return ""
	}

	searchExpression := reExpression.FindStringSubmatch(message)

	if len(searchExpression) > 0 {
		return imageSearch(searchExpression[3], animated, false)
	}
	return ""
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
	q.Set("q", expr)

	if animated {
		q.Set("imgtype", "animated")
	}

	if faces {
		q.Set("imgtype", "face")
	}

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

	var results map[string]interface{}
	if err = json.Unmarshal(body, &results); err != nil {
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
