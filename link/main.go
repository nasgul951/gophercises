package main

import (
	"os"

	"golang.org/x/net/html"
)

func main() {
	r, err := os.Open("ex1.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	t := html.NewTokenizer(r)

}
