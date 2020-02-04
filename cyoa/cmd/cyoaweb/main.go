package main

import (
	"flag"
	"fmt"
	"gophercises/cyoa"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

func main() {
	port := flag.Int("port", 3000, "the port to start the CYOA web application on")
	file := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()
	fmt.Printf("Using the story in %s.\n", *file)

	f, err := os.Open(*file)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JSONStory(f)
	if err != nil {
		panic(err)
	}

	tpl := template.Must(template.ParseFiles("story-template.html"))
	h := cyoa.NewHandler(story,
		cyoa.WithPathFunc(pathFn),
		cyoa.WithTemplate(tpl))
	mux := http.NewServeMux()
	mux.Handle("/story/", h)
	fmt.Printf("starting the server at: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))

	fmt.Printf("%+v\n", story)
}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}

	return path[len("/story/"):]
}
