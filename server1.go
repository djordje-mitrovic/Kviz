package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	//"time"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	questions  []string
	answers    []string
	options    [][]string
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
	fmt.Println("Server je na adresi", s.listenAddr)

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Prihvati gresku:", err)
			continue
		}

		go s.handleClient(conn)
	}
}

func (s *Server) getMessagesFromUser() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Unesi 10 pitanja, 4 opcije po pitanju, i tacan odgovor:")

	// Unos pitanja
	for i := 0; i < 10; i++ {
		fmt.Printf("Unesi pitanje %d: ", i+1)
		msg, _ := reader.ReadString('\n')
		s.questions = append(s.questions, msg)

		// Unos ponuđenih odgovora
		options := make([]string, 4)
		for j := 0; j < 4; j++ {
			fmt.Printf("Unesi opciju %d za pitaje %d: ", j+1, i+1)
			option, _ := reader.ReadString('\n')
			options[j] = option
		}
		s.options = append(s.options, options)

		// Unos tačnog odgovora
		fmt.Printf("Unesi broj tacnog odgovora (1-4) za pitanje %d: ", i+1)
		var correctAnswer int
		_, err := fmt.Scanf("%d\n", &correctAnswer)
		if err != nil {
			fmt.Println("Pogresan ulaz")
		}
		// Dodajemo odgovarajući ponuđeni odgovor
		s.answers = append(s.answers, fmt.Sprintf("%d", correctAnswer))
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	// Interaktivno slanje pitanja i čekanje odgovora
	for i := 0; i < 2; i++ {
		// Šaljemo pitanje klijentu
		conn.Write([]byte(s.questions[i]))
		//time.Sleep(1 * time.Second)

		// Šaljemo ponuđene odgovore
		for j, option := range s.options[i] {
			conn.Write([]byte(fmt.Sprintf("%d) %s", j+1, option)))
		}

		// Čekamo odgovor od klijenta
		reader := bufio.NewReader(conn)
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Greska prilikom citanja odgovora klijenta:", err)
			return
		}

		// Proveravamo tačnost odgovora
		tacno := "Odgovor je netačan\n"
		//num, _ := strconv.Atoi(response)

		fmt.Printf("%s se poredi sa %s\n", strings.TrimSpace(response), s.answers[i])
		if strings.TrimSpace(response) == s.answers[i] {
			tacno = "Odgovor je tačan\n"
		}

		// Šaljemo rezultat klijentu
		_, err = conn.Write([]byte(tacno))
		if err != nil {
			fmt.Println("Greska prilikom slanja odgovora klijentu:", err)
			return
		}
	}
}

func main() {
	server := NewServer(":8082")
	log.Fatal(server.Start())
}
