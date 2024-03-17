package internal

import (
	"fmt"
	"os/exec"
)

type stdOutHandler struct{}

func (handler *stdOutHandler) Output(text string) {
	fmt.Println(text)
}

func NewStdOutHandler() *stdOutHandler {
	return &stdOutHandler{}
}

type speakHandler struct{}

func (handler *speakHandler) Output(text string) {
	cmd := exec.Command("say", text)
	cmd.Output()
}

func NewSpeakHandler() *speakHandler {
	return &speakHandler{}
}
