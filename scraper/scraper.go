package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gocolly/colly"
)

type Entry struct {
	Title    string `json:"title"`
	Question string `json:"question"`
}

func Scrape(uri string) {
	var titles []string
	var questions []string

	entries := []Entry{}

	fmt.Println("Scraping from:", uri)

	c := colly.NewCollector()

	c.OnHTML("h3", func(e *colly.HTMLElement) {
		titles = append(titles, e.Text)
	})

	c.OnHTML("div.box2", func(e *colly.HTMLElement) {
		questions = append(questions, e.Text)
	})

	c.Visit(uri)

	fmt.Println(len(questions))

	if len(titles) == len(questions) {
		for i := range titles {
			entry := Entry{Title: titles[i], Question: strings.TrimSpace(questions[i])}
			entries = append(entries, entry)
		}
	}

	j, _ := json.MarshalIndent(entries, "", "	")
	fmt.Println(string(j))

	_ = ioutil.WriteFile("entries.json", j, 0644)
}
