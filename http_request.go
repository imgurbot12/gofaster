package gofaster

/***Functions***/
import (
	"fmt"
	"net"
	"net/textproto"
	"net/url"
	"strings"
)

/***Variables***/
var buf = make([]byte, 1024)

//Request : basic http request object
type Request struct {

	// complete data objects
	Method     string
	Protocol   string
	RequestURI string
	Headers    textproto.MIMEHeader

	// optional data objects
	Form  url.Values
	Query url.Values

	//connection reader
	conn    net.Conn
	sBuffer *textproto.Reader
}

//badRequestError : custom error sent during request parsing error
type badRequestError struct {
	what string
	str  string
}

/***Functions***/

//(*badRequestError).Error : required function
func (e *badRequestError) Error() string {
	return fmt.Sprintf("%s %q", e.what, e.str)
}

//parseRequestLine : parse request line from http-request
func parseRequestLine(line string) (method, uri, proto string, ok bool) {
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	if s1 < 0 || s2 < 0 {
		return
	}
	s2 += s1 + 1
	return line[:s1], line[s1+1 : s2], line[s2+1:], true
}

//parseHTTPVersion : detect invalid http-version
func parseHTTPVersion(version string) bool {
	if version == "HTTP/1.0" || version == "HTTP/1.1" {
		return true
	}
	return false
}

/***Methods***/

//(*Request).parseRequest : parse request and set results
func (req *Request) parseRequest() error {
	//temporary varaibles
	var ok bool
	var err error
	var line string
	//attempt to read first line of request
	line, err = req.sBuffer.ReadLine()
	if err != nil {
		return &badRequestError{"Buffer Readline Error!", line}
	}
	//parse first line
	if req.Method, req.RequestURI, req.Protocol, ok = parseRequestLine(line); !ok {
		return &badRequestError{"Malformed HTTP Request!", line}
	}
	//check http-version before moving on
	if ok = parseHTTPVersion(req.Protocol); !ok {
		return &badRequestError{"Malformed HTTP Version!", req.Protocol}
	}
	//get headers
	if req.Headers, err = req.sBuffer.ReadMIMEHeader(); err != nil {
		return err
	}
	return nil
}

//(*Request).RemoteAddr : collect remote address
func (req *Request) RemoteAddr() string {
	return req.conn.RemoteAddr().String()
}

//(*Request).ParseQuery : parses query perameters and appends it (*Request).Query
func (req *Request) ParseQuery() error {
	//pre-generate temp variable
	var err error
	// attempt to index place where query starts
	n := strings.Index(req.RequestURI, "?")
	// if n == 0; '?' was not found
	if n == 0 {
		return nil
	}
	// attempt to parse query from RequestURI
	req.Query, err = url.ParseQuery(req.RequestURI[n+1:])
	return err
}

//(*Request).ParseForm : parses form and appends it to (*Request).Form
func (req *Request) ParseForm() error {
	// if method is not POST, exit
	if req.Method != "POST" {
		return &badRequestError{"Not a POST Request!", req.Method}
	}
	// attempt to read-data
	len, err := req.sBuffer.R.Read(buf)
	// if line is read -> parse-query
	if err == nil {
		req.Form, err = url.ParseQuery(string(buf[:len]))
	}
	return err
}
