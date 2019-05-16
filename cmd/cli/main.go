package main

import (
	"flag"
	"fmt"
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
	us := sitemap.New()

	urlStr := os.Args[len(os.Args)-1]
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}
	parsers := parser.Pool(*parallel + 1)

	urls := parser.P(parsers, urlStr)
	if urls == nil {
		fmt.Println("could not get parser")
		os.Exit(2)
	}

	walkers := walker.Pool((2 * *parallel) + 1)

	wg := &sync.WaitGroup{}
	for k := range urls {
		wg.Add(1)
		w := <-walkers
		go w.Deep(0, maxDepth, k, us, walkers, parsers, wg)
	}

	wg.Wait()
	us.Save(*outputFile)
}
