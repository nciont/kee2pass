package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"netelevate.com/kee2pass"
	"os"
	"os/exec"
	tmpl "text/template"
)

var opts = struct {
	Input  *string
	Prefix *string
	Dryrun *bool
}{
	flag.String("input", "", "KeePass XML file to convert"),
	flag.String("prefix", "", "Prepended to item names"),
	flag.Bool("dry-run", false, "Read-only mode"),
}

// template to use
// todo keep in an external file, use flags to customize
var templateString = `
{{- .Pass }}

{{if .Title -}}
{{.Title}}
================================
{{- end}}
User: {{.User}}
Pass: {{.Pass}}
 URL: {{.URL}}

UUID: {{.UUID}}

{{- if .Notes }}

-----------[ NOTES ]------------
{{.Notes}}
{{end}}
`

var template *tmpl.Template

func init() {
	var e error
	flag.Parse()

	if !validOpts() {
		usage()
		os.Exit(0)
	}

	if e = checkOpts(); e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	// init template
	template, e = tmpl.New("main").Parse(templateString)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

func main() {
	input, e := os.Open(*opts.Input)

	if e != nil {
		fmt.Println("Could not open input file; permissions?")
		os.Exit(2)
	}

	converter := kee2pass.NewConverter(&kee2pass.Settings{Data: input})
	writer := kee2pass.NewWriter(*opts.Dryrun, *opts.Prefix, template)

	for {
		if item, e := converter.Next(); e == nil {
			if e := writer.Save(item); e != nil {
				fmt.Println(e)
			}
		} else {
			if e == io.EOF {
				fmt.Println("Done")
				os.Exit(0)
			}

			fmt.Println(e)
			os.Exit(1)
		}
	}
}

func validOpts() bool {
	if *opts.Input == "" {
		return false
	}

	return true
}

func checkOpts() error {
	// check input
	if info, e := os.Stat(*opts.Input); e == nil {
		if !info.Mode().IsRegular() {
			return errors.New("Input must point at file")
		}
	} else {
		return e
	}

	// check pass binary
	if _, e := exec.LookPath("pass"); e != nil {
		return e
	}

	return nil
}

func usage() {
	fmt.Println("Usage: kee2pass [opts] -input=path")
	flag.PrintDefaults()
}
