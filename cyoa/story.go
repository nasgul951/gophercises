package cyoa

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
)

func init() {
	tpl = template.Must(template.ParseFiles("template.html"))
}

var tpl *template.Template

// HandlerOption : to modify handler options
type HandlerOption func(h *handler)

// WithTemplate : specify a custom template
func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

//WithPathFunc : specify a path function
func WithPathFunc(f func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = f
	}
}

// NewHandler : get a new handler for the story
func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl, defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)

	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}

// JSONStory : parse json into a Story
func JSONStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

// Story : the story
type Story map[string]Chapter

// Chapter : a chapter
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// Option : an option
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
