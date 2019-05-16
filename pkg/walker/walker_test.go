package walker

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/sah4ez/sitemap/pkg/parser"
	"github.com/sah4ez/sitemap/pkg/sitemap"
)

func TestDeep(t *testing.T) {
	data := `
<html>
<title>
title
</title>
<body>
<a href="/hello"/>
</body>
</html>
`
	srv := &http.Server{Addr: "127.0.0.1:10000"}
	go srv.ListenAndServe()
	defer srv.Shutdown(context.TODO())

	http.HandleFunc("/3", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(data))
	})

	time.Sleep(1 * time.Second)
	parsers := parser.Pool(1)
	walkers := Pool(1)

	us := sitemap.New()

	w := <-walkers
	wg := &sync.WaitGroup{}
	wg.Add(1)
	w.Deep(0, 1, "http://127.0.0.1:10000/3", us, walkers, parsers, wg)
	wg.Wait()
	select {
	case <-time.After(20 * time.Second):
		t.Errorf("timeout return parser to pool")
	case _ = <-walkers:
		t.Log("walker returned to pool")
	}

}
