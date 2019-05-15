package walker

import (
	"fmt"
	"sync"

	"github.com/sah4ez/sitemap/pkg/parser"
	"github.com/sah4ez/sitemap/pkg/sitemap"
)

type Walk func(int, string, *sitemap.Urlset, chan *Walk, *sync.WaitGroup)

func (w *Walk) Deep(level int, maxDepth *int, urlStr string, us *sitemap.Urlset, pool chan *Walk, wg *sync.WaitGroup) {
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
	p := parser.P(urlStr)
	if p == nil {
		return
	}

	l := level + 1
	for k := range p.Urls {
		w.Deep(l, maxDepth, k, us, nil, nil)
		us.Add(k)
	}
}
