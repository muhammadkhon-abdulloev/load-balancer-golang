package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(5)
	go func() {
		serveBackend("web1", ":5001")
		wg.Done()
	}()

	go func() {
		serveBackend("web2", ":5002")
		wg.Done()
	}()

	go func() {
		serveBackend("web3", ":5003")
		wg.Done()
	}()

	timer := time.NewTimer(time.Second * 10)
	go func() {
		serveBackendWithSleep("web4", ":5004", *timer)
		wg.Done()
	}()
	
	go func() {
		
		serveBackend("web5", ":5005")
		wg.Done()
	}()
	wg.Wait()
}

func serveBackend(name string, port string) {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Backend server name:%v\n", name)
		fmt.Fprintf(w, "Response header:%v\n", r.Header)
	}))
	http.ListenAndServe(port, mux)
}

func serveBackendWithSleep(name string, port string, sleep time.Timer) {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Backend server name:%v\n", name)
		fmt.Fprintf(w, "Response header:%v\n", r.Header)
	}))
	http.ListenAndServe(port, mux)
}