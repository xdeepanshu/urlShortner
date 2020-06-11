package main

import (
	"fmt"
	"github.com/xdeepanshu/urlShortner/store"
	"net/http"
	"time"
)

var (
	ds = store.NewDataStore()
)

func Add(rw http.ResponseWriter, r *http.Request) {

}
func Redirect(rw http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(rw, "Hello World")
}

func main() {
	ch := make(chan error)
	go func(ch <-chan error) {
		for {
			select {
			case <-time.After(time.Duration(5 * time.Second)):
				fmt.Println("Server is running .....")
			case err := <-ch:
				fmt.Println("Error while listening: %s", err)
			}
		}
	}(ch)

	http.HandleFunc("/", Redirect)
	http.HandleFunc("/add", Add)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		ch <- fmt.Errorf("Error %s while listening\n", err)
		return
	}
}
