package output

import (
	// "log"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/jaesbit/ffuf/pkg/ffuf"
)

type Session struct {
	handle    *os.File
	enabled   bool
	size      int
	delimiter string
}

func NewSession(conf *ffuf.Config) *Session {
	var sess Session
	if conf.OutputFormat == "full" {
		f, err := os.OpenFile(conf.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		sess.handle = f
		if err != nil {
			sess.enabled = false
		} else {
			sess.enabled = true
			sess.size = 54
			sess.delimiter = strings.Repeat("~", sess.size) + "\n"
		}

	} else {
		sess.enabled = false
	}
	return &sess

}

func (sess *Session) Close() error {

	if sess.enabled {
		defer sess.handle.Close()
	}
	return nil
}

func (sess *Session) Update(res ffuf.Response) error {

	if sess.enabled {
		sess.handle.WriteString(sess.delimiter)
		sess.WriteResume(&res)
		sess.handle.WriteString("----\n")
		sess.WriteRequest(res.Request)
		sess.handle.WriteString("\n----\n")
		sess.WriteResponse(&res)
		sess.handle.WriteString("\n")
	}
	return nil
}

func (sess *Session) WriteResume(req *ffuf.Response) {
	var data = fmt.Sprintf(
		"%s %v\nMethod: %s, Size: %v, Words: %v\n",
		req.Request.Url,
		req.StatusCode,
		req.Request.Method,
		req.ContentLength,
		req.ContentWords,
	)
	sess.handle.WriteString(data)
}

func (sess *Session) WriteResponse(req *ffuf.Response) {

	var head = fmt.Sprintf("HTTP/1.1 %v %s\n",
		req.StatusCode,
		http.StatusText(int(req.StatusCode)),
	)

	for key, value := range req.Headers {
		head += fmt.Sprintf("%s: %s\n", key, strings.Join(value, " "))

	}
	head += "\n"
	sess.handle.WriteString(head)
	sess.handle.Write(req.Data)
}

func (sess *Session) WriteRequest(req *ffuf.Request) {

	u, err := url.Parse(req.Url)

	if err != nil {
		panic(err)
	}

	var head = fmt.Sprintf("%s %s HTTP/1.1\n", req.Method, u.Path)
	for key, value := range req.Headers {
		head += fmt.Sprintf("%s: %s\n", key, value)
	}

	if !strings.Contains(head, "Host:") {
		var port = ""
		if u.Port() != "" {
			port = fmt.Sprintf(":%v", u.Port())
		}
		head += fmt.Sprintf("Host: %s%s\n", u.Hostname(), port)
	}

	head += "\n"
	sess.handle.WriteString(head)

	if len(req.Data) > 0 {
		sess.handle.Write(req.Data)
	}

}
