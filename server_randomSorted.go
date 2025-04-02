package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"

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
	// Povezivanje sa bazom podataka
	db, err := sql.Open("sqlite3", "quiz.db")
	if err != nil {
		log.Fatal("Greška pri povezivanju sa bazom:", err)
	}
	defer db.Close()

	// Učitavanje nasumičnih pitanja iz baze
	s.loadRandomQuestionsFromDB(db)

	// Startovanje servera
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

func (s *Server) loadRandomQuestionsFromDB(db *sql.DB) {
	// SQL upit za učitavanje nasumičnih pitanja iz baze
	query := `
	SELECT pitanja, odgovor1, odgovor2, odgovor3, odgovor4, tacan_odgovor
	FROM questions
	ORDER BY RANDOM()
	LIMIT 3;` // Ovde određujemo da želimo 3 nasumična pitanja

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("Greška pri učitavanju pitanja iz baze:", err)
	}
	defer rows.Close()

	// Čitanje pitanja iz rezultata upita
	for rows.Next() {
		var pitanje, odgovor1, odgovor2, odgovor3, odgovor4 string
		var tacanOdgovor int
		if err := rows.Scan(&pitanje, &odgovor1, &odgovor2, &odgovor3, &odgovor4, &tacanOdgovor); err != nil {
			log.Fatal("Greška pri čitanju pitanja:", err)
		}

		// Dodavanje pitanja u odgovarajuće nizove
		s.questions = append(s.questions, pitanje)
		s.options = append(s.options, []string{odgovor1, odgovor2, odgovor3, odgovor4})
		s.answers = append(s.answers, fmt.Sprintf("%d", tacanOdgovor))
	}
	if err := rows.Err(); err != nil {
		log.Fatal("Greška pri obradi rezultata upita:", err)
	}

	fmt.Println("Pitanja su uspešno učitana iz baze.")
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	poeni := 0

	// Interaktivno slanje pitanja i čekanje odgovora
	for i := 0; i < len(s.questions); i++ {
		// Šaljemo pitanje klijentu
		conn.Write([]byte(s.questions[i] + "\n"))

		// Šaljemo ponuđene odgovore
		for j, option := range s.options[i] {
			conn.Write([]byte(fmt.Sprintf("%d) %s\n", j+1, option)))
		}

		// Čekamo odgovor od klijenta
		reader := bufio.NewReader(conn)
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Greška prilikom čitanja odgovora od klijenta:", err)
			return
		}

		// Proveravamo tačnost odgovora
		tacno := "Odgovor je netačan\n"
		if strings.TrimSpace(response) == s.answers[i] {
			tacno = "Odgovor je tačan\n"
			poeni++
		}

		// Šaljemo rezultat klijentu
		_, err = conn.Write([]byte(tacno))
		if err != nil {
			fmt.Println("Greška prilikom slanja odgovora klijentu:", err)
			return
		}
	}

	conn.Write([]byte(fmt.Sprintf("Ostvarili ste %d poena.\n", poeni)))
}

func main() {
	server := NewServer(":8082")
	log.Fatal(server.Start())
}
