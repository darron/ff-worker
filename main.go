package main

import (
	"io"
	"log"
	"net/http"

	"github.com/darron/ff/config"
	"github.com/syumai/workers"
)

type HTTPService struct {
	conf *config.App
}

func main() {
	// Get conf
	var opts []config.OptFunc
	opts = append(opts, config.WithPort("8000"))
	opts = append(opts, config.WithLogger("debug", "text"))
	conf, err := config.Get(opts...)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		msg := "Hello!"
		w.Write([]byte(msg))
	})
	http.HandleFunc("/echo", func(w http.ResponseWriter, req *http.Request) {
		io.Copy(w, req.Body)
	})
	workers.Serve(nil) // use http.DefaultServeMux
}
