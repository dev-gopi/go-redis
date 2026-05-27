package network

import (
	"io"
	"log"
	"net"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/router"
)

type Server struct {
	address string
}

func NewServer(addr string) *Server {
	return &Server{
		address: addr,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	log.Printf("TCP server started on %s", s.address)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Accept failed: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	cl := client.NewClient(conn)

	client.Manager.Add(cl)

	defer client.Manager.Remove(cl.ID)

	log.Printf(
		"Client connected: %s (%s)",
		cl.ID,
		conn.RemoteAddr(),
	)

	for {

		cmd, err := protocol.ParseRESP(cl.Reader)
		if err != nil {

			if err == io.EOF {
				log.Printf(
					"Client disconnected: %s",
					cl.ID,
				)
				return
			}

			log.Printf(
				"Connection failed [%s]: %v",
				cl.ID,
				err,
			)

			conn.Write([]byte(
				protocol.Error("invalid request"),
			))

			return
		}

		log.Printf(
			"Client %s command: %+v",
			cl.ID,
			cmd,
		)

		response := router.Handle(cl, cmd)

		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Printf(
				"Write failed [%s]: %v",
				cl.ID,
				err,
			)
			return
		}
	}
}

// func (s *Server) handleConnection(conn net.Conn) {
// 	defer conn.Close()

// 	log.Printf("Client connected: %s", conn.RemoteAddr())

// 	reader := bufio.NewReader(conn)

// 	for {
// 		cmd, err := protocol.ParseRESP(reader)
// 		if err != nil {

// 			if err == io.EOF {
// 				log.Printf("Client disconnected: %s", conn.RemoteAddr())
// 				return
// 			}

// 			log.Printf(
// 				"Connection failed [%s]: %v",
// 				conn.RemoteAddr(),
// 				err,
// 			)

// 			conn.Write([]byte(
// 				protocol.Error("invalid request"),
// 			))

// 			return
// 		}

// 		log.Printf(
// 			"Command from %s: %+v",
// 			conn.RemoteAddr(),
// 			cmd,
// 		)

// 		response := router.Handle(cmd)

// 		_, err = conn.Write([]byte(response))
// 		if err != nil {
// 			log.Printf(
// 				"Write failed [%s]: %v",
// 				conn.RemoteAddr(),
// 				err,
// 			)
// 			return
// 		}
// 	}
// }
