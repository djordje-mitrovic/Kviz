package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

type Room struct {
	code        string
	clients     []net.Conn
	questions   []string
	answers     []string
	options     [][]string
}

type Server struct {
	listenAddr string
	ln         net.Listener
	rooms      map[string]*Room
	mu         sync.Mutex
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		rooms:      make(map[string]*Room),
	}
}

func (s *Server) generateRoomCode() string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	code := make([]rune, 5)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}

func (s *Server) Start() error {
	db, err := sql.Open("sqlite3", "quiz.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s.ln, err = net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer s.ln.Close()
	fmt.Println("Server je pokrenut na", s.listenAddr)

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Println("Greška pri prihvatanju konekcije:", err)
			continue
		}
		go s.handleConnection(conn, db)
	}
}

func (s *Server) handleConnection(conn net.Conn, db *sql.DB) {
	reader := bufio.NewReader(conn)
	command, _ := reader.ReadString('\n')
	command = strings.TrimSpace(command)

	if command == "CREATE_ROOM" {
		roomCode := s.generateRoomCode()
		s.mu.Lock()
		// Kreiranje sobe sa pitanjima iz baze
		questions, options, answers := s.getQuestionsForRoom(db)
		s.rooms[roomCode] = &Room{
			code:      roomCode,
			clients:   []net.Conn{conn},
			questions: questions,
			options:   options,
			answers:   answers,
		}
		s.mu.Unlock()
		conn.Write([]byte("ROOM_CODE " + roomCode + "\n"))
	} else if strings.HasPrefix(command, "JOIN_ROOM") {
		parts := strings.Split(command, " ")
		if len(parts) != 2 {
			conn.Write([]byte("INVALID_COMMAND\n"))
			return
		}
		roomCode := parts[1]
		s.mu.Lock()
		room, exists := s.rooms[roomCode]
		if !exists {
			conn.Write([]byte("ROOM_NOT_FOUND\n"))
			s.mu.Unlock()
			return
		}
		if len(room.clients) >= 2 {
			conn.Write([]byte("ROOM_FULL\n"))
			s.mu.Unlock()
			return
		}
		room.clients = append(room.clients, conn)
		if len(room.clients) == 2 {
			go s.startGame(room)
		}
		s.mu.Unlock()
		conn.Write([]byte("JOINED_ROOM\n"))
	} else {
		conn.Write([]byte("UNKNOWN_COMMAND\n"))
	}
}

func (s *Server) getQuestionsForRoom(db *sql.DB) ([]string, [][]string, []string) {
	var questions []string
	var options [][]string
	var answers []string

	rows, err := db.Query("SELECT pitanja, odgovor1, odgovor2, odgovor3, odgovor4, tacan_odgovor FROM questions ORDER BY RANDOM() LIMIT 10;")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var pitanje, o1, o2, o3, o4 string
		var odgovor int
		if err := rows.Scan(&pitanje, &o1, &o2, &o3, &o4, &odgovor); err != nil {
			log.Fatal(err)
		}
		questions = append(questions, pitanje)
		options = append(options, []string{o1, o2, o3, o4})
		answers = append(answers, fmt.Sprintf("%d", odgovor))
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return questions, options, answers
}

func (s *Server) startGame(room *Room) {
	p1 := room.clients[0]
	p2 := room.clients[1]

	r1 := bufio.NewReader(p1)
	r2 := bufio.NewReader(p2)
	s1, s2 := 0, 0
	for i, q := range room.questions {
		p1.Write([]byte(q + "\n"))
		p2.Write([]byte(q + "\n"))
		for _, opt := range room.options[i] {
			p1.Write([]byte(opt + "\n"))
			p2.Write([]byte(opt + "\n"))
		}
		resp1, _ := r1.ReadString('\n')
		resp2, _ := r2.ReadString('\n')
		resp1 = strings.TrimSpace(resp1)
		resp2 = strings.TrimSpace(resp2)
		if resp1 == room.answers[i] {
    			s1++
		}
		if resp2 == room.answers[i] {
    			s2++
		}
		p1.Write([]byte("TACAN_ODGOVOR " + room.answers[i] + "\n"))
		p2.Write([]byte("TACAN_ODGOVOR " + room.answers[i] + "\n"))
	}

	if s1 > s2 {
		p1.Write([]byte(fmt.Sprintf("Pobedio si sa %d poena!\n", s1)))
		p2.Write([]byte(fmt.Sprintf("Izgubio si sa %d poena!\n", s2)))
	} else if s2 > s1 {
		p2.Write([]byte(fmt.Sprintf("Pobedio si sa %d poena!\n", s2)))
		p1.Write([]byte(fmt.Sprintf("Izgubio si sa %d poena!\n", s1)))
	} else {
		p1.Write([]byte(fmt.Sprintf("Nereseno! Osvojio si %d poena!\n", s1)))
		p2.Write([]byte(fmt.Sprintf("Nereseno! Osvojio si %d poena!\n", s2)))
	}
}



func main() {
	server := NewServer(":8082")
	log.Fatal(server.Start())
}
