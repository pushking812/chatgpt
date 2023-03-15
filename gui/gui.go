package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/pushking812/chatgpt/gpt"
	cc "github.com/pushking812/chatgpt/gpt/chat-complete"
	tt "github.com/pushking812/chatgpt/gpt/speech-to-text"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("ChatGPT")

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		// Create a label with a message to display
		label := widget.NewLabel("Missing API KEY")

		// Create a button to close the window
		button := widget.NewButton("Close", func() {
			myApp.Quit()
		})

		// Create a layout and add the label and button to it
		layout := container.NewVBox(label, button)

		// Set the layout as the window's content
		myWindow.SetContent(layout)

		// Show the window
		myWindow.ShowAndRun()
		return
	}

	chatRequest:=gpt.NewRequestType(cc.GetAnswer, apiKey, "10s")
	sttRequest:=gpt.NewRequestType(tt.GetAnswer, apiKey, "10s")

	// UI elements
	questionEntry := widget.NewEntry()
	questionEntry.SetPlaceHolder("Type your question here")
	//questionEntry.Resize(fyne.NewSize(350, 50))

	question := ""
	answer := ""

	answerLabel := widget.NewLabel("")
	//answerLabel.Resize(fyne.NewSize(350, 400))

	speechModeRadio := widget.NewRadioGroup([]string{"Default", "Speech"}, func(s string) {
		fmt.Println("selected:", s)
		if s == "Speech" {
			// Show file open dialog
			dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err == nil && reader != nil {
					defer reader.Close()

					fmt.Println("Selected file:", reader.URI().Path())

					answer, err = sttRequest.SendRequest(reader.URI().Path())
					if err != nil {
						answerLabel.SetText(err.Error())
					}
					answerLabel.SetText(answer)
				}
			}, myWindow)
		}
	})
	//speechModeRadio.Resize(fyne.NewSize(350, 30))

	sendButton := widget.NewButton("Send", func() {
		question = questionEntry.Text
		questionParam := validateQuestion(question)

		switch questionParam {
		case "quit":
			myApp.Quit()
		case "", " ":
			return
		default:
			var err error
			answer, err = chatRequest.SendRequest(question)
			if err != nil {
				answerLabel.SetText(err.Error())
			}
			answerLabel.SetText(answer)
			questionEntry.SetText("")
		}
	})

	// UI layout
	hor := container.New(layout.NewVBoxLayout(), questionEntry, sendButton, speechModeRadio, answerLabel)
	content := container.New(layout.NewMaxLayout(), hor)

	myWindow.SetContent(content)

	myWindow.Resize(fyne.NewSize(600, 400))
	myWindow.Show()

	myApp.Run()
	time.Sleep(500 * time.Millisecond)
}

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func validateQuestion(question string) string {
	quest := strings.Trim(question, " ")
	keywords := []string{"", " "}
	for _, x := range keywords {
		if quest == x {
			return ""
		}
	}
	return quest
}
