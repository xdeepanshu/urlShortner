package main

import (
	"flag"
	"fmt"
	"github.com/xdeepanshu/urlShortner/store"
	"net/http"
	"strings"
	"time"
)

var (
	ds   = store.NewDataStore()
	form = `<html>
		<form method="POST" action="/add">
		URL: <input type="text" name="url"> 
		<button type="submit">Add</button> 
		</form>
		</html>
		`
	port = flag.String("port", ":8080", "port for runnig the server on")
)

func Add(rw http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	fmt.Println(url)
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
	go func(ch <-chan error) {
		fmt.Printf("Server is running ... on port %s\n", (*port)[1:])
		for {
			select {
			case <-time.After(time.Duration(20 * time.Second)):
				fmt.Printf("Server is running ... on port %s\n", (*port)[1:])
			case err := <-ch:
				fmt.Println("Error while listening: %s", err)
				return
			}
		}
	}(ch)

	http.HandleFunc("/", Redirect)
	http.HandleFunc("/add", Add)
	if err := http.ListenAndServe(*port, nil); err != nil {
		ch <- fmt.Errorf("Error %s while listening\n", err)
		return
	}
}
