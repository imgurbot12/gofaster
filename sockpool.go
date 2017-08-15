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

//sockPool : object to maintain and create socket threads to handle inbound connections
type sockPool struct {
	waitGroup sync.WaitGroup // isRunning: used to stop listener threads after they have been spawned
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
		//buffer object reused for reading requests
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
		conn.SetDeadline(time.Now().Add(deadline))
		// handle connection
		if err == nil {
			// make buffers
			request.sBuffer.R = bufio.NewReader(conn)
			// parse request
			err = request.parseRequest()
			// send request to appropriate handlers
			if err != nil {
				errorHandler("Bad Request", &response)
			} else {
				handler(&request, &response)
			}
			// write response
			response.Make(conn)
			// reset variables
			response = Response{
				statusCode: 200,
			}
		}
		//close connection
		conn.Close()
	}
}

//(*sockPool).Stop : halts all current socket workers
func (s *sockPool) Wait() {
	s.waitGroup.Wait()
}