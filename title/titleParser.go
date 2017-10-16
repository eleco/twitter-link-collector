package title

import (
	"golang.org/x/net/html"
	"net/url"
	"net/http"
	"eleco/twitter-link-collector/logging"
)

var Logs *logging.Logger

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}
	return "", false
}

func GetHtmlTitle(urlstr string) (string, bool) {

	parsed, err := url.Parse(urlstr)
	if err != nil {
		Logs.Warnf("unable to parse url: %s", urlstr, err)
		return "", false
	}
	if parsed.Scheme=="" {
		parsed.Scheme = "http"
	}

	resp, err := http.Get(urlstr)
	if err != nil {
		Logs.Warnf("unable to get: %s", urlstr, err)
		return "", false
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		Logs.Warnf("unable to parse html at: %s", urlstr, err)
		return "", false
	}

	return traverse(doc)
}
