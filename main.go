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
)

const (
	maxLinksInMemory  = 10
)

var (
	consumerKey       = getenv("TWITTER_CONSUMER_KEY")
	consumerSecret    = getenv("TWITTER_CONSUMER_SECRET")
	accessToken       = getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = getenv("TWITTER_ACCESS_TOKEN_SECRET")
	mailUser          = getenv("MAIL_USER")
	mailPassword      = getenv("MAIL_PASSWORD")
	mailHost          = getenv("MAIL_HOST")
	mailPort          = getenv("MAIL_PORT")
	mailRecipient     = getenv("MAIL_RECIPIENT")

)

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

type Tuple [2]string

var links = make(map[string]string)
var logs = &logging.Logger{logrus.New()}
var api *anaconda.TwitterApi

func main() {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api = anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	api.SetLogger(logs)
	title.Logs = logs

	urlCh := make(chan anaconda.Tweet)
	emailCh := make(chan Tuple)

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
		logs.Infof(t.Text)
		urlCh <- t
	}
}

func gatherLinks(ch chan Tuple) {
	for {
		e := <-ch
		logs.Infof("new link title:%s url:%s", e[1], e[0])
		links[e[1]] = e[0]
		if len(links) > maxLinksInMemory {
			sendEmail(links)
			links = make(map[string]string)
		}
	}
}

func htmlParser(inCh chan anaconda.Tweet, outCh chan Tuple) {
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
				outCh <- Tuple{s, t}
			}
		}
	}
}

func sendEmail(links map[string]string) {

	buffer := bytes.NewBufferString("");
	i:=1;
	for k,v := range links {
		fmt.Fprint(buffer, i,". <a href=",v, ">", k, "</a><br>")
		i++;
	}

	auth := smtp.PlainAuth("", mailUser, mailPassword, mailHost);
	from := mail.Address{"twitter-link-collector", "no-reply"}
	to := mail.Address{"", mailRecipient}

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = "Links from the Twitter collector"
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString(buffer.Bytes())

	err := smtp.SendMail(
		mailHost+":"+mailPort,
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
	)
	if err != nil {
		log.Fatal(err)
	}

}
