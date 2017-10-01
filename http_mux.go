package gofaster

import (
	"log"
	"sort"
)

/***Variables***/

//muxEntry : entry for handler used to match given pattern to function
type muxEntry struct {
	pattern string                    //pattern for handler to match too
	Handler func(*Request, *Response) //function request is passed to after match
	isDir   bool                      //stored variable used to match pattern faster
	length  int                       //stored variable used to match pattern faster
}

//dedicated muxEntry slice with sorting capabilities
type muxEntries []*muxEntry

//holds pattern map for all mux-entries
type muxer struct {
	patterns map[string]int //list of existing patterns
	entries  muxEntries     //entries used hold handler and match patterns
}

/***Functions***/

//NewServeMux : spawn muxer instance to handle requests
func NewServeMux() *muxer {
	return &muxer{
		patterns: make(map[string]int),
		entries:  muxEntries{},
	}
}

//Len : used to sort muxEntries
func (slice muxEntries) Len() int {
	return len(slice)
}

//Less : used to sort muxEntries
func (slice muxEntries) Less(i, j int) bool {
	return slice[i].length > slice[j].length
}

//Swap : used to sort muxEntries
func (slice muxEntries) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

/***Methods***/

//(*muxEntry).match : attempt to match pattern and path based on muxEntry
func (e *muxEntry) match(pattern, path string) bool {
	// match if directory
	if e.isDir {
		return len(path) >= e.length && path[0:e.length] == pattern
	}
	// match if other
	return pattern == path
}

//(*muxer).Match : match requested pattern and return the corresponding handler
func (m *muxer) match(path string) func(*Request, *Response) {
	//pre-defined handler for loop exposure
	var handler func(*Request, *Response)
	// iterate patterns to find best result
	for _, v := range m.entries {
		if !v.match(v.pattern, path) {
			continue
		}
		handler = v.Handler
		break
	}
	return handler
}

//(*muxer).serve : matches path and runs corresponding handler
func (m *muxer) serve(req *Request, resp *Response) {
	// attempt to match
	handler := m.match(req.RequestURI)
	// run hander or 404 depending on results of match
	if handler != nil {
		handler(req, resp)
	} else {
		errorHandler(404, req, resp)
	}
}

//(*muxer).ServeHTTP : matches path and runs corresponding handler
func (m *muxer) ServeHTTP() func(*Request, *Response) {
	if len(m.entries) == 0 {
		log.Fatalln("Mux Error! There are no existing Patterns!")
	}
	return m.serve
}

//(*muxer).HandleFunc : add handler to muxer under given pattern
func (m *muxer) HandleFunc(pattern string, handler func(*Request, *Response)) {
	// check for possible errors
	if pattern == "" {
		log.Fatalf("Mux ERROR! Invalid Pattern: %q\n", pattern)
	}
	if handler == nil {
		log.Fatalln("Mux ERROR! Handler is NULL")
	}
	if _, ok := m.patterns[pattern]; ok {
		log.Fatalln("Mux ERROR! This Pattern already Exists!")
	}
	// pre-calculate if pattern is a directory
	var isDir bool
	n := len(pattern)
	if pattern[n-1] == '/' {
		isDir = true
	}
	// append pattern-entry to current slice
	m.entries = append(m.entries, &muxEntry{pattern: pattern, Handler: handler, length: n, isDir: isDir})
	// append pattern to quick lookup map
	m.patterns[pattern] = 0
	//sort entries
	sort.Sort(m.entries)
}
