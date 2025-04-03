

package main

import(
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)


func main () {

	db, err := sql.Open("sqlite3", "quiz.db")		
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	createTable(db)
}

func createTable(db *sql.DB) {

	query := `CREATE TABLE IF NOT EXISTS questions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		pitanja TEXT NOT NULL,
		odgovor1 TEXT NOT NULL,
		odgovor2 TEXT NOT NULL,
		odgovor3 TEXT NOT NULL,
		odgovor4 TEXT NOT NULL,
		tacan_odgovor INTEGER NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tabela je kreirana ili vec postoji.")

}