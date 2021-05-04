package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
	// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
	//colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"),
	//colly.Async(true),
	)
	c1 := c.Clone()
	c2 := c.Clone()
	c.OnHTML("ul figure a", func(e *colly.HTMLElement) {
		var imageLink string
		// e.ForEach("a", func(i int, h *colly.HTMLElement) {
		// 	imageLink = h.Attr("href")
		// 	return
		// })
		imageLink = e.Attr("href")
		fmt.Printf("image link found: %s\n", imageLink)
		c1.Visit(imageLink)
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c2.OnResponse(func(r *colly.Response) {
		ss := strings.Split(r.Request.URL.String(), "/")
		if len(ss) == 0 {
			return
		}
		caption := ss[len(ss)-1]
		fmt.Printf("Downloading image: %s\n", caption)
		f, err := os.Create("./image/" + caption)
		if err != nil {
			panic(err)
		}
		io.Copy(f, bytes.NewReader(r.Body))
	})
	c1.OnHTML("img[id=wallpaper]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		c2.Visit(link)
	})
	c.Visit("https://wallhaven.cc/toplist?page=7")
}
