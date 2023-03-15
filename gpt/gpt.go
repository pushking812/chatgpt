package gpt

import (
	"fmt"
	"time"
)

type Request struct {
	h        Handler
	apiKey   string
	timeout  string
	answerCh chan string
	errCh    chan error
}

func NewRequestType(hr Handler, ak, to string) *Request {
	return &Request{
		h:        hr,
		apiKey:   ak,
		timeout:  to,
		answerCh: make(chan string),
		errCh:    make(chan error),
	}
}

func (r *Request) SendRequest(question string) (string, error) {
	if question == "" {
		panic("Message is empty")
	}

	go func() {
		answer, err := r.h(r.apiKey, question, r.timeout)
		if err != nil {
			r.errCh <- err
		} else {
			r.answerCh <- answer
		}
	}()

	for {
		select {
		case answer := <-r.answerCh:
			return answer, nil
		case err := <-r.errCh:
			return "", fmt.Errorf("error: %s", err)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

type Handler = func(apikey, question, timeout string) (string, error)
