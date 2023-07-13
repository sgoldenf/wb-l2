package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var (
	attrMap = map[string]struct{}{
		"cite": {}, "data": {}, "href": {}, "src": {}, "srcset": {},
	}
	tagMap = map[string]struct{}{
		"blockquote": {}, "del": {}, "ins": {}, "q": {}, "object": {}, "a": {},
		"area": {}, "base": {}, "link": {}, "audio": {}, "embed": {},
		"iframe": {}, "img": {}, "input": {}, "script": {}, "source": {},
		"track": {}, "video": {},
	}
)

func main() {
	r := flag.Bool("r", false, "recursive")
	flag.Parse()
	paths := flag.Args()
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	c := crawler{visited: map[string]struct{}{}, pwd: pwd, recursive: *r}
	for _, path := range paths {
		c.crawl(path)
	}
}

type crawler struct {
	http.Client
	visited   map[string]struct{}
	pwd       string
	recursive bool
}

func (c *crawler) crawl(link string) {
	if _, ok := c.visited[link]; ok {
		return
	}
	c.visited[link] = struct{}{}

	fmt.Println(link)

	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
		link = "http://" + link
	}

	resp, err := c.Get(link)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	if c.recursive {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		c.crawlRecursive(bodyBytes, link)
	} else {
		save(resp.Body, "index.html")
	}
}

func (c *crawler) crawlRecursive(bodyBytes []byte, link string) {
	links := parseLinks(bytes.NewReader(bodyBytes))

	path := strings.TrimPrefix(link, "http://")
	path = strings.TrimPrefix(path, "https://")

	if len(links) != 0 {
		if err := os.MkdirAll(path, 0777); err != nil {
			log.Println(err)
			return
		}
	}

	save(bytes.NewReader(bodyBytes), path)

	for _, l := range links {
		if !strings.Contains(link, l) {
			if l[0] == '/' {
				c.crawl(link + l)
			} else {
				c.crawl(l)
			}
		}
	}
}

func parseLinks(body io.Reader) (links []string) {
	tokenizer := html.NewTokenizer(body)
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return
		case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if _, ok := tagMap[token.Data]; !ok {
				continue
			}
			link := parseToken(token)
			if link != "" {
				link = strings.TrimPrefix(link, "//")
				links = append(links, link)
			}
		}
	}
}

func parseToken(t html.Token) string {
	for _, attr := range t.Attr {
		if _, ok := attrMap[attr.Key]; ok {
			return attr.Val
		}
	}
	return ""
}

func save(body io.Reader, path string) {
	fileInfo, err := os.Stat(path)
	if err == nil && fileInfo.IsDir() {
		path = filepath.Join(path, "index.html")
	}

	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, body)
	if err != nil {
		log.Println(err)
	}
}
