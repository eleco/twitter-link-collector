package main

import (
	"testing"
	"os"
	"log"
	"time"
)


func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup(){

	err := os.Setenv("TWITTER_CONSUMER_KEY","x")
	if (err!=nil) {
		log.Fatal(err)
	}
	os.Setenv("TWITTER_CONSUMER_SECRET","x")
	os.Setenv("TWITTER_ACCESS_TOKEN","x")
	os.Setenv("TWITTER_ACCESS_TOKEN_SECRET","x")
	os.Setenv("MAIL_USER","x")
	os.Setenv("MAIL_PASSWORD","x")
	os.Setenv("MAIL_HOST","x")
	os.Setenv("MAIL_PORT","x")
	os.Setenv("MAIL_RECIPIENT","x")
}


func Test(t *testing.T) {

	emailCh := make(chan Link)
	defer close(emailCh)

	go gatherLinks(emailCh)

	l1 := Link {"title","url",3, "screen_name"}
	l2 := Link {"title","url",3, "screen_name"}

	emailCh <- l1
	emailCh <- l2

	time.Sleep(time.Millisecond * 2)

	if links["title"].favourites!=6 {
		t.Fatalf("invalid favourites count, actual: %d, expected %d", links["title"].favourites, 6)
	}


}


