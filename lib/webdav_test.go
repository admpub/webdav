package lib

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/admpub/pp/ppnocolor"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/webdav"
)

func TestWebDAVHandler(t *testing.T) {
	os.MkdirAll(`./testdata`, os.ModePerm)
	os.MkdirAll(`./testdata/a/music`, os.ModePerm)
	os.MkdirAll(`./testdata/a/video`, os.ModePerm)
	os.MkdirAll(`./testdata/a/book`, os.ModePerm)
	os.MkdirAll(`./testdata/b`, os.ModePerm)
	os.MkdirAll(`./testdata/c`, os.ModePerm)
	os.MkdirAll(`./testdata/d`, os.ModePerm)
	noSniff := true
	//defer os.RemoveAll(`./testdata`)
	wfs := WebDavDir{
		Dir: webdav.Dir(`./testdata/`),
		User: &User{
			Username: `Test`,
			Password: `TEST`,
			Scope:    `/`,
			Rules: []*Rule{
				{
					Allow: true,
					Path:  `/a/music`,
				},
				{
					Allow: false,
					Path:  `/a/video`,
				},
				{
					Allow: false,
					Path:  `/a/book`,
				},
				{
					Allow: true,
					Path:  `/d/`,
				},
			},
			Handler: nil,
		},
		NoSniff: noSniff,
	}
	uh := &webdav.Handler{
		Prefix: ``,
		FileSystem: WebDavDir{
			Dir:     webdav.Dir(wfs.User.Scope),
			User:    wfs.User,
			NoSniff: noSniff,
		},
		LockSystem: webdav.NewMemLS(),
	}
	wfs.User.Handler = uh
	h := webdav.Handler{
		Prefix:     ``,
		FileSystem: wfs,
		LockSystem: webdav.NewMemLS(),
	}

	req, err := http.NewRequest(`PROPFIND`, `/a`, nil)
	assert.NoError(t, err)
	req.Header.Add("Depth", "1")
	req.Body = io.NopCloser(bytes.NewReader(nil))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, webdav.StatusMulti, rec.Code)
	result, err := parseXML(rec.Body)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	ppnocolor.Println(result)
	//rec.Code, rec.Body.String(),rec.Header
}

type props struct {
	Status      string   `xml:"DAV: status"`
	Name        string   `xml:"DAV: prop>displayname,omitempty"`
	Type        xml.Name `xml:"DAV: prop>resourcetype>collection,omitempty"`
	Size        string   `xml:"DAV: prop>getcontentlength,omitempty"`
	ContentType string   `xml:"DAV: prop>getcontenttype,omitempty"`
	ETag        string   `xml:"DAV: prop>getetag,omitempty"`
	Modified    string   `xml:"DAV: prop>getlastmodified,omitempty"`
}

type response struct {
	Href  string  `xml:"DAV: href"`
	Props []props `xml:"DAV: propstat"`
}

func parseXML(data io.Reader) ([]*response, error) {
	decoder := xml.NewDecoder(data)
	result := []*response{}
	for t, _ := decoder.Token(); t != nil; t, _ = decoder.Token() {
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "response" {
				resp := &response{}
				if err := decoder.DecodeElement(resp, &se); err == nil {
					result = append(result, resp)
				}
			}
		}
	}
	return result, nil
}
