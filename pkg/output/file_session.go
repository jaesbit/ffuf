package output

import (
	// "log"
	"os"
	"fmt"
	"strings"

	"ffuf/pkg/ffuf"
)

type Session struct {
	handle  *os.File
	enabled bool
	size  int
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

		sess.WriteRequest(res.Request, res.StatusCode)
		sess.WriteResponse(&res)


	}
	return nil
}

func (sess *Session) WriteResponse(req *ffuf.Response) {

	var head =  ""
	for key, value := range req.Headers {
		head += fmt.Sprintf("%s: %s\n", key, value)

	}
	head += "\n\n"
	if _, err := sess.handle.WriteString(head); err != nil {
	}
	if _, err := sess.handle.Write(req.Data); err != nil {
	}
	// head = strings.Repeat("*", sess.size)
	head = "\n\n"
	if _, err := sess.handle.WriteString(head); err != nil {
	}
}


func (sess *Session) WriteRequest(req *ffuf.Request, status int64) {

	var head = strings.Repeat("~", sess.size)
	head += fmt.Sprintf("\n%s %s %v\n", req.Method, req.Url, status)
	for key, value := range req.Headers {
		head += fmt.Sprintf("%s: %s\n", key, value)

	}
	head += strings.Repeat("~", sess.size)
	head += "\n\n"
	if _, err := sess.handle.WriteString(head); err != nil {
	}

}

