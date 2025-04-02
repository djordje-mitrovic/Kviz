package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

// Ubacuje pitanja u bazu
func InsertQuestions(db *sql.DB) {
	query := `
	INSERT INTO questions (id, pitanja, odgovor1, odgovor2, odgovor3, odgovor4, tacan_odgovor)
	VALUES 
	(1, 'Koji je glavni grad Francuske?', 'Madrid', 'London', 'Pariz', 'Berlin', 3),
	(2, 'Kad je bio kosovski boj?', '1389', '1375', '1412', '1403', 1),
	(3, 'Koji je hemijski simbol za vodonik?', 'O', 'H', 'C', 'N', 2);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Greška pri ubacivanju podataka:", err)
	}

	fmt.Println("Pitanja su uspešno dodata u bazu!")
}

//
func main() {
	db, err := sql.Open("sqlite3", "quiz.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	InsertQuestions(db)

}

