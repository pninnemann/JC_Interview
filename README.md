# JumpCloud Interview Assignment

This is a small library written in golang which takes in actions and times and returns the average time for each unique action. Actions are input through the addAction
method by providing a json string of the form `{"action":"<String>", "time":<Number>}`. This can be called multiple times with the same or different actions
and will continue to keep a running tally of all actions inputed. Averages for all actions will be computed and returned through the getStats method as
a json string of the form `[{"action":"<String>", "avg":<Number>},{"action":"<String>", "avg":<Number>}]`.

To use this library, simply copy the runningStats.go file into your source directory and call the addAction or getStats methods from your own program.

Also Included with this library is a simple RestAPI interface to allow data to be input at runtime. See below for usage.

## Building and Running
Pre-req: Ensure Golang is installed on your computer

Run the API directly

```make run```

Create an executable in ./bin

```make build```

## Testing

Unit tests for the two library methods are included in `stats_test.go` to run these tests

```make test```

## Interacting with the API

Use Curl or Postman to call any of the available methods using json input `'{"method":"<method>, "params":<JSON INPUT>}'`

### getAPI
Returns a list of available methods

`curl -d '{"method":"getAPI"}' http://0.0.0.0:8080/json-rpc`

### addAction
Adds one or more action to the running tally, returns the number of actions added

`curl -d '{"method":"addAction", "params":{"action":"jump", "time":100}}' http://0.0.0.0:8080/json-rpc`

`curl -d '{"method":"addAction", "params":[{"action":"jump", "time":100}, {"action":"run", "time":75}, {"action":"jump", "time":200}]}' http://0.0.0.0:8080/json-rpc`

### getStats
Returns the averages of all the unique inputs

`curl -d '{"method":"getStats"}' http://0.0.0.0:8080/json-rpc`

### getTotals
Returns the running totals and the counts of all the unique inputs

`curl -d '{"method":"getTotals"}' http://0.0.0.0:8080/json-rpc`

### resetStats
Clears all the previous inputs

`curl -d '{"method":"resetStats"}' http://0.0.0.0:8080/json-rpc`