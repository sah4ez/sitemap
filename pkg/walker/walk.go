package walker

import (
	"fmt"
	"sync"

	"github.com/sah4ez/sitemap/pkg/parser"
	"github.com/sah4ez/sitemap/pkg/sitemap"
)

type Walk struct{}

func (w Walk) Deep(level int, maxDepth int, urlStr string, us *sitemap.Urlset, walkers chan Walk, parsers chan *parser.Parser, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
		if walkers != nil {
			walkers <- w
		}
	}()

	fmt.Println("level", level, "depth", maxDepth)
	if level >= maxDepth {
		return
	}
	urls := parser.P(parsers, urlStr)
	if urls == nil {
		return
	}

	l := level + 1
	for k := range urls {
		w.Deep(l, maxDepth, k, us, nil, parsers, nil)
		us.Add(k)
	}
}

func Pool(size int) chan Walk {
	pool := make(chan Walk, size)

	for i := 0; i < size; i++ {
		pool <- Walk{}
	}
	return pool
}
