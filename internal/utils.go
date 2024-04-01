package internal

import (
	"fmt"
	"io"
	"net/http"
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

// Inspects the HTTP response and panics if there is anything
// a 4xx or 5xx response.
// Includes the reponse body in the error message
func InspectAPIResponsePayload(resp *http.Response) {
	if resp.StatusCode < 400 {
		return
	}
	// There is actually an error, so let's inspect the body
	// and return a helpful message
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body) // ignore any error reading body
	message := fmt.Sprintf("got status code %d when hitting %v: %v", resp.StatusCode, resp.Request.URL, string(body))
	panic(ManagedError{message})
}

// picks the first non zero value passed
func Coalesce[T string | int](values ...T) T {
	null_val := *new(T)
	for _, value := range values {
		if value != null_val {
			return value
		}
	}
	return null_val
}
