package main

import (
	"fmt"
	"log"
	"net/http"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Port  int    `long:"port" short:"p" description:"Port" default:"8080"`
	Realm string `long:"message" short:"m" description:"WWW-Authenticate realm"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	addr := fmt.Sprintf(":%d", opts.Port)

	http.HandleFunc("/login", login)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	user, pass, ok := r.BasicAuth()
	if !ok {
		unauthorized(w, "Admin Panel")
		return
	} else if user != "user" && pass != "pass" {
		unauthorized(w, "Admin Panel")
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func unauthorized(w http.ResponseWriter, realm string) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	w.WriteHeader(401)
}
