package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", handle_ipx800).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe("10.0.2.15:8085", nil)
}

func handle_ipx800(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Body)
	// if err := json.NewDecoder(r.Body).Decode(&ch); err != nil {
	// 	http.Error(w, "bad channel", http.StatusBadRequest)
	// 	return
	// }

}
