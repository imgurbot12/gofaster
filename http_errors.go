package gofaster

import "log"

/***Variables***/
var _400 = []byte("<html>\r\n<head><title>400 Bad Request</title></head>\r\n<body bgcolor=\"white\">\r\n<center><h1>400 Bad Request</h1></center>\r\n<hr><center>nginx</center>\r\n</body>\r\n</html>\r\n")
var _404 = []byte("<html>\r\n<head><title>404 Not Found</title></head>\r\n<body bgcolor=\"white\">\r\n<center><h1>404 Not Found</h1></center>\r\n<hr><center>nginx</center>\r\n</body>\r\n</html>\r\n")

/***Functions***/

//errorHandler : base error handler that then responds with the best response accordingly
func errorHandler(err int64, req *Request, resp *Response) {
	switch err {
	case 400:
		badRequest(req, resp)
	case 404:
		request404(req, resp)
	}
}

//badRequest : handle bad requests
func badRequest(req *Request, resp *Response) {
	log.Printf("Request - 400 - %-15s %s\n", req.RemoteAddr(), req.RequestURI)
	resp.SetRawBytes(_400)
}

//request404 : handles 404 response
func request404(req *Request, resp *Response) {
	log.Printf("Request - 404 - %-15s %s\n", req.RemoteAddr(), req.RequestURI)
	resp.StatusCode(404)
	resp.SetHeader("Server", "nginx")
	resp.SetHeader("Date", GetDate())
	resp.SetHeader("Content-Type", "text/html")
	resp.SetHeader("Content-Length", "162")
	resp.SetHeader("Connection", "close")
	resp.SetBodyBytes(_404)
}
