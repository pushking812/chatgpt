module github.com/pushking812/chatgpt/cli

go 1.20

replace github.com/pushking812/chatgpt/gpt => ../gpt

require (
	github.com/pushking812/chatgpt/gpt v0.0.0
	github.com/spf13/cobra v1.6.1
)

require (
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/sashabaranov/go-openai v1.5.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)
