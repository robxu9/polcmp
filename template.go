package main

import (
	"bytes"
	"errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/shurcooL/go/github_flavored_markdown"
	"honnef.co/go/js/console"
)

var (
	// ErrParse indicates that parsing
	ErrParse = errors.New("candidate: could not parse from markdown")
)

// Category refers to the grouping of issues.
type Category struct {
	Name   string
	Issues map[string]string
}

// Candidate represents the basic structure of the Markdown data file.
type Candidate struct {
	Name      string
	Positions []*Category
}

// CandidateFromMarkdown translates a markdown file that follows the template
// into a Candidate struct if possible. If not, it errors out.
func CandidateFromMarkdown(markdown string) (*Candidate, error) {
	result := github_flavored_markdown.Markdown([]byte(markdown))
	console.Log(string(result))

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(result))
	if err != nil {
		return nil, err
	}

	name := doc.Find("h1").Text()
	var positions []*Category

	doc.Find("table").Each(func(i1 int, s1 *goquery.Selection) {
		category := s1.Find("th[align*='right']").Text()
		console.Log("header: %s", category)

		header := true
		headerText := ""

		collection := make(map[string]string)

		s1.Find("td").Each(func(i2 int, s2 *goquery.Selection) {
			if header {
				headerText = s2.Text()
			} else {
				collection[headerText] = s2.Text()
				console.Log("%s: %s", headerText, s2.Text())
			}

			header = !header
		})

		positions = append(positions, &Category{
			Name:   category,
			Issues: collection,
		})
	})

	return &Candidate{
		Name:      name,
		Positions: positions,
	}, nil
}
