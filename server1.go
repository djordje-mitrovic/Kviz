package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	messages   []string
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}

func (s *Server) Start() error {
	s.getMessagesFromUser()

	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln
	fmt.Println("Server is listening on", s.listenAddr)

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		go s.handleClient(conn)
	}
}

func (s *Server) getMessagesFromUser() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter 4 messages to send to the client:")
	
	fmt.Println("Enter the question:")
	msg, _ := reader.ReadString('\n')
	s.messages = append(s.messages, msg)
	
	for i := 0; i < 4; i++ {
		fmt.Printf(i+1, ": ")
		msg, _ := reader.ReadString('\n')
		s.messages = append(s.messages, msg)
	}
	fmt.Printf("tacan odgovor: ")
	msg, _ := reader.ReadString('\n')
	s.messages = append(s.messages, msg)
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	for _, msg := range s.messages {
		conn.Write([]byte(msg))
		time.Sleep(1 * time.Second)
	}

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading response from client:", err)
		return
	}

	fmt.Println("Client selected:", response)
	if(response == s.messages[5]) {
		fmt.Println("Odgovor je tacan")
	} else {
		fmt.Println("Odgovor je netacan")
	}
}

func main() {
	server := NewServer(":8082")
	log.Fatal(server.Start())
}
