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
	http.HandleFunc("/", homePage)
	http.HandleFunc("/json-rpc", restAPI)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	startListener()
}
