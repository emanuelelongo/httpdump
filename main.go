package main

import (
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
)

func main() {
	portFlag := flag.String("port", "", "Listening port")
	strictPathFlag := flag.Bool("strict-path", false, "Preserve request path exactly as received (no normalization)")
	flag.Parse()

	port := *portFlag
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	if *strictPathFlag {
		log.Println("Strict path mode enabled: no path normalization")
		server := &http.Server{
			Handler: rawPathHandler(handler),
		}
		listener, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Listening on port %s", port)
		log.Fatal(server.Serve(listener))
	} else {
		log.Println("Strict path mode disabled: standard Go mux (may normalize)")
		http.HandleFunc("/", handler)
		log.Printf("Listening on port %s", port)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}
}

func rawPathHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawPath != "" {
			r.RequestURI = r.URL.RawPath
		}
		next.ServeHTTP(w, r)
	})
}
