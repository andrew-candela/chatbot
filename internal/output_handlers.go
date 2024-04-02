package internal

import (
	"fmt"
	"os/exec"

	"github.com/chzyer/readline"
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

type AGIHandler struct {
	OtherLLM *LLM
}

func (handler *AGIHandler) Output(text string) {
	// make the user opt in to continue the insanity
	rl, err := readline.New("\n")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	rl.Readline()
	handler.OtherLLM.Converse(text)
}
