package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

const _letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = _letterBytes[rand.Intn(len(_letterBytes))]
	}
	return string(b)
}

func main() {
	// Instantiate default collector
	c1 := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		//colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"),
		//colly.Async(true),
		colly.Debugger(&debug.LogDebugger{
			Flag: log.Lshortfile,
		}))

	c1.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	c2 := c1.Clone()
	c3 := c1.Clone()

	// extensions.RandomUserAgent(c1)
	// extensions.RandomUserAgent(c2)
	// extensions.RandomUserAgent(c3)

	c1filterStr1 := "body>main>div:first-of-type>section:first-of-type>ul>li>figure>a"
	c1.OnHTML(c1filterStr1, func(e *colly.HTMLElement) {
		imageLink := ""
		imageLink = e.Attr("href")
		// log.Printf("image detail link found: %s\n", imageLink)
		c2.Visit(imageLink)
	})

	c2.OnHTML("img[id=wallpaper][src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		// log.Printf("image data link found: %s\n", link)
		c3.Visit(link)
	})
	c3.OnResponse(func(r *colly.Response) {
		ss := strings.Split(r.Request.URL.String(), "/")
		if len(ss) == 0 {
			// log.Printf("invalid image data link : %s\n", r.Request.URL.String())
			return
		}
		caption := ss[len(ss)-1]
		// log.Printf("Downloading image: %s\n", caption)
		f, err := os.Create(`.\image\` + caption)
		if err != nil {
			// log.Printf("open file error: %s\n", err.Error())
			return
		}
		defer f.Close()
		io.Copy(f, bytes.NewReader(r.Body))
	})
	c3.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	//c1.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})
	// c1.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})
	// c2.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})
	// c3.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	// c1.Wait()
	// c2.Wait()
	// c3.Wait()

	for i := 0; i < 1; i++ {
		url := fmt.Sprintf("https://wallhaven.cc/toplist?page=%d", i+1)
		c1.Visit(url)
		fmt.Printf("第%d页访问结束", i+1)
	}

}
