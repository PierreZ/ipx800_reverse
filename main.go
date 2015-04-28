package main

import (
	"encoder/xml"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type IPX800 struct {
	Led0 int `xml:"led0"`
	Led0 int `xml:"led0"`
	Led0 int `xml:"led0"`
	Led0 int `xml:"led0"`
	Led0 int `xml:"led0"`
}

type Proxy struct {
	Dsnames        []string  `json:"dsnames"`
	Dstypes        []string  `json:"dstypes"`
	Host           string    `json:"host"`
	Interval       float64   `json:"interval"`
	Plugin         string    `json:"plugin"`
	PluginInstance string    `json:"plugin_instance"`
	Time           float64   `json:"time"`
	Type           string    `json:"type"`
	TypeInstance   string    `json:"type_instance"`
	Values         []float64 `json:"values"`
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", handle_ipx800).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe("10.0.2.15:8085", r)
}

func handle_ipx800(w http.ResponseWriter, r *http.Request) {

	var ipx IPX800
	fmt.Println(r.Body)
	if err := xml.NewDecoder(r.Body).Decode(&ipx); err != nil {
		http.Error(w, "Error during decode", http.StatusBadRequest)
		return
	}

	// TODO: trnasform IPX800 to proxy

}
