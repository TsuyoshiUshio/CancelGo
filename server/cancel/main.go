package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func randomServer(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(3)
	w.WriteHeader(http.StatusOK)
	if value == 0 {
		time.Sleep(time.Second * 6)
		w.Write([]byte(fmt.Sprintf("You've got a slow response")))
	} else {
		w.Write([]byte(fmt.Sprintf("You've got a quick response")))
	}
}

func main() {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("%s not found\n", r.URL)))
	})
	router.HandleFunc("/api/{name}", randomServer).Methods("GET")
	fmt.Println("...Server stareted")
	fmt.Println("try http://localhost:39090/api/hello")
	http.ListenAndServe(":39090", router)

}
