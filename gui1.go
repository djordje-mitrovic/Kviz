package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Kviz")

	conn, err := net.Dial("tcp", "localhost:8082")
	if err != nil {
		fmt.Println("Konekcija neuspela:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	mainContent := container.NewVBox()
	w.SetContent(mainContent)

	var loadNextQuestion func()

	loadNextQuestion = func() {
		// Čitanje pitanja
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Greška prilikom primanja pitanja:", err)
			return
		}
		question := strings.TrimSpace(msg)

		// Ako dobijemo broj kao pitanje (kraj igre), tretiraj kao broj poena
		if isNumber(question) {
			labelResult := widget.NewLabel(fmt.Sprintf("Osvojili ste %s poena!", question))
			cLabelResult := container.NewCenter(labelResult)

			mainContent.Objects = []fyne.CanvasObject{
				layout.NewSpacer(),
				cLabelResult,
				layout.NewSpacer(),
			}
			mainContent.Refresh()
			return
		}

		labelQuestion := widget.NewLabel(question)
		cLabelQuestion := container.NewCenter(labelQuestion)

		// Čitanje opcija
		var option [4]string
		buttons := make([]*widget.Button, 4)
		for j := 0; j < 4; j++ {
			option[j], err = reader.ReadString('\n')
			if err != nil {
				fmt.Println("Greška prilikom primanja opcije sa servera:", err)
				return
			}
			option[j] = strings.TrimSpace(option[j])
		}

		labelTacno := widget.NewLabel("")
		cLabelTacno := container.NewCenter(labelTacno)

		sendAnswer := func(answer int) {
			for _, btn := range buttons {
				btn.Disable()
			}

			_, err = conn.Write([]byte(fmt.Sprintf("%d\n", answer)))
			if err != nil {
				fmt.Println("Greška prilikom slanja odgovora:", err)
				return
			}

			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Greška prilikom čitanja odgovora sa servera:", err)
				return
			}

			labelTacno.SetText(strings.TrimSpace(message))
			mainContent.Refresh()

			time.AfterFunc(1*time.Second, func() {
				loadNextQuestion()
			})
		}

		// Dugmići za odgovore
		buttons[0] = widget.NewButton(option[0], func() { sendAnswer(1) })
		buttons[1] = widget.NewButton(option[1], func() { sendAnswer(2) })
		buttons[2] = widget.NewButton(option[2], func() { sendAnswer(3) })
		buttons[3] = widget.NewButton(option[3], func() { sendAnswer(4) })

		row1 := container.NewGridWithColumns(2, buttons[0], buttons[1])
		row2 := container.NewGridWithColumns(2, buttons[2], buttons[3])

		mainContent.Objects = []fyne.CanvasObject{
			layout.NewSpacer(),
			cLabelQuestion,
			layout.NewSpacer(),
			row1,
			row2,
			layout.NewSpacer(),
			cLabelTacno,
		}
		mainContent.Refresh()
	}

	loadNextQuestion()

	w.Resize(fyne.NewSize(600, 450))
	w.ShowAndRun()
}

// isNumber proverava da li je string broj (da razlikujemo poene od pitanja)
func isNumber(s string) bool {
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return len(s) > 0
}
