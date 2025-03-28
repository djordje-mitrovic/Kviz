package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8082")
	if err != nil {
		fmt.Println("Connection failed:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	receivedMessages := make([]string, 5) // Definišemo niz za čuvanje poruka

	for i := 0; i < 5; i++ {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error receiving message from server:", err)
			return
		}
		receivedMessages[i] = msg // Čuvamo poruku u nizu
	}

	// Ispisujemo sve primljene poruke nakon što su sve sačuvane
	i := 0
	for _, msg := range receivedMessages {
		if(i == 0){
			fmt.Print("question", msg)
		} else {
			fmt.Print(i, ") ", msg)
		}
		i++
	}


	fmt.Print("Select a message to return: ")
	userInput := bufio.NewReader(os.Stdin)
	selectedMessage, _ := userInput.ReadString('\n')

	// Proveravamo da li se unesena poruka nalazi u nizu receivedMessages
	messageFound := false
	for _, msg := range receivedMessages {
		if selectedMessage == msg {
			messageFound = true
			break
		}
	}

	// Ako je poruka pronađena u nizu, šaljemo je serveru, inače prijavljujemo grešku
	if messageFound {
		_, err := conn.Write([]byte(selectedMessage))
		if err != nil {
			fmt.Println("Error sending response to server:", err)
		}
	} else {
		fmt.Println("Error: The selected message is not in the list.")
	}
}
