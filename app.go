package main

import (
	"log"
	"net/http"
	"os"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"fmt"
)

func main() {
	router := httprouter.New()

	router.GET("/", home)

	fmt.Printf("Hello, world.\n")
   log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

func home(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	widgets := Widgets{
		Widget{Name: "Foo", Description: "This is a foo"},
		Widget{Name: "Bar", Description: "This is a bar"},
		Widget{Name: "Baz", Description: "This is a Baz"},
	}

	cj, _ := json.Marshal(widgets)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", cj)
}