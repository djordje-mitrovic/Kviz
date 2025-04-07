package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	//"os"
	"strings"
	//"time"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
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
	db, err := sql.Open("sqlite3", "quiz.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	s.getQuestions(db)
	
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

func (s *Server) getQuestions(db *sql.DB) {
	rows, err := db.Query("SELECT pitanja, odgovor1, odgovor2, odgovor3, odgovor4, tacan_odgovor FROM questions ORDER BY RANDOM() LIMIT 10;")
	if err != nil {
		log.Fatal(err)
	}

	//pitanje1, a, b, c, d, 1
	//pitanje2, a, b, c, d, 2

	for rows.Next() {
		var pitanje string
		var opcija1, opcija2, opcija3, opcija4 string
		var odgovor int

		// Skeniraj podatke iz trenutnog reda u promenljive
		if err := rows.Scan(&pitanje, &opcija1, &opcija2, &opcija3, &opcija4, &odgovor); err != nil {
			log.Fatal(err)	
		}
		
		s.questions = append(s.questions, pitanje)
		s.options = append(s.options, []string{opcija1, opcija2, opcija3, opcija4})
		s.answers = append(s.answers, fmt.Sprintf("%d", odgovor))
	}
	if err := rows.Err(); err != nil {
		log.Fatal("Greska pri obradi podataka iz baze.")
	}
	fmt.Println("Pitanja su uspesno ucitana.")

}


// CTRL + K + C
// func (s *Server) getMessagesFromUser() {
// 	//reader := bufio.NewReader(os.Stdin)

// 	//fmt.Println("Unesi 10 pitanja, 4 opcije po pitanju i tacan odgovor:")

// 	// Unos pitanja
// 	for i := 0; i < 3; i++ {
// 		fmt.Printf("Unesi pitanje %d: ", i+1)
// 		msg, _ := reader.ReadString('\n')
// 		s.questions = append(s.questions, msg)

// 		// Unos ponuđenih odgovora
// 		options := make([]string, 4)
// 		for j := 0; j < 4; j++ {
// 			fmt.Printf("Unesi opciju %d za pitanje %d: ", j+1, i+1)
// 			option, _ := reader.ReadString('\n')
// 			options[j] = option
// 		}
// 		s.options = append(s.options, options)

// 		// Unos tačnog odgovora
// 		fmt.Printf("Unesi broj tacnog odgovora (1-4) za pitanje %d: ", i+1)
// 		var correctAnswer int
// 		_, err := fmt.Scanf("%d\n", &correctAnswer)
// 		if err != nil {
// 			fmt.Println("Pogresan ulaz")
// 		}
// 		// Dodajemo odgovarajući ponuđeni odgovor
// 		s.answers = append(s.answers, fmt.Sprintf("%d", correctAnswer))
// 	}
// }

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	poeni := 0

	// Interaktivno slanje pitanja i čekanje odgovora
	for i := 0; i < len(s.questions); i++ {
		// Šaljemo pitanje klijentu
		conn.Write([]byte(s.questions[i] + "\n"))
		//time.Sleep(1 * time.Second)

		// Šaljemo ponuđene odgovore
		for _, option := range s.options[i] {
			conn.Write([]byte(fmt.Sprintf("%s\n", option)))
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

		//fmt.Printf("%s se poredi sa %s\n", strings.TrimSpace(response), s.answers[i])
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