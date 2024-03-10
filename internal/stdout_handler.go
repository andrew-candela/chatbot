package internal

import "fmt"

type stdOutHandler struct{}

func (handler *stdOutHandler) Output(text string) {
	fmt.Println(text)
}

func NewStdOutHandler() *stdOutHandler {
	return &stdOutHandler{}
}
