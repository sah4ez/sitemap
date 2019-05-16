package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var (
	Nil = struct{}{}
)

type Parser struct {
	Base   string
	Url    string
	Urls   map[string]struct{}
	Client *http.Client
}

func (p *Parser) Parse(n *html.Node) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "base":
			for _, a := range n.Attr {
				if a.Key == "href" {
					p.Base = a.Val
					fmt.Println("base", a.Val)
					break
				}
			}
		case "a":
			for _, a := range n.Attr {
				if a.Key == "href" {
					if strings.HasPrefix(a.Val, "mailto:") {
						break
					}
					if strings.HasPrefix(a.Val, "/") {
						u, _ := url.Parse(a.Val)
						if u.Scheme == "" {
							u := p.Url + a.Val
							if _, ok := p.Urls[u]; !ok {
								p.Urls[u] = Nil
							}
						} else {
							u := p.Base + a.Val
							if _, ok := p.Urls[u]; !ok {
								p.Urls[u] = Nil
							}
						}
					} else {
						u, _ := url.Parse(a.Val)
						if u.Scheme == "" {
							u := p.Url + a.Val
							if _, ok := p.Urls[u]; !ok {
								p.Urls[u] = Nil
							}
						} else {
							u := a.Val
							if _, ok := p.Urls[u]; !ok {
								p.Urls[u] = Nil
							}
						}

					}
					break
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.Parse(c)
	}
}

func (p *Parser) Get() ([]byte, error) {
	resp, err := p.Client.Get(p.Url)
	if err != nil {
		return nil, fmt.Errorf("error http get: %s", err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error read body: %s", err)
	}
	resp.Body.Close()
	return b, nil
}

func Pool(size int) chan *Parser {
	pool := make(chan *Parser, size)

	for i := 0; i < size; i++ {
		pool <- New()
	}
	return pool
}

func P(pool chan *Parser, u string) map[string]struct{} {
	p := <-pool
	defer func() { pool <- p }()
	p.Url = u
	b, err := p.Get()
	if err != nil {
		fmt.Println("error get:", err)
	}
	if b == nil {
		return nil
	}

	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		fmt.Println("error read body:", err)
		return nil
	}
	p.Parse(doc)
	return p.Urls
}

func New() *Parser {
	return &Parser{
		Urls: map[string]struct{}{},
		Client: &http.Client{
			Timeout: time.Second * 15,
		},
	}
}
