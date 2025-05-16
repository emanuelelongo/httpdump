package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func main() {
	portFlag := flag.String("port", "", "Listening port")
	flag.Parse()

	port := *portFlag
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, false)
		if err != nil {
			http.Error(w, "Error dumping the request", http.StatusInternalServerError)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading the request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		log.Printf("%s\n\n%s", dump, string(body))

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write(dump)
	})

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
