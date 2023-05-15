package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	order_id := params["orderID"]
	log.Printf("received order for id: %s", order_id)
	ctx := r.Context()

	// simulate work
	doWork(ctx)

	json := simplejson.New()
	json.Set("order_id", order_id)
	json.Set("status", "received")

	payload, err := json.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func doWork(ctx context.Context) {
	r := rand.Intn(1729)
	time.Sleep(time.Duration(r) * time.Microsecond)
}

func catchAllHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Invalid URI path. Use /shipping/{orderID} to create shipping for order"))
}

// Route declaration
func router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/shipping/{orderID}", handler)

	r.HandleFunc("/", catchAllHandler)

	return r
}

// Initiate web server
func main() {
	router := router()

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:9100",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
