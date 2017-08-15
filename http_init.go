package gofaster

import (
	"log"
	"time"
)

/***Variables***/
var timeFormat = time.RFC1123

var commonStatusCodes = map[int64][]byte{
	200: []byte(" OK\r\n"),
	400: []byte(" Bad Request\r\n"),
	404: []byte(" Not Found\r\n"),
	500: []byte(" Internal Server Error\r\n"),
}
var allStatusCodes = map[int64][]byte{
	100: []byte("Continue\r\n"),
	101: []byte("Switching Protocols\r\n"),
	102: []byte("Processing\r\n"),
	200: []byte("OK\r\n"),
	201: []byte("Created\r\n"),
	202: []byte("Accepted\r\n"),
	203: []byte("Non-Authoritative Information\r\n"),
	204: []byte("No Content\r\n"),
	205: []byte("Reset Content\r\n"),
	206: []byte("Partial Content\r\n"),
	207: []byte("Multi-Status\r\n"),
	208: []byte("Already Reported\r\n"),
	226: []byte("IM Used\r\n"),
	300: []byte("Multiple Choices\r\n"),
	301: []byte("Moved Permanently\r\n"),
	302: []byte("Found\r\n"),
	303: []byte("See Other\r\n"),
	304: []byte("Not Modified\r\n"),
	305: []byte("Use Proxy\r\n"),
	307: []byte("Temporary Redirect\r\n"),
	308: []byte("Permanent Redirect\r\n"),
	400: []byte("Bad Request\r\n"),
	401: []byte("Unauthorized\r\n"),
	402: []byte("Payment Required\r\n"),
	403: []byte("Forbidden\r\n"),
	404: []byte("Not Found\r\n"),
	405: []byte("Method Not Allowed\r\n"),
	406: []byte("Not Acceptable\r\n"),
	407: []byte("Proxy Authentication Required\r\n"),
	408: []byte("Request Timeout\r\n"),
	409: []byte("Conflict\r\n"),
	410: []byte("Gone\r\n"),
	411: []byte("Length Required\r\n"),
	412: []byte("Precondition Failed\r\n"),
	413: []byte("Request Entity Too Large\r\n"),
	414: []byte("Request-URI Too Long\r\n"),
	415: []byte("Unsupported Media Type\r\n"),
	416: []byte("Requested Range Not Satisfiable\r\n"),
	417: []byte("Expectation Failed\r\n"),
	422: []byte("Unprocessable Entity\r\n"),
	423: []byte("Locked\r\n"),
	424: []byte("Failed Dependency\r\n"),
	426: []byte("Upgrade Required\r\n"),
	428: []byte("Precondition Required\r\n"),
	429: []byte("Too Many Requests\r\n"),
	431: []byte("Request Header Fields Too Large\r\n"),
	500: []byte("Internal Server Error\r\n"),
	501: []byte("Not Implemented\r\n"),
	502: []byte("Bad Gateway\r\n"),
	503: []byte("Service Unavailable\r\n"),
	504: []byte("Gateway Timeout\r\n"),
	505: []byte("HTTP Version Not Supported\r\n"),
	506: []byte("Variant Also Negotiates\r\n"),
	507: []byte("Insufficient Storage\r\n"),
	508: []byte("Loop Detected\r\n"),
	510: []byte("Not Extended\r\n"),
	511: []byte("Network Authentication Required\r\n"),
}

/***Functions***/

//GetDate : get current date date (w/ format for http requests)
func GetDate() string {
	return time.Now().Format(timeFormat)
}

//ListenAndServe : serve server FOREVER! (A Really Long Time...)
func ListenAndServe(address string, handler func(*Request, *Response)) {
	var n int = 10
	sp := &sockPool{}
	sp.Spawn("tcp4", address, n, handler)
	log.Printf("- Started Server: %s, with %d Workers\n", address, n)
	sp.Wait()
}
