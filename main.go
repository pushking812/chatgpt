package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	cc "github.com/pushking812/workout/chatgpt/chat-complete"
	tt "github.com/pushking812/workout/chatgpt/gpt-speech-to-text"
	"github.com/spf13/cobra"
)

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func main() {
	log.SetOutput(new(NullWriter))
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("Missing API KEY")
	}

	answerCh := make(chan string)
	errCh := make(chan error)

	var h handler

	rootCmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Chat with ChatGPT in console.",
		Run: func(cmd *cobra.Command, args []string) {
			scanner := bufio.NewScanner(os.Stdin)
			quit := false
			speech := false

			fmt.Println("Type your question to ChatGPT or `speech` to set speech-to-text mod (type `quit` to exit)")

			for !quit {
				fmt.Println("")
				if speech {
					fmt.Print("input mp3 filename: ")
				} else {
					fmt.Print("input question: ")
				}

				if !scanner.Scan() {
					break
				}

				question := scanner.Text()
				questionParam := validateQuestion(question)

				switch questionParam {
				case "quit":
					quit = true
				case "speech":
					speech = true
				case "", " ":
					speech = false
				default:
					h = cc.GetAnswer

					if speech {
						h = tt.GetAnswer
					}

					go func() {
						answer, err := h(apiKey, question)
						if err != nil {
							errCh <- err
						} else {
							answerCh <- answer
						}
					}()

					for done := false; !done; {
						select {
						case answer := <-answerCh:
							fmt.Printf("output: %s\n", answer)
							done = true
						case err := <-errCh:
							fmt.Println("Error:", err)
							done = true
						default:
							time.Sleep(100 * time.Millisecond)
						}
					}
				}
			}
		},
	}

	log.Fatal(rootCmd.Execute())
}

type handler = func(apikey, question string) (string, error)

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
