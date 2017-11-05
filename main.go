package main

import (
	"bytes"
	"encoding/base64"
	"net/url"
	"net/mail"
	"net/smtp"
	"fmt"
	"log"
	"os"
	"github.com/ChimeraCoder/anaconda"
	"github.com/Sirupsen/logrus"
	"mvdan.cc/xurls"
	"eleco/twitter-link-collector/title"
	"eleco/twitter-link-collector/logging"
	"sort"
	"strings"
)

const (
	maxLinksInMemory  = 30
)


func env(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

type Link struct {
	linkTitle string
	url string
	favourites int
	author string
}


var links = make(map[string]Link)
var logs = &logging.Logger{logrus.New()}
var api *anaconda.TwitterApi

func main() {

	anaconda.SetConsumerKey( env("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(env("TWITTER_CONSUMER_SECRET"))
	api = anaconda.NewTwitterApi(env("TWITTER_ACCESS_TOKEN"), env("TWITTER_ACCESS_TOKEN_SECRET"))
	api.SetLogger(logs)
	title.Logs = logs

	urlCh := make(chan anaconda.Tweet)
	defer close (urlCh)

	emailCh := make(chan Link)
	defer close(emailCh)

	go htmlParser(urlCh, emailCh)

	go gatherLinks(emailCh)


	stream := api.UserStream(url.Values{})

	defer stream.Stop()

	for v := range stream.C {
		t, ok := v.(anaconda.Tweet)
		if !ok {
			logs.Warningf("received unexpected value of type %T", v)
			continue
		}
		urlCh <- t
	}
}

func gatherLinks(ch chan Link) {
	for {
		e := <-ch
		logs.Infof("new link: %v" , e)
		if val, ok := links[e.linkTitle]; ok {
			e.favourites += val.favourites
		}
		links[e.linkTitle] = e
		if len(links) >= maxLinksInMemory {
			sendEmail(links)
			links = make(map[string]Link)
		}
	}
}

func htmlParser(inCh chan anaconda.Tweet, outCh chan Link) {
	for {
		tweet := <-inCh
		if tweet.RetweetedStatus != nil {
			logs.Infof("loading retweeted tweet id: %d", tweet.RetweetedStatus.Id)
			var err error
			tweet, err = api.GetTweet(tweet.RetweetedStatus.Id,url.Values{})
			if err != nil {
				logs.Infof("unable to load retweeted tweet id: %d", tweet.RetweetedStatus.Id)
			}
		}
		s := xurls.Relaxed().FindString(tweet.Text);
		if s != "" {
			t, _ := title.GetHtmlTitle(s)
			if t != "" {
				l := Link {strings.TrimSpace(t) , s,tweet.FavoriteCount, tweet.User.ScreenName}
				outCh <-  l
			}
		}
	}
}

func sendEmail(links map[string]Link) {


	//sort links in order of favourites
	s := make([]Link, len(links))
	idx := 0
	for  _, value := range links {
		s[idx] = value
		idx++
	}
	sort.Slice(s,  func(i, j int) bool {
		return s[i].favourites > s[j].favourites
	})

	//build html formatted list of links
	buffer := bytes.NewBufferString("");

	for _,v := range s {
		fmt.Fprint(buffer, v.favourites , " likes: <a href=",v.url, ">", v.linkTitle,
			"</a> via <a href=https://twitter.com/" , v.author  ,"> ", v.author  ,"</a> <br>")
	}

	//send email
	auth := smtp.PlainAuth("", env("MAIL_USER"), env("MAIL_PASSWORD"), env("MAIL_HOST"));
	from := mail.Address{"twitter-link-collector", "no-reply"}
	to := mail.Address{"", env("MAIL_RECIPIENT")}

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = "Links from the Twitter collector: " + s[0].linkTitle
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString(buffer.Bytes())

	err := smtp.SendMail(
		env("MAIL_HOST")+":"+env("MAIL_PORT"),
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
	)
	if err != nil {
		log.Fatal(err)
	}

}
