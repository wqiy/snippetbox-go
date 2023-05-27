package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	// Define a new command-line flag
	addr := flag.String("addr", ":4000", "HTTP network address")
  // parse the command-line flag
	flag.Parse()


	// Use log.New() to create a logger for writing information messages.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	
	// Create a logger for writing error messages. 
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() function to register the file server as the handler
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Register the other application routes as normal.
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Initialize a new http.Server struct.
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
