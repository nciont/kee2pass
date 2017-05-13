package tests

import (
	"github.com/nciont/kee2pass"
	"os"
	"strings"
	"testing"
)

func TestConverter(t *testing.T) {
	datafile, e := os.Open("testdata/data.xml")

	if e != nil {
		t.Error("error opening file")
	}

	kp := kee2pass.NewConverter(&kee2pass.Settings{Data: datafile})
	entry, e := kp.Next()

	if e != nil {
		t.Fatal("XML parse error")
	}

	if entry == nil {
		t.Fatal("Failed to make an entry")
	}

	// verify data in first entry

	if entry.UUID != "1WgLuXO5YEeOn2OGZT8HwQ==" {
		t.Error("Invalid UUID")
	}

	if !strings.HasPrefix(entry.Notes, "testnote testnote") {
		t.Error("Invalid Notes")
	}

	if entry.Pass != "Secret<Pass>" {
		t.Error("Invalid Pass")
	}

	if entry.Title != "Entry Title" {
		t.Error("Invalid Title")
	}

	if entry.URL != "https://bing.com/" {
		t.Error("Invalid URL")
	}

	if entry.User != "testuser" {
		t.Error("Invalid User")
	}
}
