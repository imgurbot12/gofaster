package gofaster

/***Variables***/
var _400 = []byte("<html>\r\n<head><title>400 Bad Request</title></head>\r\n<body bgcolor=\"white\">\r\n<center><h1>400 Bad Request</h1></center>\r\n<hr><center>nginx</center>\r\n</body>\r\n</html>\r\n")
var _404 = []byte("<html>\r\n<head><title>404 Not Found</title></head>\r\n<body bgcolor=\"white\">\r\n<center><h1>404 Not Found</h1></center>\r\n<hr><center>nginx</center>\r\n</body>\r\n</html>\r\n")

/***Functions***/

//errorHandler : base error handler that then responds with the best response accordingly
func errorHandler(err string, response *Response) {
	switch err {
	case "Bad Request":
		badRequest(response)
	case "404":
		request404(response)
	}
}

//badRequest : handle bad requests
func badRequest(resp *Response) {
	resp.SetRawBytes(_400)
}

func request404(resp *Response) {
	resp.StatusCode(400)
	resp.SetHeader("Server", "nginx")
	resp.SetHeader("Date", GetDate())
	resp.SetHeader("Content-Type", "text/html")
	resp.SetHeader("Content-Length", "162")
	resp.SetHeader("Connection", "close")
	resp.SetBodyBytes(_404)
}
