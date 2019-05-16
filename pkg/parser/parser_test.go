package parser

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/html"
)

func TestSuccessParse(t *testing.T) {
	data := `
<html>
<title>
title
</title>
<body>
<base href="http://example.com"/>
<a href="/hello"/>
<a href="/hello2"/>
<a href="file:///home/sah4ez/go/src/github.com/sah4ez/sitemap/pkg/parser/test/data2.html"/>
<a href="mailto:a@emxaple.com"/>
<a href="/"/>
<b href="/b"/>
</body>
</html>
`
	doc, err := html.Parse(strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	exp := map[string]struct{}{
		"/hello":  struct{}{},
		"/hello2": struct{}{},
		"/":       struct{}{},
		"file:///home/sah4ez/go/src/github.com/sah4ez/sitemap/pkg/parser/test/data2.html": struct{}{},
	}

	p := New()
	p.Parse(doc)
	for k := range p.Urls {
		if _, ok := exp[k]; !ok {
			t.Errorf("path not contains: %s", k)
		}
	}
	tag := "http://example.com"
	if p.Base != tag {
		t.Errorf("invalid base tag: %s exp: %s", p.Base, tag)
	}
}

func TestPool(t *testing.T) {
	srv := &http.Server{Addr: "127.0.0.1:10000"}
	go srv.ListenAndServe()
	defer srv.Shutdown(context.TODO())

	http.HandleFunc("/1", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

	time.Sleep(1 * time.Second)

	pool := Pool(1)
	P(pool, "http://127.0.0.1:10000/1")
	select {
	case <-time.After(20 * time.Second):
		t.Errorf("timeout return parser to pool")
	case _ = <-pool:
		t.Log("parser returned to pool")
	}
}

func TestPoolNotReturn(t *testing.T) {
	srv := &http.Server{Addr: "127.0.0.1:10000"}
	go srv.ListenAndServe()
	defer srv.Shutdown(context.TODO())

	http.HandleFunc("/2", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

	time.Sleep(1 * time.Second)

	pool := Pool(1)
	P(pool, "http://127.0.0.1:10000/2")
	<-pool
	select {
	case <-time.After(1 * time.Second):
		t.Log("timeout return parser to pool it's ok")
	case _ = <-pool:
		t.Error("parser returned to pool")
	}
}
