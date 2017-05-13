package kee2pass

import (
	"encoding/xml"
	"errors"
	"io"
	"github.com/nciont/kee2pass/entities"
)

// TagEntry entry tag name
const TagEntry = "Entry"

// TagGroup entry grouping tag
const TagGroup = "Group"

// NewConverter ctor
func NewConverter(settings *Settings) *Converter {
	var converter = &Converter{
		data: settings.Data,
	}

	return converter.init()
}

// Converter Converts KeePass XML to pass format
type Converter struct {
	data    io.Reader
	decoder *xml.Decoder
}

func (kp *Converter) init() *Converter {
	kp.decoder = xml.NewDecoder(kp.data)
	return kp
}

// Next read next entry
func (kp *Converter) Next() (*entities.XMLEntry, error) {
	var stok, e = kp.moveToNextEntry()

	if e != nil {
		return nil, e
	}

	item := &entities.XMLEntry{}
	if e = kp.decoder.DecodeElement(item, stok); e == nil {
		return item, nil
	}

	return nil, e
}

// moveToNextEntry moves to next entry tag.
func (kp *Converter) moveToNextEntry() (*xml.StartElement, error) {
	for {
		token, e := kp.decoder.Token()

		if e != nil {
			return nil, e
		}

		if stok, ok := token.(xml.StartElement); ok && stok.Name.Local == TagEntry {
			return &stok, nil
		}
	}

	return nil, errors.New("Unreachable line")
}
