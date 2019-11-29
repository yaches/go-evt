package main

import (
	"errors"
	"fmt"
	"os"

	"go-evt/evt"
)

func main() {
	if len(os.Args) != 2 {
		panic(errors.New("usage"))
	}

	file, err := os.Open(os.Args[1])
	defer file.Close()
	if err != nil {
		panic(err)
	}

	h, records, err := evt.ParseEvt(file)
	if err != nil {
		panic(err)
	}
	fmt.Println(h)
	for _, r := range records {
		fmt.Println(r)
	}
}
