package gofaster

import (
	"bufio"
	"log"
	"net"
	"net/textproto"
	"sync"
	"time"

	"github.com/valyala/tcplisten"
)

/***Variables***/

//TODO: allow for http/1.0 vs http/1.1 response
//TODO: extra requests are not handled by super fast sockets, they are handled via fasthttp.Server.serveConn

//sockPool : object to maintain and create socket threads to handle inbound connections
type sockPool struct {
	waitGroup  sync.WaitGroup // isRunning: used to stop listener threads after they have been spawned
	readerPool sync.Pool
}

/***Methods***/

//(*sockPool).Spawn : spawn `n` number of sockets with SO_REUSEADDR
func (s *sockPool) Spawn(network, address string, n int, handler func(*Request, *Response)) {
	for i := 0; i < n; i++ {
		go s.listen(network, address, handler)
		s.waitGroup.Add(1)
	}
}

//(*sockPool).listen : spawn socket and send inbound connections to queue
func (s *sockPool) listen(network, address string, handler func(*Request, *Response)) {
	//build config for tcp-listener
	cfg := &tcplisten.Config{
		ReusePort:   true,
		DeferAccept: true,
		FastOpen:    true,
	}
	// open listener
	ln, err := cfg.NewListener(network, address)
	if err != nil {
		log.Fatalf("- Socket FAILED TO INIT! Error: %s\n", err)
	}
	defer ln.Close()
	defer s.waitGroup.Done()

	// make preaollocated variables
	var (
		//re-used connection object
		conn net.Conn
		//deadline for rw on socket (adds connection timeout)
		deadline = 5 * time.Second
		//buffer objects reused for reading requests
		sBuffer = &textproto.Reader{}
		//requst object reused for every request
		request = Request{
			sBuffer: sBuffer,
		}
		//response object reused for every request
		response = Response{
			statusCode: 200,
		}
	)
	// pass connection to queue
	for {
		// accept conection and set deadline
		conn, err = ln.Accept()
		conn.SetDeadline(AproxTimeNow().Add(deadline))
		// handle connection
		if err == nil {
			// make buffers
			request.sBuffer.R = bufio.NewReader(conn)
			// parse request
			err = request.parseRequest()
			// send request to appropriate handlers
			if err != nil {
				errorHandler(400, &response)
			} else {
				handler(&request, &response)
			}
			// write response
			response.Make(conn)
			// reset variables
			response = Response{
				statusCode: 200,
			}
			// s.putReader(request.sBuffer.R)
		}
		//close connection
		conn.Close()
	}
}

//(*sockPool).getReader : return bufio.Reader instance for use in request handler
func (s *sockPool) getReader(c net.Conn) *bufio.Reader {
	v := s.readerPool.Get()
	if v == nil {
		return bufio.NewReader(c)
	}
	r := v.(*bufio.Reader)
	r.Reset(c)
	return r
}

//(*sockPool).putReader : return bufio.Reader instance to readerPool
func (s *sockPool) putReader(b *bufio.Reader) {
	s.readerPool.Put(b)
}

//(*sockPool).Stop : halts all current socket workers
func (s *sockPool) Wait() {
	s.waitGroup.Wait()
}
