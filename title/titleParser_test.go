package title

import (
	"testing"
)

func TestHtmlToRst(t *testing.T) {

	title, ok := GetHtmlTitle("http://www.bbc.co.uk/news");
	if (title !="Home - BBC News") {
		t.Fatalf("failed to parse title ok: %t  title: %s", ok, title)
	}
}

//https://t.co/l50X3dtArh

