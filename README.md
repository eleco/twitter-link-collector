**twitter-link-collector**

A tool to extract url links from the user's twitter timeline, and email them back to the user's mailbox


**Todo**
- More tests !
- extract url, title and image from the links
- Dockerise


**Installation**

- Create a twitter app here: https://apps.twitter.com/

- Export the following environment parameters
       - TWITTER_CONSUMER_KEY, TWITTER_CONSUMER_SECRET, TWITTER_ACCESS_TOKEN, TWITTER_ACCESS_TOKEN_SECRET: from the Twitter app console
       - MAIL_USER, MAIL_PASSWORD, MAIL_HOST, MAIL_PORT, MAIL_RECIPIENT : the email client settings

- `go run Main.go`

- the app will email links collected by batches of 10.