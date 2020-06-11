package main

import (
	"flag"
	"fmt"
	"github.com/xdeepanshu/urlShortner/store"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	logger = log.New(os.Stdout, "url-shortner", log.LstdFlags)
	ds     = store.NewDataStore("storage.gob", logger)
	form   = `<html>
		<form method="POST" action="/add">
		URL: <input type="text" name="url"> 
		<button type="submit">Add</button> 
		</form>
		</html>
		`
	port       = flag.String("port", ":8080", "port for runnig the server on")
	serverEcho = flag.Int("echo", 10, "time in secs for echo of server status")
)

func Add(rw http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	logger.Println(url)
	protoPresent := strings.Contains(url, "http://") || strings.Contains(url, "https://")
	if url == "" || !protoPresent {
		fmt.Fprintln(rw, form)
		return
	}
	short := ds.Put(url)
	fmt.Fprintf(rw, "http://localhost%s/%s", *port, short)

}
func Redirect(rw http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	url, err := ds.Get(key)
	if err != nil {
		http.Error(rw, fmt.Sprint(err), http.StatusNotFound)
		return
	}
	http.Redirect(rw, r, url, http.StatusFound)
}

func main() {
	flag.Parse()
	ch := make(chan error)
	go func(ch <-chan error, logger *log.Logger) {
		logger.Printf("Server is running ... on port %s\n", (*port)[1:])
		for {
			select {
			case <-time.After((time.Duration(*serverEcho)) * time.Second):
				logger.Printf("Server is running ... on port %s\n", (*port)[1:])
			case err := <-ch:
				logger.Println("Error while listening: %s", err)
				return
			}
		}
	}(ch, logger)

	http.HandleFunc("/", Redirect)
	http.HandleFunc("/add", Add)
	if err := http.ListenAndServe(*port, nil); err != nil {
		ch <- fmt.Errorf("Error %s while listening\n", err)
		return
	}
}
