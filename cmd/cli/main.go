package main

import (
	"flag"
	"os"
	"strings"
	"sync"

	"github.com/sah4ez/sitemap/pkg/parser"
	"github.com/sah4ez/sitemap/pkg/sitemap"
	"github.com/sah4ez/sitemap/pkg/walker"
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

func main() {
	urlStr := os.Args[len(os.Args)-1]
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	p := parser.P(urlStr)

	us := sitemap.New()

	pool := make(chan *walker.Walk, *parallel+1)

	for i := 0; i < *parallel+1; i++ {
		var w walker.Walk
		pool <- &w
	}

	wg := &sync.WaitGroup{}
	for k := range p.Urls {
		wg.Add(1)
		w := <-pool
		go w.Deep(0, maxDepth, k, us, pool, wg)
	}

	wg.Wait()
	us.Save(*outputFile)
}
