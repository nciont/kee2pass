package entities

import (
	"encoding/xml"
	"errors"
	"io"
)

// XMLEntry a single keepass XML entry
type XMLEntry struct {
	UUID  string
	Notes string
	Pass  string
	Title string
	URL   string
	User  string
}

// UnmarshalXML implements Unmarshaler interface
func (entry *XMLEntry) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (e error) {
	var token xml.Token
	for {
		token, e = d.Token()

		switch e {
		case nil:
			break
		case io.EOF:
			return nil
		default:
			return e
		}

		if sToken, ok := token.(xml.StartElement); ok {
			entry.decodeTag(d, &sToken)
		}
	}

	return
}

func (entry *XMLEntry) decodeTag(d *xml.Decoder, tag *xml.StartElement) (e error) {
	var kv XMLKeyValue
	switch tag.Name.Local {
	case "String":
		if e = d.DecodeElement(&kv, tag); e == nil {
			entry.setKeyValue(&kv)
		}
		return
	case "UUID":
		token, e := d.Token()
		if e != nil {
			return e
		}
		if chardata, ok := token.(xml.CharData); ok {
			entry.UUID = string([]byte(chardata))
			break
		}
		return errors.New("Unexpected token type")
	}

	d.Skip() // move to next 
	return
}

func (entry *XMLEntry) setKeyValue(kv *XMLKeyValue) {
	switch kv.Key {
	case "Notes":
		entry.Notes = kv.Value
	case "Password":
		entry.Pass = kv.Value
	case "Title":
		entry.Title = kv.Value
	case "URL":
		entry.URL = kv.Value
	case "UserName":
		entry.User = kv.Value
	}
}
