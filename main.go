package main

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"honnef.co/go/js/console"
)

const (
	url  = "https://github.com/polcmp/2016_usa_presidential.git"
	cors = "https://cors-anywhere.herokuapp.com/"
)

var (
	// ErrNoTemplate is thrown if no TEMPLATE.md can be found in the git repository
	ErrNoTemplate = errors.New("no TEMPLATE.md found")
	// ErrBadData is thrown when there is a bad dataset
	ErrBadData = errors.New("bad dataset")
	jQ         = jquery.NewJQuery
	git        = js.Global.Get("GitApi")
)

func updateProgress(info string) {
	jQ("#progress").SetText(info)
}

func errorFunc(err interface{}) {
	console.Log(err)
	jQ("#main").SetHtml(`<p class="text-danger text-center">We couldn't retrieve the data for you! (　ﾟДﾟ)＜!! Refresh to try again?</p>`)
	panic(err)
}

func readEntries(entries chan *js.Object) {
	var candidates []*Candidate
	var template *Candidate

	for v := range entries {
		if v.Get("isFile").Bool() { // we only want files
			name := v.Get("name").String()
			if name == "LICENSE" || name == "README.md" {
				// skip
				continue
			}
			updateProgress(fmt.Sprintf("parsing %v...", name))

			chRet := make(chan *js.Object, 1)

			v.Call("file", func(file *js.Object) {
				go func() {
					chRet <- file
					close(chRet)
				}()
			}, func(err *js.Object) {
				close(chRet)
				go errorFunc(err)
			})

			file := <-chRet
			if file == nil {
				// the error function was called, jump ship!
				return
			}

			reader := js.Global.Get("FileReader").New()

			waitReader := make(chan struct{}, 1)
			reader.Set("onloadend", func() {
				go func() {
					waitReader <- struct{}{}
					close(waitReader)
				}()
			})

			reader.Call("readAsText", file)
			<-waitReader

			str := reader.Get("result").String()
			console.Log(name)
			console.Log(str)

			result, err := CandidateFromMarkdown(str)
			if err != nil {
				errorFunc(err)
			}

			if name == "TEMPLATE.md" {
				template = result
			} else {
				candidates = append(candidates, result)
			}
		}
	}

	if template == nil {
		errorFunc(ErrNoTemplate)
	}

	// translate candidates into the table format
	var tables []*Table
	for k, v := range template.Positions {
		// build list of candidates
		var c []string
		p := make(map[string][]string)

		for _, candidate := range candidates {
			c = append(c, candidate.Name)
		}

		for i := range v.Issues {
			var stances []string
			for _, candidate := range candidates {
				stances = append(stances, candidate.Positions[k].Issues[i])
			}
			p[i] = stances
		}

		tables = append(tables, &Table{
			Name:       v.Name,
			Candidates: c,
			Positions:  p,
		})
	}

	console.Log(tables)

	var buf bytes.Buffer
	err := tableTmpl.Execute(&buf, tables)
	if err != nil {
		errorFunc(ErrBadData)
	}

	jQ("#main").SetHtml(buf.String())

	console.Log("done")
}

func dirReadAll(dirReader *js.Object, consumer chan *js.Object) {
	dirReader.Call("readEntries", func(e []*js.Object) {
		go func() {
			if len(e) != 0 {
				for _, v := range e {
					consumer <- v
				}
				go dirReadAll(dirReader, consumer)
			} else {
				close(consumer)
			}
		}()
	}, errorFunc)
}

func afterClone(dir *js.Object) {
	updateProgress("parsing through the list of candidates...")
	dirReader := dir.Call("createReader")
	entries := make(chan *js.Object, 10)
	go dirReadAll(dirReader, entries)
	go readEntries(entries)
}

func cloneProgress(progress *js.Object) {
	percentage := progress.Get("pct").Int()
	msg := progress.Get("msg").String()

	updateProgress(fmt.Sprintf("[%v] %v", percentage, msg))
}

func afterDirectory(directory *js.Object) {
	jQ("#progress").SetText(`retrieving the list of candidates...`)
	// All git functions go here, including parsing and stuff.
	cloneOps := map[string]interface{}{
		"dir":   directory,
		"url":   cors + url,
		"depth": 1,
		"progress": func(prog *js.Object) {
			go cloneProgress(prog)
		},
	}

	git.Call("clone", cloneOps, func() {
		go afterClone(directory)
	}, errorFunc)
}

func afterDelete(fs, dir *js.Object) {
	dir.Call("removeRecursively", func() {
		go afterFilesystem(fs, false)
	}, errorFunc)
}

func afterFilesystem(fs *js.Object, del bool) {
	fs.Get("root").Call("getDirectory", "/checkoutloc", map[string]interface{}{
		"create": true,
	}, func(dir *js.Object) {
		if del {
			go afterDelete(fs, dir)
		} else {
			go afterDirectory(dir)
		}
	})
}

func main() {
	updateProgress("starting...")
	// request a filesystem
	js.Global.Call("requestFileSystem", js.Global.Get("TEMPORARY"), float64(5*1024*1024*1024), func(fs *js.Object) {
		go afterFilesystem(fs, true)
	})
}
