package handlers

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/gookit/sux"
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

// Status get status code
func (w *statusWriter) Status() int {
	return w.status
}

// WriteHeader write status code
func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Write data to resp
func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

// RequestLogger middleware
func RequestLogger() sux.HandlerFunc {
	skip := map[string]int{
		// "/": 1,
		"/health": 1,
		"/status": 1,
	}

	return func(c *sux.Context) {
		// start time
		start := time.Now()

		// rewrite the resp
		sw := &statusWriter{ResponseWriter: c.Resp}
		c.Resp = sw

		// Process request
		c.Next()

		path := c.URL().Path

		// Log only when path is not being skipped
		if _, ok := skip[path]; ok {
			return
		}

		// log post/put data
		// postData := ""
		// if c.Req.Method != "GET" {
		// 	buf, _ := c.RawData()
		// 	postData = string(buf)
		// }

		mColor := colorForMethod(c.Req.Method)
		codeColor := colorForStatus(sw.Status())

		fmt.Printf(
			// 2006-01-02 15:04:05 [sux] GET /articles 200 10.0.0.1 "use-agent" 0.034ms
			// `%s %s %s %d %s "%s" %sms` + "\n",
			"%s %s %s [%s] %s %sms\n",
			start.Format("2006/01/02 15:04:05"),
			c.ClientIP(),
			mColor.Render(c.Req.Method),
			codeColor.Render(sw.Status()),
			c.Req.RequestURI,
			// c.Header("User-Agent"),
			calcElapsedTime(start),
		)
	}
}

// calcElapsedTime 计算运行时间消耗 单位 ms(毫秒)
func calcElapsedTime(startTime time.Time) string {
	return fmt.Sprintf("%.3f", time.Since(startTime).Seconds()*1000)
}

func colorForStatus(code int) color.Color {
	switch {
	case code >= 200 && code < 300:
		return color.FgGreen
	case code >= 300 && code < 400:
		return color.FgCyan
	case code >= 400 && code < 500:
		return color.FgYellow
	default:
		return color.FgRed
	}
}

func colorForMethod(method string) color.Color {
	switch method {
	case "GET":
		return color.FgBlue
	case "POST":
		return color.FgCyan
	case "PUT":
		return color.FgYellow
	case "DELETE":
		return color.FgRed
	case "PATCH":
		return color.FgGreen
	case "HEAD":
		return color.FgMagenta
	case "OPTIONS":
		return color.FgWhite
	default:
		return color.FgDefault
	}
}