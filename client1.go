package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

var broj_pitanja = 2

func main() {
	conn, err := net.Dial("tcp", "localhost:8082")
	if err != nil {
		fmt.Println("Konekcija neuspela:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Interaktivno odgovaranje na pitanja
	for i := 0; i < broj_pitanja; i++ {
		// Prima pitanje od servera
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Greska prilikom primanja poruke sa servera:", err)
			return
		}
		
		// Ispisuje pitanje
		fmt.Printf("Pitanje %d: %s", i+1, msg)

		// Prima ponuđene odgovore
		for j := 0; j < 4; j++ {
			option, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Greska prilikom primanja opcije sa servera:", err)
				return
			}
			fmt.Print(option)
		}

		// Unos odgovora
		fmt.Print("Izaberite jedan od ponuđenih odgovora (1-4): ")
		var selectedOption int
		_, err = fmt.Scanf("%d\n", &selectedOption)
		if err != nil || selectedOption < 1 || selectedOption > 4 {
			fmt.Println("Pogresan unos")
			return
		}

		// Šalje odgovor serveru
		_, err = conn.Write([]byte(fmt.Sprintf("%d\n", selectedOption)))
		if err != nil {
			fmt.Println("Greksa prilikom slanja odgovora na server:", err)
			return
		}

		// Čitanje odgovora od servera (da li je tačan odgovor)
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Greska prilikom citanja poruke sa servera:", err)
			return
		}

		// Ispisivanje rezultata
		fmt.Println("Odgovor od servera:", strings.TrimSpace(message))
	}

	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Greska prilikom citanja poruke sa servera:", err)
		return
	}

	fmt.Println("Odgovor od servera: Osvojili ste", strings.TrimSpace(message), "poena")

}
