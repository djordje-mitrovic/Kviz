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
	code        string			//jedinstveni broj sobe
	clients     []net.Conn			//lista klijenata
	questions   []string			//pitanja iz baze
	answers     []string			//tacni odgovori
	options     [][]string			//opcije koje su na raspolaganju
}

type Server struct {
	listenAddr string			//adresa na kojoj server slusa
	ln         net.Listener			//komponenta koja prihvata nove mrezne konekcie
	rooms      map[string]*Room		//mapa stringova u u pokazivace na Room strukturu za pracenje soba
	mu         sync.Mutex			//katanac za sinhronizaciju
}

func NewServer(listenAddr string) *Server {		//konstruktor za server
	return &Server{
		listenAddr: listenAddr,
		rooms:      make(map[string]*Room),
	}
}

func (s *Server) generateRoomCode() string {		//funkcija koja generise random kod sobe
	rand.Seed(time.Now().UnixNano())
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	code := make([]rune, 5)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}

func (s *Server) Start() error {
	db, err := sql.Open("sqlite3", "quiz.db")	//otvaramo bazu sa pitanjima
	if err != nil {					//provera greske
		log.Fatal(err)
	}
	defer db.Close()				//zatvaranje baze nakon sto zavrsimo sa njom

	s.ln, err = net.Listen("tcp", s.listenAddr)	//pokrecemo server koji osluskuje dolazne mrezne konekcije
	if err != nil {					//provera greske
		return err
	}
	defer s.ln.Close()				//gasenje servera
	fmt.Println("Server je pokrenut na", s.listenAddr)

	for {
		conn, err := s.ln.Accept()		//prihvatanje novih mreznih konekcija
		if err != nil {				//provera greske
			log.Println("Greška pri prihvatanju konekcije:", err)
			continue
		}
		go s.handleConnection(conn, db)		//pokretanje rada sa klijentom u gorutinu
	}
}

func (s *Server) handleConnection(conn net.Conn, db *sql.DB) {	//funkcija za rad sa klijentom
	reader := bufio.NewReader(conn)			//stvara novi bafer koji omogucava efikasno citanje sa conn objekta koji je mrezni konekcija sa klijentom
	command, _ := reader.ReadString('\n')		//citaju se podaci sa objekta reader dok se ne naidje na novi red
	command = strings.TrimSpace(command)		//uklanjamo beline oko primljenog stringa

	if command == "CREATE_ROOM" {			//ukoliko je klijent izabrao formiranje nove sobe
		roomCode := s.generateRoomCode()	//generisemo jedinstveni kod sobe
		//kako je uzimanje pitanja kriticna sekcija, potrebno je da zakljucamo katanac
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
	} else if strings.HasPrefix(command, "JOIN_ROOM") {	//ukoliko klijent bira da se pridruzi sobi
		parts := strings.Split(command, " ")	//citamo kod koji unosi klijent
		if len(parts) != 2 {
			conn.Write([]byte("INVALID_COMMAND\n"))
			return
		}
		roomCode := parts[1]
		//kriticna sekcija jer moze dvaklijenta da pokusavaju da se nakace na istu sobu u isto vreme
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
			//Kad u sobu udju dva igraca, pocinje igra
			go s.startGame(room)
		}
		s.mu.Unlock()
		conn.Write([]byte("JOINED_ROOM\n"))
	} else {
		conn.Write([]byte("UNKNOWN_COMMAND\n"))
	}
}

//Citanje pitanja iz baze
func (s *Server) getQuestionsForRoom(db *sql.DB) ([]string, [][]string, []string) {
	var questions []string
	var options [][]string
	var answers []string

	//Ucitamo 10 random pitanja iz baze
	rows, err := db.Query("SELECT pitanja, odgovor1, odgovor2, odgovor3, odgovor4, tacan_odgovor FROM questions ORDER BY RANDOM() LIMIT 10;")
	if err != nil {
		log.Fatal(err)
	}
	//Dodajemo pitanja, potencijalne odgovore, i tacne odgovore u nizove
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

	//Vracamo niz pitanja, potencijalnih odgovora i tacnih odgovora
	return questions, options, answers
}

//pocetak igre nakon sto su oba klijenta usla u sobu
func (s *Server) startGame(room *Room) {
	p1 := room.clients[0] //prvi klijent
	p2 := room.clients[1] //drugi klijent

	r1 := bufio.NewReader(p1)
	r2 := bufio.NewReader(p2)
	s1, s2 := 0, 0 //inicijalizacija poena oba igraca
	//izlistavanje pitanja
	for i, q := range room.questions {
		p1.Write([]byte(q + "\n")) //Saljemo pitanje prvogm igracu
		p2.Write([]byte(q + "\n")) //Saljemo pitanje drugom igracu
		for _, opt := range room.options[i] { // Saljemo 4 ponudjena odgovora i jednom i drugom
			p1.Write([]byte(opt + "\n")) 
			p2.Write([]byte(opt + "\n"))
		}
		//Citamo odgovore od oba igraca
		resp1, _ := r1.ReadString('\n')
		resp2, _ := r2.ReadString('\n')

		resp1 = strings.TrimSpace(resp1)
		resp2 = strings.TrimSpace(resp2)

		//Ako je odgovor koji je dao igrac tacan odgovor, azuriramo poene
		if resp1 == room.answers[i] {
    			s1++
		}
		if resp2 == room.answers[i] {
    			s2++
		}
		//Saljemo klijentima tacne odgovore
		p1.Write([]byte("TACAN_ODGOVOR " + room.answers[i] + "\n"))
		p2.Write([]byte("TACAN_ODGOVOR " + room.answers[i] + "\n"))
	}

	//Uporedjivanje rezultata
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
