package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	rootDir string
	db      Database
)

func setupHandlers(router *mux.Router) {
	GET := router.Methods("GET", "HEAD").Subrouter()
	POST := router.Methods("POST").Subrouter()

	GET.HandleFunc("/", httpIndex)
	GET.HandleFunc("/m/{id}", httpMeetup)

	POST.HandleFunc("/new", apiNew)
}

func dontListDirs(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/") {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			h.ServeHTTP(w, r)
		})
}

func main() {
	// get executable path
	exec, err := filepath.Abs(os.Args[0])
	if err != nil {
		panic(err)
	}
	rootDir = filepath.Dir(exec)

	// Command line parameters
	bind := flag.String("port", ":4242", "Address to bind to")
	mongo := flag.String("dburl", "localhost", "MongoDB servers, separated by comma")
	dbname := flag.String("dbname", "lily", "MongoDB database to use")
	flag.StringVar(&rootDir, "root", rootDir, "The HTTP server root directory")
	flag.Parse()

	db = InitDatabase(*mongo, *dbname)

	router := mux.NewRouter()
	setupHandlers(router)
	http.Handle("/", router)
	http.Handle("/static/", dontListDirs(http.StripPrefix("/static/",
		http.FileServer(http.Dir(rootDir+"/static")))))
	log.Printf("Listening on %s\r\nServer root: %s\r\n", *bind, rootDir)
	http.ListenAndServe(*bind, nil)
}
