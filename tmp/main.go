package main

import (
	"errors"
	"fmt"
)

func main() {
	err := baseString("errorf")
	if err != nil {
		fmt.Println(fmt.Errorf("%w", err))
	}
	PrintString("hello")
}

func baseString(hello string) error {
	if hello == "hello" {
		return nil
	} else {
		return errors.New("error")
	}
}

func PrintString(hello string) {
	for range hello {
		fmt.Println(hello)
	}
}
