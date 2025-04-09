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
		return color.NRGBA{R: 0, G: 0, B: 128, A: 109}
	case "text":
		return color.NRGBA{R: 20, G: 20, B: 20, A: 255}
	default:
		return theme.DefaultTheme().Color(n, v)
	}
}
func (m myTheme) Font(style fyne.TextStyle) fyne.Resource  { return theme.DefaultTheme().Font(style) }
func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource { return theme.DefaultTheme().Icon(name) }
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
	a := app.New()
	a.Settings().SetTheme(&myTheme{})

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

	title := widget.NewLabel("🎉 Dobrodošao u Kviz! 🎉")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	btnCreate := widget.NewButton("Kreiraj sobu", func() {
		var err error
		conn, err = net.Dial("tcp", "localhost:8082")
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		conn.Write([]byte("CREATE_ROOM\n"))
		reader := bufio.NewReader(conn)
		roomCode, _ := reader.ReadString('\n')
		roomCode = strings.TrimSpace(strings.TrimPrefix(roomCode, "ROOM_CODE "))

		waitLabel := widget.NewLabel("Čeka se drugi igrač")
		startAnimation(waitLabel)

		info := widget.NewLabel("Tvoj kod sobe: " + roomCode)
		info.Alignment = fyne.TextAlignCenter

		content.Objects = []fyne.CanvasObject{
			layout.NewSpacer(),
			info,
			waitLabel,
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

		var err error
		conn, err = net.Dial("tcp", "localhost:8082")
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		conn.Write([]byte("JOIN_ROOM " + code + "\n"))
		reader := bufio.NewReader(conn)
		resp, _ := reader.ReadString('\n')
		resp = strings.TrimSpace(resp)
		if resp == "ROOM_NOT_FOUND" || resp == "ROOM_FULL" {
			dialog.ShowInformation("Greška", resp, w)
			return
		}

		waitLabel := widget.NewLabel("Čeka se početak")
		startAnimation(waitLabel)

		content.Objects = []fyne.CanvasObject{
			layout.NewSpacer(),
			waitLabel,
			layout.NewSpacer(),
		}
		content.Refresh()
		waitForStart(conn, reader, content, w)
	})

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

func waitForStart(conn net.Conn, reader *bufio.Reader, content *fyne.Container, w fyne.Window) {
	go func() {
		loadQuiz(conn, reader, content, w)
	}()
}

func loadQuiz(conn net.Conn, reader *bufio.Reader, content *fyne.Container, w fyne.Window) {
	var loadNextQuestion func()

	loadNextQuestion = func() {
		question, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		question = strings.TrimSpace(question)

		// Kraj kviza
		if strings.HasPrefix(question, "Pobedio") || strings.HasPrefix(question, "Izgubio") || strings.HasPrefix(question, "Nereseno") {
			labelEnd := widget.NewLabel(question)
			labelEnd.Alignment = fyne.TextAlignCenter
			labelEnd.TextStyle = fyne.TextStyle{Bold: true}

			btnBack := widget.NewButton("Vrati se u meni", func() {
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

		labelQuestion := widget.NewLabel(question)
		labelQuestion.Alignment = fyne.TextAlignCenter
		labelQuestion.TextStyle = fyne.TextStyle{Bold: true}

		options := make([]string, 4)
		for i := 0; i < 4; i++ {
			options[i], _ = reader.ReadString('\n')
			options[i] = strings.TrimSpace(options[i])
		}

		resultLabel := widget.NewLabel("")
		buttons := []*widget.Button{}
		for i, opt := range options {
			idx := i + 1
			optCopy := opt
			btn := widget.NewButton(optCopy, func() {
				conn.Write([]byte(fmt.Sprintf("%d\n", idx)))
				feedback, _ := reader.ReadString('\n')
				resultLabel.SetText(strings.TrimSpace(feedback))
				for _, b := range buttons {
					b.Disable()
				}
				content.Refresh()
				time.AfterFunc(2*time.Second, func() {
					loadNextQuestion()
				})
			})
			btn.Resize(fyne.NewSize(0, 60))
			buttons = append(buttons, btn)
		}

		row1 := container.NewGridWithColumns(2, buttons[0], buttons[1])
		row2 := container.NewGridWithColumns(2, buttons[2], buttons[3])

		questionContainer := container.NewVBox(
			layout.NewSpacer(),
			container.NewCenter(labelQuestion),
			layout.NewSpacer(),
			row1,
			row2,
			layout.NewSpacer(),
			resultLabel,
		)

		content.Objects = []fyne.CanvasObject{questionContainer}
		content.Refresh()
	}

	loadNextQuestion()
}

