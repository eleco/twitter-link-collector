**twitter-link-collector**

A tool to extract url links from the user's twitter timeline, and email them to a predefined mail box (in batch of 10).


**Todo**
- More tests !
- extract url, title and image from the links
- Dockerise

**Pre-requisites**

- A twitter app as defined here: https://apps.twitter.com/
- Smtp email details: server, port and user/password.

**Installation - Option 1: run locally**

- Export the following environment parameters
       - TWITTER_CONSUMER_KEY, TWITTER_CONSUMER_SECRET, TWITTER_ACCESS_TOKEN, TWITTER_ACCESS_TOKEN_SECRET: from the Twitter app console
       - MAIL_USER, MAIL_PASSWORD, MAIL_HOST, MAIL_PORT, MAIL_RECIPIENT : the email client settings

- `go run Main.go`


**Installation - Option 2: use systemd on linux**

- Edit the twitter-link-collector.service file with the relevant account details

- sudo mv twitter-link-collector.service /etc/systemd/system

- sudo systemctl start twitter-link-collector

