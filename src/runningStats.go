package main

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
)

type ActionTotal struct {
	Action string
	Time float64
	Count int
}

type ActionInput struct {
	Action string `json:"action"`
	Time float64 `json:"time"`
}

type ActionStat struct {
	Action string `json:"action"`
	Avg float64 `json:"avg"`
}

var Mux sync.Mutex
var totals = make(map[string]ActionTotal)

/**
 * This function accepts a json serialized string of the form below and maintains an average time
 * for each action.
 *
 * @parms [in] string    json data of the form {"action":"<String>", "time":<Number>}
 * returns     error     in the case json data is unable to be parsed
 */
func addAction (input string) error {
	var newAction ActionInput
	err := json.Unmarshal([]byte(input), &newAction)
	if err != nil {
		return errors.New("Unable to parse JSON Data")
	}
	Mux.Lock()
	if actionTotal, ok := totals[newAction.Action]; ok {
		actionTotal.Count++
		actionTotal.Time += newAction.Time
		totals[newAction.Action] = actionTotal
	} else {
		actionTotal := ActionTotal{newAction.Action,newAction.Time, 1}
		totals[newAction.Action] = actionTotal
	}
	Mux.Unlock()
	return nil
}

/**
 * Function that accepts no input returns a serialized json array of the average
 * time for each action that has been provided to the addAction function.
 *
 * returns     String     json array of average times or the form [{"action":"<String>", "avg":<Number>}]
 */
func getStats() string {
	var stats []ActionStat
	Mux.Lock()
	for _, action := range totals {
		var stat = ActionStat{action.Action, action.Time/float64(action.Count)}
		stats = append(stats, stat)
	}
	Mux.Unlock()

	resJson, err := json.Marshal(stats)
	if err != nil {
		log.Println("Error occured marshalling stats", err)
		return "Error occured marshalling stats"
	}
	return string(resJson)
}