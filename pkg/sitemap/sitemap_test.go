package sitemap

import (
	"bytes"
	"encoding/xml"
	"testing"
)

func TestMrashal(t *testing.T) {
	str := `
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
   <url>
      <loc>http://www.example.com/</loc>
   </url>
   <url>
      <loc>http://www.example.com/2/</loc>
   </url>
   <url>
      <loc>http://www.example.com/3/</loc>
   </url>
</urlset>
	`
	exp := Urlset{
		Url: []Url{
			{Loc: "http://www.example.com/"},
			{Loc: "http://www.example.com/2/"},
			{Loc: "http://www.example.com/3/"},
		},
	}
	sm := Urlset{}
	err := xml.Unmarshal([]byte(str), &sm)
	if err != nil {
		t.Fatal(err)
	}
	act, _ := xml.Marshal(&sm)
	exp2, _ := xml.Marshal(&exp)
	if !bytes.Equal(act, exp2) {
		t.Errorf("not equal")
	}
}
