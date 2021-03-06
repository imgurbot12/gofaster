package gofaster

import (
	"bufio"
	"log"
	"net"
	"net/textproto"
	"sync"
	"time"
)

/***Variables***/

//TODO: need to reduce the number of allocations, any := or var at all is an allocation
//TODO: allow for http/1.0 vs http/1.1 response
//TODO: might want to re-modle connection handling to allow for keep-alive

//sockPool : object to maintain and create socket threads to handle inbound connections
type sockPool struct {
	listenerFunc func(string, string) (net.Listener, error)
	waitGroup    sync.WaitGroup // isRunning: used to stop listener threads after they have been spawned
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
	// open listener using config
	ln, err := s.listenerFunc(network, address)
	if err != nil {
		log.Fatalf("- Socket FAILED TO INIT! Error: %s\n", err)
	}
	defer ln.Close()
	defer s.waitGroup.Done()

	// make preallocated variables
	var (
		//re-used connection object
		conn net.Conn
		//deadline for rw on socket (adds connection timeout)
		deadline = 5 * time.Second
		//buffer objects reused for reading requests
		sBuffer = &textproto.Reader{}
		//request object reused for every request
		request = Request{
			sBuffer: sBuffer,
			bbuf: make([]byte, 1024),
		}
		//response object reused for every request
		response = Response{
			statusCode: 200,
		}
	)
	// pass connection to queue
	for {
		// accept connection and set deadline
		conn, err = ln.Accept()
		conn.SetDeadline(AproxTimeNow().Add(deadline))
		// handle connection
		if err == nil {
			// make buffers
			request.conn = conn
			request.sBuffer.R = bufio.NewReader(conn)
			// parse request
			err = request.parseRequest()
			// send request to appropriate handlers
			if err != nil {
				errorHandler(400, &request, &response)
			} else {
				handler(&request, &response)
			}
			// write response
			response.Make(conn)
			// reset variables
			response = Response{
				statusCode: 200,
			}
			request.Query = nil
			request.Form = nil
			request.Cookies = nil
		}
		//close connection
		conn.Close()
	}
}

//(*sockPool).Stop : halts all current socket workers
func (s *sockPool) Wait() {
	s.waitGroup.Wait()
}
