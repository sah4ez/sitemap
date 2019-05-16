package sitemap

import (
	"encoding/xml"
	"fmt"
	"html"
	"io/ioutil"
	"sync"

	"github.com/sah4ez/sitemap/pkg/util"
)

const (
	Header = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	Offset = "    "
)

type Urlset struct {
	l sync.RWMutex

	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Url     []Url    `xml:"url"`

	set map[string]bool
}

type Url struct {
	Loc string `xml:"loc"`
}

func (us *Urlset) Add(u string) {
	us.l.Lock()
	defer us.l.Unlock()

	if _, ok := us.set[u]; !ok {
		us.set[u] = true
		us.Url = append(us.Url, Url{Loc: html.EscapeString(u)})
	}
}

func (us *Urlset) Save(path string) error {
	us.l.RLock()
	defer us.l.RUnlock()

	b, err := xml.Marshal(us)
	if err != nil {
		return fmt.Errorf("marshalling error: %s", err)
	}
	pretty, _ := util.Prettify(string(b), Offset)
	pretty = Header + pretty
	err = ioutil.WriteFile(path, []byte(pretty), 0644)
	if err != nil {
		return fmt.Errorf("marshalling error: %s", err)
	}
	return nil
}

func New() *Urlset {
	return &Urlset{
		l:     sync.RWMutex{},
		Url:   make([]Url, 0),
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		set:   map[string]bool{},
	}
}
