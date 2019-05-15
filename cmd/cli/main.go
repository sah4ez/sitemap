package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/html"

	"github.com/sah4ez/sitemap/pkg/parser"
	"github.com/sah4ez/sitemap/pkg/sitemap"
)

var (
	parallel   = flag.Int("parallel", 1, "number of parallel workers to navigate through site")
	outputFile = flag.String("output-file", "./sitemap.xml", "output file path")
	maxDepth   = flag.Int("max-depth", 2, "max depth of url navigation recursion")
)

func init() {
	flag.Parse()
	if *parallel < 1 {
		*parallel = 1
	}
	if *maxDepth < 1 {
		*maxDepth = 1
	}
}

type walk func(int, string, *sitemap.Urlset, chan walk, *sync.WaitGroup)

func (w *walk) deep(level int, urlStr string, us *sitemap.Urlset, pool chan *walk, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
		if pool != nil {
			pool <- w
		}
	}()

	fmt.Println("level", level, "depth", *maxDepth)
	if level >= *maxDepth {
		return
	}
	p := P(urlStr)
	if p == nil {
		return
	}

	l := level + 1
	for k := range p.Urls {
		w.deep(l, k, us, nil, nil)
		us.Add(k)
	}
}

func main() {
	urlStr := os.Args[len(os.Args)-1]
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	p := P(urlStr)

	us := sitemap.New()

	pool := make(chan *walk, *parallel+1)

	for i := 0; i < *parallel+1; i++ {
		var w walk
		pool <- &w
	}

	wg := &sync.WaitGroup{}
	for k := range p.Urls {
		wg.Add(1)
		w := <-pool
		go w.deep(0, k, us, pool, wg)
	}

	wg.Wait()
	us.Save(*outputFile)
}

func P(u string) *parser.Parser {
	p := parser.New(u)
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
		os.Exit(2)
	}
	p.Parse(doc)
	return p
}
