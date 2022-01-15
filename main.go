package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Option struct {
	Text string `json:"text"`
	Arc string `json:"arc"`
}
type Arc struct {
	Title string `json:"title"`
	Story []string `json:"story"`
	Options []Option `json:"options"`
}
type HTMLHandler struct {
	Arc Arc
	Template *template.Template
}
func(h HTMLHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if err := h.Template.Execute(rw, h.Arc); err != nil {
			fmt.Println(err)
			log.Fatal("could not execute template to response writer")
	}
}
func NewHandler(a Arc, t *template.Template) HTMLHandler {
	return HTMLHandler{a, t}
}

func main() {
	mux := http.NewServeMux()
	arcs := createArcs()
	createHandlers(mux, arcs)
	http.ListenAndServe(":3000", mux)
}

func createArcs() map[string]Arc {
	data, err := os.ReadFile("gopher.json")
	if err != nil {
		log.Fatal("Could not read json file")
	}
	type Arcs map[string]Arc
	var arcs Arcs
	err = json.Unmarshal(data, &arcs)
	for name := range arcs {
		fmt.Println(name)
	}
	if err != nil {
		log.Fatal("could not unmarshal json file")
	}
	return arcs
}
func createHandlers(s *http.ServeMux, arcs map[string]Arc) {
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		log.Fatal("could not read template file")
	}
	for name, a := range arcs {
		(*s).Handle("/"+name, NewHandler(a, tmpl))
	}
	(*s).HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		http.Redirect(rw, r, "/intro", http.StatusSeeOther)
	})
}