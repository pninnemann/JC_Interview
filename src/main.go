package main

import (
	"fmt"
	"log"
	"net/http"
)

func homePage(res http.ResponseWriter, req *http.Request)  {
	fmt.Fprint(res, "Homepage Hit")
}

func startListener() {
	log.Println("RestAPI Server running on http://0.0.0.0:8080/json-rpc")
	http.HandleFunc("/", homePage)
	http.HandleFunc("/json-rpc", restAPI)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	startListener()
}
