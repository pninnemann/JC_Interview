package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type ValidMethods struct {
	Methods []string `json:"methods"`
}

type Call struct {
	Method string `json:"method"`
	Params json.RawMessage `json:"params"`
	Res http.ResponseWriter
}

func restAPI(res http.ResponseWriter, req *http.Request, )  {
	var newCall Call
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil || len(reqBody) == 0 {
		res.WriteHeader(400)
		log.Println("Body Cannot be empty")
		fmt.Fprintf(res, "Body Cannot be empty")
		return
	}

	err = json.Unmarshal(reqBody, &newCall)
	if err != nil {
		res.WriteHeader(400)
		log.Println("Unable to parse JSON Data")
		fmt.Fprintf(res, "Unable to Parse JSON Data")
		return
	}
	newCall.Res = res
	method := reflect.ValueOf(&newCall).MethodByName(strings.Title(newCall.Method))
	if !method.IsValid() {
		log.Printf("Method %s Does not exist\n", newCall.Method)
		res.WriteHeader(400)
		fmt.Fprintf(res, "Method %s does not exist, run getAPI for a list of valid Methods", newCall.Method)
		return
	}

	method.Call([]reflect.Value{})
}

func (call *Call) AddAction() {
	var newActions []ActionInput
	var newAction ActionInput

	if len(call.Params) == 0 {
		log.Println("Method 'addAction' must be called with Params")
		call.Res.WriteHeader(400)
		fmt.Fprintf(call.Res, "Method 'addAction' must be called with Params")
		return
	}

	if call.Params[0] == '[' {
		log.Println("Multiple Actions added")
		json.Unmarshal(call.Params, &newActions)
	} else {
		log.Println("One Action Added")
		json.Unmarshal(call.Params, &newAction)
		newActions = append(newActions, newAction)
	}

	var count int
	for _, action := range newActions {
		resJson, err := json.Marshal(action)
		if err != nil {
			log.Println("Error occured marshalling validMethods", err)
			return
		}
		addAction(string(resJson))
		count++
	}

	log.Println(count, "new actions added")
	call.Res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(call.Res, "{ \"actionsAdded\": %d }", count)
}

func (call *Call) GetStats(){
	log.Println("GetStats Called")

	if len(totals) == 0 {
		log.Println("No Stats logged")
		call.Res.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(call.Res, "[]")
		return
	}

	output := getStats()

	log.Println(output)
	call.Res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(call.Res, output)
}

func (call *Call) GetTotals() {
	log.Println("GetTotals Called")
	var totals []ActionTotal

	if len(totals) == 0 {
		log.Println("No Stats logged")
		call.Res.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(call.Res, "[]")
		return
	}

	Mux.Lock()
	for _, action := range totals {
		totals = append(totals, action)
	}
	Mux.Unlock()

	resJson, err := json.Marshal(totals)
	if err != nil {
		log.Println("Error occured marshalling validMethods", err)
		call.Res.WriteHeader(500)
		return
	}
	log.Println(string(resJson))
	call.Res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(call.Res, string(resJson))
}

func (call *Call) ResetStats() {
	log.Println("ResetStats Called")
	Mux.Lock()
	for action := range totals {
		delete(totals, action)
	}
	Mux.Unlock()
	log.Println(totals)
}

func (call *Call) GetAPI() {
	callType := reflect.TypeOf(call)
	var validMethods ValidMethods
	for i := 0; i < callType.NumMethod(); i++ {
		validMethods.Methods = append(validMethods.Methods, callType.Method(i).Name)
	}

	resJson, err := json.Marshal(validMethods)
	if err != nil {
		log.Println("Error occured marshalling validMethods", err)
		call.Res.WriteHeader(500)
		return
	}
    log.Println(string(resJson))
	call.Res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(call.Res, string(resJson))
}
