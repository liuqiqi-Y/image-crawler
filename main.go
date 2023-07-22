package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

const _letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandomString() string {
	str, err := randomString(15)
	if err != nil {
		return _letterBytes
	}
	return str
}

const _parallelistm = 2

func main() {

	makeDir()

	// Instantiate default collector
	c1 := colly.NewCollector(
	// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
	//colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"),
	//colly.Async(true),
	// colly.Debugger(&debug.LogDebugger{
	// 	Flag: log.Lshortfile,
	// }),
	)
	c2 := colly.NewCollector()
	// c2.Async = true
	// c2.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: _parallelistm})

	c3 := colly.NewCollector()
	// c3.Async = true
	// c3.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: _parallelistm})
	c3.SetRequestTimeout(0)

	extensions.RandomUserAgent(c1)
	extensions.RandomUserAgent(c2)
	extensions.RandomUserAgent(c3)

	c1filterStr1 := "body>main>div:first-of-type>section:first-of-type>ul>li>figure>a"
	c1.OnHTML(c1filterStr1, func(e *colly.HTMLElement) {
		imageLink := ""
		imageLink = e.Attr("href")
		c2.Visit(imageLink)
	})

	// c2.OnRequest(func(r *colly.Request) {
	// 	a := RandomString()
	// 	r.Headers.Set("User-Agent", a)
	// })

	c2.OnHTML("img[id=wallpaper][src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		c3.Visit(link)
	})

	// c3.OnRequest(func(r *colly.Request) {
	// 	r.Headers.Set("User-Agent", RandomString())
	// })

	c3.OnResponse(func(r *colly.Response) {
		ss := strings.Split(r.Request.URL.String(), "/")
		if len(ss) == 0 {
			return
		}
		caption := ss[len(ss)-1]

		f, err := os.Create(`.\images\` + caption)
		if err != nil {
			return
		}
		defer f.Close()
		io.Copy(f, bytes.NewReader(r.Body))
	})
	c1.OnError(func(r *colly.Response, err error) {
		log.Printf("c1: request to %s error: %s\n", r.Request.URL, err.Error())
	})
	c2.OnError(func(r *colly.Response, err error) {
		log.Printf("c2: request to %s error: %s\n", r.Request.URL, err.Error())
	})
	c3.OnError(func(r *colly.Response, err error) {
		log.Printf("c3: request to %s error: %s\n", r.Request.URL, err.Error())
	})
	// c1.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})
	// c1.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})
	// c2.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})
	// c3.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	for i := 0; i < 1; i++ {
		url := fmt.Sprintf("https://wallhaven.cc/toplist?page=%d", i+1)
		c1.Visit(url)
		fmt.Printf("第%d页访问结束\n", i+1)
	}

}

func makeDir() {
	dirName := "images"

	// 检查目录是否存在
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		// 目录不存在，创建它
		err := os.Mkdir(dirName, 0755) // 0755表示目录权限，这里是读、写、执行权限给所有用户
		if err != nil {
			panic(err.Error())
		}
	}
}

func randomString(length int) (string, error) {
	// 计算生成随机字符串所需的字节数
	numBytes := (length * 6) / 8 // base64编码后，每6个比特位对应一个字符

	// 生成随机字节序列
	randomBytes := make([]byte, numBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// 使用base64编码将字节序列转换为字符串
	randomString := base64.RawURLEncoding.EncodeToString(randomBytes)

	// 截取指定长度的字符串
	if len(randomString) > length {
		randomString = randomString[:length]
	}

	return randomString, nil
}
