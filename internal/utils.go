package internal

import (
	"fmt"
	"os"
)

type ManagedError struct {
	Message string
}

func PanicOnErr(e error) {
	if e != nil {
		panic(e)
	}
}

func CatchPanicAndExit() {
	if r := recover(); r != nil {
		if m, ok := r.(ManagedError); ok {
			fmt.Println(m.Message)
			fmt.Println("Exiting...")
			os.Exit(1)
		}
		panic(r)
	}
}

// Prepends the given elements to the beginning of a slice
func Prepend[T any](slice []T, elements ...T) []T {
	return append(elements, slice...)
}
