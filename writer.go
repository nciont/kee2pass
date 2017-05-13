package kee2pass

import (
	"fmt"
	"math/rand"
	"github.com/nciont/kee2pass/entities"
	"os/exec"
	"strconv"
	"text/template"
	"time"
)

// NewWriter Initializes a new writer
func NewWriter(dryrun bool, prefix string, tmpl *template.Template) *Writer {
	return &Writer{
		dryrun,
		prefix,
		tmpl,
		make(map[string]struct{}),
		rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Writer saves parsed items to password store
type Writer struct {
	dryrun    bool
	prefix    string
	tmpl      *template.Template
	itemNames map[string]struct{}
	rand      *rand.Rand
}

// Save stores an item
func (w *Writer) Save(item *entities.XMLEntry) (e error) {
	// Deal with dupe names
	itemName := w.prefix + GetItemName(item)
	if _, ok := w.itemNames[itemName]; ok {
		itemName += "." + strconv.FormatUint(w.rand.Uint64(), 36)
		fmt.Printf("Duplicate item name; renamed to %s\n", itemName)
	}
	w.itemNames[itemName] = struct{}{}

	if w.dryrun {
		fmt.Printf("DRYRUN: not writing item %s\n", itemName)
		return
	}

	cmd := exec.Command("pass", "insert", "--multiline", itemName)
	cmdIn, e := cmd.StdinPipe()

	if e != nil {
		return	
	}

	go func() {
		defer cmdIn.Close()
		e := w.tmpl.Execute(cmdIn, item)
		if e != nil {
			fmt.Println(e, itemName)
		}
	}()

	return cmd.Run()
}
