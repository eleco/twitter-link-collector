package title

import (
	"testing"
	"golang.org/x/net/html"
)

func TestHtmlToRst(t *testing.T) {

	title, ok := GetHtmlTitle("http://www.bbc.co.uk/news");
	if (title !="Home - BBC News") {
		t.Fatalf("failed to parse title ok: %t  title: %s", ok, title)
	}
}

func TestShouldParseEmptyTitleNode_withoutCrashing(t *testing.T) {

	node := html.Node{
		Type:     html.ElementNode,
		Data:     "title",
	}
	traverse(&node)
}



