package main

import (
	"net/http"

	"github.com/imgurbot12/gofaster"
	"github.com/valyala/fasthttp"
)

var body = []byte("Hello World!")

func runGoFaster() {
	gofaster.ListenAndServe(":8080", func(_ *gofaster.Request, resp *gofaster.Response) {
		resp.SetHeader("Server", "nginx")
		resp.SetHeader("Date", "Tue, 15 Aug 2017 02:57:36 GMT")
		resp.SetHeader("Content-Type", "text/plain")
		resp.SetHeader("Content-Length", "12")
		resp.SetHeader("Connection", "close")
		resp.SetHeader("Accept-Ranges", "bytes")
		resp.SetBodyBytes(body)
	})
}

func runFastHttp() {
	fasthttp.ListenAndServe(":8080", func(ctx *fasthttp.RequestCtx) {
		// ctx.Response.Header.Set("Server", "nginx")
		// ctx.Response.Header.Set("Date", "Tue, 15 Aug 2017 02:57:36 GMT")
		// ctx.Response.Header.Set("Content-Type", "text/plain")
		// ctx.Response.Header.Set("Content-Length", "12")
		// ctx.Response.Header.Set("Connection", "close")
		// ctx.Response.Header.Set("Accept-Ranges", "bytes")
		ctx.Response.SetBody(body)
	})
}

func runNetHttp() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Server", "nginx")
		w.Header().Set("Date", "Tue, 15 Aug 2017 02:57:36 GMT")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", "12")
		w.Header().Set("Connection", "close")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Write(body)
	})
	http.ListenAndServe(":8080", nil)
}

func main() {
	// runGoFaster()
	// runFastHttp()
	runNetHttp()
}
