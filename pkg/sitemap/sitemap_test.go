package sitemap

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"testing"

	"github.com/sah4ez/sitemap/pkg/util"
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
	exp := New()
	exp.Add("http://www.example.com/")
	exp.Add("http://www.example.com/2/")
	exp.Add("http://www.example.com/3/")

	sm := New()
	err := xml.Unmarshal([]byte(str), &sm)
	if err != nil {
		t.Fatal(err)
	}
	act, _ := xml.Marshal(&sm)
	exp2, _ := xml.Marshal(&exp)
	if !bytes.Equal(act, exp2) {
		t.Errorf("not equal \n%s\n%s", string(act), string(exp2))
	}
}

func TestSave(t *testing.T) {
	str := `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>http://www.example.com/</loc></url><url><loc>http://www.example.com/2/</loc></url><url><loc>http://www.example.com/3/</loc></url></urlset>`
	sm := New()
	err := xml.Unmarshal([]byte(str), &sm)
	if err != nil {
		t.Fatal(err)
	}
	path := "/tmp/sitemap.xml"
	sm.Save(path)
	b, _ := ioutil.ReadFile(path)
	pretty, err := util.Prettify(str, Offset)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal([]byte(pretty), b) {
		t.Errorf("not equal \n%s\n%s", pretty, string(b))
	}
}
