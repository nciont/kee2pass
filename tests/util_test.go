package tests

import (
	kp "github.com/nciont/kee2pass"
	en "github.com/nciont/kee2pass/entities"
	"testing"
)

func TestGetItemName(t *testing.T) {
	for item, expected := range getItemNameTestData() {
		if actual := kp.GetItemName(item); actual != expected {
			t.Errorf("Expected <%s>, got <%s>", expected, actual)
		}
	}
}

func getItemNameTestData() map[*en.XMLEntry]string {
	return map[*en.XMLEntry]string{
		&en.XMLEntry{}: "", // empty strings
		&en.XMLEntry{ User: "user", URL: "zxcv", Title: "wrong" }: "user@zxcv", // user@host > title
		&en.XMLEntry{ Title: "  UpperCASE .// Slashes /// Spaces  " }: "uppercase.slashes.spaces", // normalization
		&en.XMLEntry{ User: "user" }: "user@???", // host unknown
	}
}

func TestPrettifyURL(t *testing.T) {
	for input, expected := range prettifyURLTestData() {
		if actual := kp.PrettifyURL(input); actual != expected {
			t.Errorf("Expected <%s>, got <%s>", expected, actual)
		}
	}
}

func prettifyURLTestData() map[string]string {
	return map[string]string{
		"": "", // empty strings
		"https://google.com/asdf/qwer": "google.com", // normal url
		"asdf": "asdf", // not url
		"http://www.example.com/2.aspx": "www.example.com", // url with subdomain
	}
}
