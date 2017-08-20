package gofaster

import (
	"bytes"
	"net"
	"strconv"
)

/***Variables***/

//TODO: Make() -> Make(conn)
//TODO: headers []byte -> headers bytes.Buffer

//TODO: may want to reacreate readmimeheader someotherway

//Response : custom http response object used to handle requests
type Response struct {

	//data objects
	body       []byte       //raw bytes of body
	headers    bytes.Buffer //raw bytes of headers
	statusCode int64        //status code for return

	//status objects
	bIsRaw         bool //turns off status and headers if true to leave raw body
	bConnClose     bool //determines if you connection should be set
	bContentType   bool //determines if content-type has been declared
	bContentLength bool //determines if content-length has been declared

	//buffer objects
	ok      bool   //used to help build message for completed response
	message []byte //used to help build message for completed response
	tbuffer []byte //temporary buffer for transactions
}

/***Methods***/

//(*Response).checkHeaders : update object checks when headers are manually specified
func (r *Response) checkHeaders(key string) {
	if key == "Content-Type" {
		r.bContentType = true
	} else if key == "Content-Length" {
		r.bContentLength = true
	} else if key == "Connection" {
		r.bConnClose = true
	}
}

//(*Response).StatusCode : sets the status code for the response
func (r *Response) StatusCode(code int64) {
	r.statusCode = code
}

//(*Response).SetHeader : write string key and value to headers byte array
func (r *Response) SetHeader(key, value string) {
	r.tbuffer = append(r.tbuffer[:0], key...)
	r.tbuffer = append(r.tbuffer, ':', ' ')
	r.tbuffer = append(r.tbuffer, value...)
	r.tbuffer = append(r.tbuffer, '\r', '\n')
	r.headers.Write(r.tbuffer)
	r.checkHeaders(key)
}

//(*Response).SetHeaderBytes : write raw bytes of key and value to header byte array
func (r *Response) SetHeaderBytes(key, value []byte) {
	r.tbuffer = append(r.tbuffer[:0], key...)
	r.tbuffer = append(r.tbuffer, ':', ' ')
	r.tbuffer = append(r.tbuffer, value...)
	r.tbuffer = append(r.tbuffer, '\r', '\n')
	r.headers.Write(r.tbuffer)
	r.checkHeaders(string(key))
}

//(*Response).SetBody : write string as body to body byte array
func (r *Response) SetBody(body string) {
	r.body = append(r.body[:0], body...)
}

//(*Response).SetBody : write raw bytes as body to body byte array
func (r *Response) SetBodyBytes(body []byte) {
	r.body = append(r.body[:0], body...)
}

//(*Response).SetRaw : write string as entire http response
func (r *Response) SetRaw(body string) {
	r.bIsRaw = true
	r.body = append(r.body[:0], body...)
}

//(*Response).SetRaw : write raw bytes as entire http response
func (r *Response) SetRawBytes(body []byte) {
	r.bIsRaw = true
	r.body = append(r.body[:0], body...)
}

//(*Response).make : build objects into a single response bytearray
func (r *Response) Make(conn net.Conn) {
	//determine if error response
	if !r.bIsRaw {
		// attempt to get message
		if r.message, r.ok = commonStatusCodes[r.statusCode]; !r.ok {
			r.message = allStatusCodes[r.statusCode]
		}
		// write status to tbuffer
		r.tbuffer = append(r.tbuffer[:0], 'H', 'T', 'T', 'P', '/', '1', '.', '1', ' ')
		r.tbuffer = append(r.tbuffer, strconv.FormatInt(r.statusCode, 10)...)
		r.tbuffer = append(r.tbuffer, r.message...)
		// write missing headers to tbuffer
		if !r.bContentType {
			r.tbuffer = append(r.tbuffer, "Content-Type"...)
			r.tbuffer = append(r.tbuffer, ':', ' ')
			r.tbuffer = append(r.tbuffer, "text/plain"...)
			r.tbuffer = append(r.tbuffer, '\r', '\n')
		}
		if !r.bContentLength {
			r.tbuffer = append(r.tbuffer, "Content-Length"...)
			r.tbuffer = append(r.tbuffer, ':', ' ')
			r.tbuffer = append(r.tbuffer, strconv.Itoa(len(r.body))...)
			r.tbuffer = append(r.tbuffer, '\r', '\n')
		}
		if !r.bConnClose {
			r.tbuffer = append(r.tbuffer, "Connection"...)
			r.tbuffer = append(r.tbuffer, ':', ' ')
			r.tbuffer = append(r.tbuffer, "Close"...)
			r.tbuffer = append(r.tbuffer, '\r', '\n')
		}
		// write given headers to tbuffer
		r.tbuffer = append(r.tbuffer, r.headers.Bytes()...)
		r.tbuffer = append(r.tbuffer, '\r', '\n')
		// write body to tbuffer
		r.tbuffer = append(r.tbuffer, r.body...)
	} else {
		// write body to tbuffer
		r.tbuffer = append(r.tbuffer[:0], r.body...)
	}
	conn.Write(r.tbuffer)
}
