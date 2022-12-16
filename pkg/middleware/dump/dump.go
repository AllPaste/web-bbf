package dump

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

func resetStrB(strB *strings.Builder) {
	strB.Reset()
}

func Dump(opts ...Option) gin.HandlerFunc {
	strBPool := sync.Pool{
		New: func() interface{} {
			return new(strings.Builder)
		},
	}

	o := &options{
		request:      true,
		response:     true,
		body:         true,
		headers:      true,
		cookies:      true,
		convertBytes: nil,
	}

	for _, opt := range opts {
		opt(o)
	}

	headerHiddenFields := make([]string, 0)
	bodyHiddenFields := make([]string, 0)

	if !o.cookies {
		headerHiddenFields = append(headerHiddenFields, "cookie")
	}

	return func(ctx *gin.Context) {
		strB := strBPool.Get().(*strings.Builder)
		if o.request {
			if o.headers {
				// dump req header
				s, err := FormatToBeautifulJson(ctx.Request.Header, headerHiddenFields)

				if err != nil {
					strB.WriteString(fmt.Sprintf("\nparse req header err \n" + err.Error()))
				} else {
					strB.WriteString("Request-Header:\n")
					strB.WriteString(string(s))
				}
			}
			if o.body {
				// dump req body
				if ctx.Request.ContentLength > 0 {
					buf, err := ioutil.ReadAll(ctx.Request.Body)
					if err != nil {
						strB.WriteString(fmt.Sprintf("\nread bodyCache err \n %s", err.Error()))
						ctx.Writer = &bodyWriter{
							bodyCache:      bytes.NewBufferString(""),
							ResponseWriter: ctx.Writer,
						}
						ctx.Next()
					}
					rdr := ioutil.NopCloser(bytes.NewBuffer(buf))
					ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
					ctGet := ctx.Request.Header.Get("Content-Type")
					ct, _, err := mime.ParseMediaType(ctGet)
					if err != nil {
						strB.WriteString(
							fmt.Sprintf("\ncontent_type: %s parse err \n %s", ctGet, err.Error()),
						)
						ctx.Writer = &bodyWriter{
							bodyCache:      bytes.NewBufferString(""),
							ResponseWriter: ctx.Writer,
						}
						ctx.Next()
					}

					switch ct {
					case gin.MIMEJSON:
						bts, err := ioutil.ReadAll(rdr)
						if err != nil {
							strB.WriteString(fmt.Sprintf("\nread rdr err \n %s", err.Error()))
							ctx.Writer = &bodyWriter{
								bodyCache:      bytes.NewBufferString(""),
								ResponseWriter: ctx.Writer,
							}
						}

						s, err := BeautifyJsonBytes(bts, bodyHiddenFields)
						if err != nil {
							strB.WriteString(fmt.Sprintf("\nparse req body err \n" + err.Error()))
							ctx.Writer = &bodyWriter{
								bodyCache:      bytes.NewBufferString(""),
								ResponseWriter: ctx.Writer,
							}
						}

						strB.WriteString("\nRequest-Body:\n")
						strB.WriteString(string(s))
					case gin.MIMEPOSTForm:
						bts, err := ioutil.ReadAll(rdr)
						if err != nil {
							strB.WriteString(fmt.Sprintf("\nread rdr err \n %s", err.Error()))
							ctx.Writer = &bodyWriter{
								bodyCache:      bytes.NewBufferString(""),
								ResponseWriter: ctx.Writer,
							}
							ctx.Next()
						}
						val, err := url.ParseQuery(string(bts))

						s, err := FormatToBeautifulJson(val, bodyHiddenFields)
						if err != nil {
							strB.WriteString(fmt.Sprintf("\nparse req body err \n" + err.Error()))
							ctx.Writer = &bodyWriter{
								bodyCache:      bytes.NewBufferString(""),
								ResponseWriter: ctx.Writer,
							}
						}
						strB.WriteString("\nRequest-Body:\n")
						strB.WriteString(string(s))
					case gin.MIMEMultipartPOSTForm:
					default:
					}
				}
			}
		}
		if o.response {
			if o.headers {
				// dump res header
				sHeader, err := FormatToBeautifulJson(ctx.Writer.Header(), headerHiddenFields)
				if err != nil {
					strB.WriteString(fmt.Sprintf("\nparse res header err \n" + err.Error()))
				} else {
					strB.WriteString("\nResponse-Header:\n")
					strB.WriteString(string(sHeader))
				}
			}

			if o.body {
				bw, ok := ctx.Writer.(*bodyWriter)
				if !ok || bw == nil {
					strB.WriteString("\nbodyWriter was override or nil, can not read bodyCache")
					if o.convertBytes != nil {
						o.convertBytes(strB.String())
					} else {
						fmt.Println(strB.String())
					}
					resetStrB(strB)
					strBPool.Put(strB)
					return
				}
				// dump res body
				if bodyAllowedForStatus(ctx.Writer.Status()) && bw.bodyCache.Len() > 0 {
					ctGet := ctx.Writer.Header().Get("Content-Type")
					ct, _, err := mime.ParseMediaType(ctGet)
					if err != nil {
						strB.WriteString(
							fmt.Sprintf("\ncontent-type: %s parse  err \n %s", ctGet, err.Error()),
						)
						if o.convertBytes != nil {
							o.convertBytes(strB.String())
						} else {
							fmt.Println(strB.String())
						}
					}
					switch ct {
					case gin.MIMEJSON:
						s, err := BeautifyJsonBytes(bw.bodyCache.Bytes(), bodyHiddenFields)
						if err != nil {
							strB.WriteString(fmt.Sprintf("\nparse bodyCache err \n" + err.Error()))
							if o.convertBytes != nil {
								o.convertBytes(strB.String())
							} else {
								fmt.Println(strB.String())
							}
						}
						strB.WriteString("\nResponse-Body:\n")

						strB.WriteString(string(s))
					case gin.MIMEHTML:
					default:
					}
				}
				strBPool.Put(strB)
				ctx.Next()
			}
		}
	}
}

type bodyWriter struct {
	gin.ResponseWriter
	bodyCache *bytes.Buffer
}

// rewrite Write()
func (w bodyWriter) Write(b []byte) (int, error) {
	w.bodyCache.Write(b)
	return w.ResponseWriter.Write(b)
}

// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}
