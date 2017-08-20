package gofaster

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/OneOfOne/xxhash"
)

/***Functions***/

//statFile : collect stats of a file and check that ftype is correct
func statFile(fname string, isdir bool) os.FileInfo {
	var err error
	var fdata os.FileInfo
	// check if file exists before creating page-handler
	if fdata, err = os.Stat(fname); os.IsNotExist(err) {
		log.Fatalf("- ServeStatic: %q does NOT exist!\n", fname)
	}
	// check if file should not be directory
	if fdata.IsDir() && !isdir {
		log.Fatalf("- ServeFile: %q is a Directory!\n", fname)
	}
	// check if file should be a directory
	if !fdata.IsDir() && isdir {
		log.Fatalf("- ServeFile: %q is NOT a Directory!\n", fname)
	}
	return fdata
}

func walkDir(dname string) []string {
	//list for all found files0
	var fileList []string
	// walk files
	filepath.Walk(dname, func(path string, f os.FileInfo, err error) error {
		if err == nil {
			if f.Mode().IsRegular() {
				fileList = append(fileList, path)
			}
		}
		return nil
	})
	return fileList
}

/***Handlers***/

//EtagHandler : supports response caching via etags (https://en.wikipedia.org/wiki/HTTP_ETag)
func EtagHandler(handler func(*Request, *Response)) func(*Request, *Response) {
	return func(req *Request, resp *Response) {
		//run original handler
		handler(req, resp)
		//generate etag as string
		hasher := xxhash.New64()
		hasher.Write(resp.body)
		etag := "W/" + strconv.FormatUint(hasher.Sum64(), 10)
		//check if etag header exists
		if etag == req.Headers.Get("If-None-Match") {
			resp.StatusCode(304)
			resp.SetBody("")
		} else {
			resp.SetHeader("Etag", etag)
			resp.SetHeader("Accept-Ranges", "bytes")
		}
	}
}

//ServeFile : serves static file content to
func ServeFile(fname string) func(*Request, *Response) {
	//variables
	var mimetype string
	var lastmodified string

	// collect file stats
	fdata := statFile(fname, false)
	// pre-calculate function varaibles
	ext := filepath.Ext(fname)
	mimetype = mime.TypeByExtension(ext)
	if n := strings.Index(mimetype, ";"); n > 0 {
		mimetype = mimetype[:n] //ignore charset returned
	}
	lastmodified = fdata.ModTime().Format(timeFormat)

	// generate function
	return func(req *Request, resp *Response) {
		content, err := ioutil.ReadFile(fname)
		if err == nil {
			resp.SetHeader("Server", "nginx")
			resp.SetHeader("Date", GetDate())
			resp.SetHeader("Content-Type", mimetype)
			resp.SetHeader("Content-Length", strconv.Itoa(len(content)))
			resp.SetHeader("Last-Modified", lastmodified)
			resp.SetHeader("Connection", "close")
			resp.SetBodyBytes(content)
		} else {
			resp.StatusCode(500)
		}
	}
}

func ServeDir(dname string) func(*Request, *Response) {
	// collect stats and confirm is-dir
	statFile(dname, true)
	files := walkDir(dname)
	// create muxer for all files given
	mux := NewServeMux()
	for _, fname := range files {
		mux.HandleFunc("/"+fname, ServeFile(fname))
	}
	return mux.ServeHTTP()
}

func StripPrefix(prefix string, handler func(*Request, *Response)) func(*Request, *Response) {
	// check for nil prefix
	if prefix == "" {
		return handler
	}
	// create handler to trip the prefix before running handler
	return func(req *Request, resp *Response) {
		req.RequestURI = strings.TrimPrefix(req.RequestURI, prefix)
		handler(req, resp)
	}
}

//CompressHandler : compresses data and transfers it under certain conditions
func CompressHandler(handler func(*Request, *Response)) func(*Request, *Response) {
	return func(req *Request, resp *Response) {
		//run original handler
		handler(req, resp)
		//if compression is accepted
		if req.Headers.Get("Accept-Encoding") != "" {
			//if the data is large enough to be worth compressing
			if len(resp.body) > 1000 {
				//compress data
				var b bytes.Buffer
				writer := gzip.NewWriter(&b)
				writer.Write(resp.body)
				writer.Close()
				//add compression header
				resp.SetHeader("Content-Encoding", "gzip")
				//replace body with compressed form
				resp.body = b.Bytes()
			}
		}
	}
}
