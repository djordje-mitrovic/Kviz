package main

import (
	"bufio"
	"fmt"
	"image/color"
	"net"
	"strings"
	"time"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type myTheme struct{}

func (m myTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case "background":
		return color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	case "button":
		return color.NRGBA{R: 128, G: 128, B: 128, A: 189}
	case "text":
		return color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	default:
		return theme.DefaultTheme().Color(n, v)
	}
}
func (m myTheme) Font(style fyne.TextStyle) fyne.Resource      { return theme.DefaultTheme().Font(style) }
func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource   { return theme.DefaultTheme().Icon(name) }
func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case "text":
		return 18
	case "padding":
		return 12
	default:
		return theme.DefaultTheme().Size(name)
	}
}

func main() {
	//Pravljenje aplikacije
	a := app.New()
	a.Settings().SetTheme(&myTheme{})

	//Pravljenje prozora i namestanje dimenzija
	w := a.NewWindow("Kviz")
	w.Resize(fyne.NewSize(600, 400))

	mainContent := container.NewVBox()
	showMainMenu(w, mainContent)

	w.SetContent(mainContent)
	w.ShowAndRun()
}

func showMainMenu(w fyne.Window, content *fyne.Container) {
	entryRoomCode := widget.NewEntry()
	entryRoomCode.SetPlaceHolder("Unesi kod sobe")

	var conn net.Conn

	title := widget.NewLabel("Dobrodošao u Kviz!")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	//Dugme za kreiranje sobe, kad se kreira soba povezujemo se na server
	btnCreate := widget.NewButton("Kreiraj sobu", func() {
		var err error
		conn, err = net.Dial("tcp", "localhost:8082")
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		//Saljemo serveru poruku da je soba napravljena
		conn.Write([]byte("CREATE_ROOM\n"))
		reader := bufio.NewReader(conn)
		roomCode, _ := reader.ReadString('\n')
		//Server salje random izgenerisan kod sobe
		roomCode = strings.TrimSpace(strings.TrimPrefix(roomCode, "ROOM_CODE "))

		waitLabel := widget.NewLabel("Čeka se drugi igrač")
		startAnimation(waitLabel)

		info := widget.NewLabel("Tvoj kod sobe: " + roomCode)
		//Postavljamo labelu na kojoj ce pisati kod koji treba drugi igrac da unese
		info.Alignment = fyne.TextAlignCenter

		//Dugme koje nas vraca nazad na pocetnu stranu
		btnBack := widget.NewButton("Nazad", func() {
			conn.Close()
			showMainMenu(w, content)
		})


		content.Objects = []fyne.CanvasObject{
			layout.NewSpacer(),
			info,
			waitLabel,
			btnBack,
			layout.NewSpacer(),
		}
		content.Refresh()
		waitForStart(conn, reader, content, w)
	})

	
	btnJoin := widget.NewButton("Pridruži se sobi", func() {
		code := strings.TrimSpace(entryRoomCode.Text)
		if code == "" {
			dialog.ShowInformation("Greška", "Unesi kod sobe", w)
			return
		}

		//Povezivanje klijenta na server
		var err error
		conn, err = net.Dial("tcp", "localhost:8082")
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
//192.168.206.159
		//Saljemo serveru komandu sa kojom ulazimo u sobu sa prvim igracem koji je napravio sobu
		conn.Write([]byte("JOIN_ROOM " + code + "\n"))
		reader := bufio.NewReader(conn)
		resp, _ := reader.ReadString('\n')
		//Citamo odgovor od servera
		resp = strings.TrimSpace(resp)
		if resp == "ROOM_NOT_FOUND" || resp == "ROOM_FULL" {
			dialog.ShowInformation("Greška", resp, w)
			return
		}


		waitLabel := widget.NewLabel("Čeka se početak")
		startAnimation(waitLabel)

		//Dugme za povratak na pocetni meni
		btnBack := widget.NewButton("Nazad", func() {
			conn.Close()
			showMainMenu(w, content)
		})

		content.Objects = []fyne.CanvasObject{
			layout.NewSpacer(),
			waitLabel,
			btnBack,
			layout.NewSpacer(),
		}
		content.Refresh()
		waitForStart(conn, reader, content, w)
	})

	//Dodavanje svih komponenti
	menu := container.NewVBox(
		layout.NewSpacer(),
		title,
		layout.NewSpacer(),
		btnCreate,
		entryRoomCode,
		btnJoin,
		layout.NewSpacer(),
	)

	content.Objects = []fyne.CanvasObject{menu}
	content.Refresh()
}

func startAnimation(label *widget.Label) {
	//Funkcija koja simulira cekanje tako sto prikazuje . pa .. pa ... i tako u krug
	go func() {
		dots := []string{"", ".", "..", "..."}
		i := 0
		for {
			time.Sleep(500 * time.Millisecond)
			label.SetText("Čeka se protivnik" + dots[i%4])
			i++
		}
	}()
}

//Funkcija koja ceka da oba igraca udju u sobu i tad ucitava pitanja
func waitForStart(conn net.Conn, reader *bufio.Reader, content *fyne.Container, w fyne.Window) {
	go func() {
		loadQuiz(conn, reader, content, w)
	}()
}

func loadQuiz(conn net.Conn, reader *bufio.Reader, content *fyne.Container, w fyne.Window) {
	var loadNextQuestion func()

	loadNextQuestion = func() {
		//Ucitavamo pitanje koje je stiglo od servera
		question, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		question = strings.TrimSpace(question)

		if strings.HasPrefix(question, "Pobedio") || strings.HasPrefix(question, "Izgubio") || strings.HasPrefix(question, "Nereseno") {
			//U slucaju da nam je stigo string koji nije pitanje, prikazijuemo samo poruku koja nas obavestava o rezultatu i imamo opciju da se vratimo na pocetak
			labelEnd := widget.NewLabel(question)
			labelEnd.Alignment = fyne.TextAlignCenter
			labelEnd.TextStyle = fyne.TextStyle{Bold: true}

			//Dugme koje vraca na pocetni meni
			btnBack := widget.NewButton("Vrati se u pocetni meni", func() {
				conn.Close()
				showMainMenu(w, content)
			})

			content.Objects = []fyne.CanvasObject{
				layout.NewSpacer(),
				labelEnd,
				layout.NewSpacer(),
				btnBack,
			}
			content.Refresh()
			return
		}

		//Pitanje koje je stiglo se postavlja na labelu
		labelQuestion := widget.NewLabel(question)
		labelQuestion.Alignment = fyne.TextAlignCenter
		labelQuestion.TextStyle = fyne.TextStyle{Bold: true}

		//Ucitavaju se ponudjeni odgovori
		options := make([]string, 4)
		for i := 0; i < 4; i++ {
			options[i], _ = reader.ReadString('\n')
			options[i] = strings.TrimSpace(options[i])
		}

		buttons := make([]*widget.Button, 4)
		buttonBoxes := make([]*fyne.Container, 4)

		for i := 0; i < 4; i++ {
		    btnIndex := i
		    opt := options[i]
		    btn := widget.NewButton(opt, nil)
		    box := container.NewMax(canvas.NewRectangle(color.White), btn)

		    btn.OnTapped = func() {
			// Kada kliknemo, prvo zakljucamo sva dugmad
			for j := 0; j < 4; j++ {
			    btns, ok := buttonBoxes[j].Objects[1].(*widget.Button)
			    if ok {
				btns.Disable()
			    }
			}

			// Obojimo izabrano dugme plavom bojom odmah
			btnSize := btn.Size()
			bg := canvas.NewRectangle(color.RGBA{R: 0, G: 122, B: 255, A: 255}) // plava
			bg.SetMinSize(btnSize)

			label := canvas.NewText(btn.Text, color.White)
			label.Alignment = fyne.TextAlignCenter
			label.TextStyle = fyne.TextStyle{Bold: true}
			label.TextSize = 18

			newBox := container.NewMax(bg, container.NewCenter(label))
			newBox.Resize(btnSize)

			buttonBoxes[btnIndex].Objects = []fyne.CanvasObject{newBox}
			buttonBoxes[btnIndex].Refresh()

			content.Refresh()

			// Sada tek šaljemo serveru
			conn.Write([]byte(fmt.Sprintf("%d\n", btnIndex+1)))

			// Cekamo odgovor sa servera (TACAN_ODGOVOR x)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(response)

			var correct int
			if strings.HasPrefix(response, "TACAN_ODGOVOR") {
			    parts := strings.Split(response, " ")
			    if len(parts) == 2 {
				fmt.Sscanf(parts[1], "%d", &correct)
			    }
			}

			// Obojimo sve dugmice kako treba
			for j := 0; j < 4; j++ {
			    var bgColor color.Color

			    if j+1 == correct {
				bgColor = color.RGBA{R: 11, G: 156, B: 49, A: 255} // Zelena
			    }
			    if j == btnIndex && btnIndex+1 != correct {
				bgColor = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Crvena
			    }

			    btnSize := fyne.NewSize(300, 40)

			    bg := canvas.NewRectangle(bgColor)
			    bg.SetMinSize(btnSize)

			    label := canvas.NewText(options[j], color.White)
			    label.Alignment = fyne.TextAlignCenter
			    label.TextStyle = fyne.TextStyle{Bold: true}
			    label.TextSize = 18

			    newBox := container.NewMax(bg, container.NewCenter(label))
			    newBox.Resize(btnSize)

			    buttonBoxes[j].Objects = []fyne.CanvasObject{newBox}
			    buttonBoxes[j].Refresh()
			}

			content.Refresh()

			time.AfterFunc(2*time.Second, func() {
			    loadNextQuestion()
			})
		    }

		    buttons[i] = btn
		    buttonBoxes[i] = box
		}


		row1 := container.NewGridWithColumns(2, buttonBoxes[0], buttonBoxes[1])
		row2 := container.NewGridWithColumns(2, buttonBoxes[2], buttonBoxes[3])

		//Dodavanje svih komponenti u GUI
		questionContainer := container.NewVBox(
			layout.NewSpacer(),
			container.NewCenter(labelQuestion),
			layout.NewSpacer(),
			row1,
			row2,
			layout.NewSpacer(),
		)

		content.Objects = []fyne.CanvasObject{questionContainer}
		content.Refresh()
	}

	loadNextQuestion()
}
