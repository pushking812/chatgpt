package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pushking812/chatgpt/gpt"
	cc "github.com/pushking812/chatgpt/gpt/chat-complete"
	tt "github.com/pushking812/chatgpt/gpt/speech-to-text"
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

	chatRequest := gpt.NewRequestType(cc.GetAnswer, apiKey, "15s")
	sttRequest := gpt.NewRequestType(tt.GetAnswer, apiKey, "15s")

	history := []string{}

	rootCmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Chat with ChatGPT in console.",
		Run: func(cmd *cobra.Command, args []string) {
			scanner := bufio.NewScanner(os.Stdin)
			quit := false
			speech := false
			answer := ""

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
				//questionParam := validateQuestion(question)

				switch question {
				case "quit":
					quit = true

					fmt.Println("\nHistory:")
					for i, v := range history {
						fmt.Printf("\n#%d\n%v\n", i+1, v)
					}

				case "speech":
					speech = true
				case "", " ":
					speech = false
				default:
					var err error
					if speech {
						answer, err = sttRequest.SendRequest(question)
					} else {
						answer, err = chatRequest.SendRequest(question)
					}

					if err != nil {
						fmt.Println("error:", err)
					}

					answer = strings.ReplaceAll(answer, "\n", "")
					answer = strings.ReplaceAll(answer, "\r", "")

					history = append(history, "question: "+question+"\nanswer: "+answer)

					fmt.Printf("output: %s\n", answer)
				}
			}
		},
	}

	log.Fatal(rootCmd.Execute())
}

// func validateQuestion(question string) string {
// 	quest := strings.Trim(question, " ")
// 	keywords := []string{"", " "}
// 	for _, x := range keywords {
// 		if quest == x {
// 			return ""
// 		}
// 	}
// 	return quest
// }
