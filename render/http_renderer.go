package render

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	// ContentType header key
	ContentType = "Content-Type"
	// ContentText represents content type text/plain
	ContentText = "text/plain"
	// ContentJSON represents content type application/json
	ContentJSON = "application/json"
	// ContentJSONP represents content type application/javascript
	ContentJSONP = "application/javascript"
	// ContentXML represents content type application/xml
	ContentXML = "application/xml"
	// ContentYAML represents content type application/x-yaml
	ContentYAML = "application/x-yaml"
	// ContentHTML represents content type text/html
	ContentHTML = "text/html"
	// ContentBinary represents content type application/octet-stream
	ContentBinary = "application/octet-stream"

	// ContentDisposition describes contentDisposition
	ContentDisposition = "Content-Disposition"

	// describes content disposition type
	dispositionInline = "inline"
	// describes content disposition type
	dispositionAttachment = "attachment"
)

type HttpOptions struct {
	Options
	// supported content types
	ContentBinary, ContentHTML, ContentXML, ContentText, ContentJSON, ContentJSONP string

	// DefaultCharset content data charset
	DefaultCharset string
	// AppendCharset on response content
	AppendCharset bool
}

// HTTPRenderer definition
type HTTPRenderer struct {
	Renderer
	opts HttpOptions
	// mark init is completed
	initialized bool
}

func NewHTTPRenderer(config ...func(*HttpOptions)) *HTTPRenderer {
	base := New()
	baseOpts := base.opts

	r := &HTTPRenderer{
		Renderer: *base,
		opts: HttpOptions{
			Options: baseOpts,

			ContentXML:    ContentXML,
			ContentText:   ContentText,
			ContentHTML:   ContentHTML,
			ContentJSON:   ContentJSON,
			ContentJSONP:  ContentJSONP,
			ContentBinary: ContentBinary,

			DefaultCharset: defaultCharset,
			AppendCharset:  true,
		},
	}

	// apply user config
	if len(config) > 0 {
		config[0](&r.opts)
	}

	if r.opts.AppendCharset {
		AppendCharset()
	}

	base = nil
	return r
}

// AppendCharset for all content types
func AppendCharset() {

}

// Empty alias method of the NoContent()
func (r *HTTPRenderer) Empty(w http.ResponseWriter) error {
	return r.NoContent(w)
}

// NoContent serve success but no content response
func (r *HTTPRenderer) NoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// String alias method of the Text()
func (r *HTTPRenderer) String(w http.ResponseWriter, status int, v string) error {
	return r.Text(w, status, v)
}

// Text serve string content as text/plain response
func (r *HTTPRenderer) Text(w http.ResponseWriter, status int, v string) error {
	w.WriteHeader(status)
	w.Header().Set(ContentType, "text/plain; charset=UTF-8")
	_, err := w.Write([]byte(v))

	return err
}

// JSON serve string content as json response
func (r *HTTPRenderer) JSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set(ContentType, r.opts.ContentJSON)
	w.WriteHeader(status)

	return r.Renderer.JSON(w, v)
}

// JSONP serve data as JSONP response
func (r *HTTPRenderer) JSONP(w http.ResponseWriter, status int, callback string, v interface{}) error {
	w.Header().Set(ContentType, r.opts.ContentJSONP)
	w.WriteHeader(status)

	bs, err := jsonMarshal(v, false, false)
	if err != nil {
		return err
	}

	if callback == "" {
		return errors.New("renderer: callback can not bet empty")
	}

	w.Write([]byte(callback + "("))
	_, err = w.Write(bs)
	w.Write([]byte(");"))

	return err
}

// XML serve data as XML response
func (r *HTTPRenderer) XML(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set(ContentType, r.opts.ContentXML)
	w.WriteHeader(status)

	return r.Renderer.XML(w, v)
}

// Binary serve data as Binary response.
// usage:
// 		var reader io.Reader
// 		reader, _ = os.Open("./README.md")
// 		r.Binary(w, http.StatusOK, reader, "readme.md", true)
func (r *HTTPRenderer) Binary(w http.ResponseWriter, status int, in io.Reader, outName string, inline bool) error {
	bs, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	dispositionType := dispositionAttachment
	if inline {
		dispositionType = dispositionInline
	}

	w.Header().Set(ContentType, r.opts.ContentBinary)
	w.Header().Set(ContentDisposition, fmt.Sprintf("%s; filename=%s", dispositionType, outName))
	w.WriteHeader(status)

	_, err = w.Write(bs)
	return err
}
