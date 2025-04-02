package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
	"math/rand"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var broj_pitanja = 2

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
	s.getMessagesFromDB()

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

func (s *Server) getMessagesFromDB() {
	db, err := sql.Open("sqlite3", "quiz.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//Mapa za proveravanje vec postojacih pitanja
	var askedQuestions = make(map[int]bool)
	
	
	// Unos pitanja
	for {
		var id int
		rand.Seed(time.Now().UnixNano())
		err = db.QueryRow("SELECT id FROM questions ORDER BY RANDOM() LIMIT 1").Scan(&id)
		
		fmt.Println("Id: ", id)
		
		if !askedQuestions[id] {
			askedQuestions[id] = true
			
			if err != nil {
				log.Fatal(err)
			}
			var pitanje, odg1, odg2, odg3, odg4, tacan_odg string
		
			err = db.QueryRow("SELECT pitanja, odgovor1, odgovor2, odgovor3, odgovor4, tacan_odgovor FROM questions WHERE id = ?", id).
			Scan(&pitanje, &odg1, &odg2, &odg3, &odg4, &tacan_odg)
		
			//fmt.Println(pitanje, " ", odg1, " ", odg2, " ", odg3, " ", odg4, " ", tacan_odg)
		
			if err != nil {
				log.Fatal(err)
			}
		
			s.questions = append(s.questions, pitanje)

			// Unos ponuđenih odgovora
			options := make([]string, 4)
		
			options[0] = odg1
			options[1] = odg2
			options[2] = odg3
			options[3] = odg4
			
			s.options = append(s.options, options)

			// Dodajemo odgovarajući ponuđeni odgovor
			s.answers = append(s.answers, tacan_odg)
		
		}
		
		if len(askedQuestions) >= broj_pitanja {
				fmt.Println("Broj postavljenih pitanja: ", len(askedQuestions))
				break
			}
		
	}
	
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	poeni := 0

	// Interaktivno slanje pitanja i čekanje odgovora
	for i := 0; i < broj_pitanja; i++ {
		// Šaljemo pitanje klijentu
		conn.Write([]byte(s.questions[i] + "\n"))
		//time.Sleep(1 * time.Second)

		// Šaljemo ponuđene odgovore
		for j, option := range s.options[i] {
			conn.Write([]byte(fmt.Sprintf("%d) %s\n", j+1, option)))
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
			poeni+=1
		}

		// Šaljemo rezultat klijentu
		_, err = conn.Write([]byte(tacno))
		if err != nil {
			fmt.Println("Greska prilikom slanja odgovora klijentu:", err)
			return
		}
	}

	conn.Write([]byte(fmt.Sprintf("%d\n", poeni)))

}

func main() {
	server := NewServer(":8082")
	log.Fatal(server.Start())
}
