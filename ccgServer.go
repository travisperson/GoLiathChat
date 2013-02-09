/************************

Go Command Chat
-Jeromy Johnson, Travis Lane
A command line chat system that 
will make it easy to set up a 
quick secure chat room for any 
number of people

************************/

package main

import (
	"container/list"
	"crypto/rand"
	"crypto/tls"
	//"crypto/x509"
	"log"
	"net"
)

type Server struct {
	connections *list.List
	messages	*list.List
	listener	net.Listener
	com		chan Packet 
	parse  chan Packet
}

func StartServer() *Server {
	s := Server{}
	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	service := "127.0.0.1:10234"
	listener, err := tls.Listen("tcp", service, &config)
	s.listener = listener
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	s.connections = list.New()
	if err != nil {
		panic(err)
	}
	s.com = make(chan Packet)   //Channel for incoming messages
	s.parse = make(chan Packet) //Channel for parsed messages to be sent
	return &s
}

func (s *Server) Listen() {
	log.Print("server: listening")
	go s.MessageWriter()
	go s.MessageHandler()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		defer conn.Close()
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		_, ok := conn.(*tls.Conn) //Type assertion
		if ok {
			log.Print("ok=true")
			/*state := tlscon.ConnectionState()
			for _, v := range state.PeerCertificates {
				log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
			}*/
		}
		s.connections.PushBack(conn)
		u := UserWithConn(conn)
		go u.Handle(s.com) //Asynchronously listen to the connection
	}
}

//Receives packets parsed from incoming connections and 
//Processes them, then sends them to be relayed
func (s *Server) MessageHandler() {
	messages := *list.New()
	for {
		p := <-s.com
		//ts := time.Unix(int64(p.timestamp), 0)
		messages.PushFront(p)
		s.parse <- p
	}
}

//Receives and parses packets and then sends them to each connection in the list
//This is where any information requested is given
func (s *Server) MessageWriter() {
	for {
		p := <-s.parse

		//for now, just write the packets back.
		for i := s.connections.Front(); i != nil; i = i.Next() {
			_, err := i.Value.(net.Conn).Write(p.getBytes())
			if err != nil {
			}
		}
	}
}

func main() {
	s := StartServer()
	s.Listen()
}
