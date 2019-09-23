// +build ignore

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/vugu/vugu/simplehttp"
)

func main() {
	wd, _ := os.Getwd()
	l := "127.0.0.1:8844"
	log.Printf("Starting HTTP Server at %q", l)
	vuguh := simplehttp.New(wd, true)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.URL.Path {
		case "/api/startup":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"startup":1}`)
			return

		case "/api/somedata":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"somedata":2}`)
			return

		case "/api/somemoredata":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"somemoredata":3}`)
			return

		}

		vuguh.ServeHTTP(w, r)
	})

	// include a CSS file
	// simplehttp.DefaultStaticData["CSSFiles"] = []string{ "/my/file.css" }
	log.Fatal(http.ListenAndServe(l, h))
}
